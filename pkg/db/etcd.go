package db

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

type etcdOnCreateFunc func(ctx context.Context, key, value []byte)

type etcdOnModifyFunc func(ctx context.Context, key, preValue, value []byte)

type etcdOnDeleteFunc func(ctx context.Context, key []byte)

type etcdClient struct {
	cli             *clientv3.Client
	requestTimeout  time.Duration
	leaseTTLTimeout int

	leaseID     clientv3.LeaseID
	leaseLiving bool

	watchers map[string]*etcdWatchers

	onKeepaliveFailure func()
	stopKeepaliveFunc  func()
	namespace          string
}

type etcdOptions struct {
	Endpoints      []string
	UserName       string
	Password       string
	Timeout        int
	RequestTimeout int
	LeaseExpire    int
	Namespace      string
}

type etcdWatchers struct {
	watcher    clientv3.Watcher
	cancelFunc context.CancelFunc
}

var (
	etcdFactory *etcdClient
	once        sync.Once
)

// GetEtcdFactory Get Etcd Factory
// etcd lease
func GetEtcdFactory(opts *etcdOptions, defaultOnKeepalive func()) (*etcdClient, error) {
	if opts == nil && etcdFactory == nil {
		return nil, fmt.Errorf("get etcd store faild")
	}

	var err error
	once.Do(func() {
		var cli *clientv3.Client
		e := &etcdClient{}

		if defaultOnKeepalive == nil {
			defaultOnKeepalive = defaultKeepaliveFailure
		}

		e.onKeepaliveFailure = defaultOnKeepalive
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   opts.Endpoints,
			DialTimeout: time.Duration(opts.Timeout) * time.Second,
			Username:    opts.UserName,
			Password:    opts.Password,
			DialOptions: []grpc.DialOption{grpc.WithBlock()},
		})

		if err != nil {
			return
		}
		e.cli = cli
		e.requestTimeout = time.Duration(opts.RequestTimeout) * time.Second
		e.leaseTTLTimeout = opts.LeaseExpire
		e.watchers = make(map[string]*etcdWatchers)
		e.namespace = opts.Namespace

		if err = e.startSession(); err != nil {
			if err = e.Close(); err != nil {
				log.Println("etcdClient closed failed")
			}
			return
		}

		etcdFactory = e
	})

	if etcdFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get etcd factory error")
	}

	return etcdFactory, nil
}

func (e *etcdClient) startSession() error {
	ctx, cancelFunc := context.WithCancel(context.TODO())
	e.stopKeepaliveFunc = cancelFunc

	resp, err := e.cli.Grant(ctx, int64(e.leaseTTLTimeout))
	if err != nil {
		return err
	}

	e.leaseID = resp.ID
	ch, err := e.cli.KeepAlive(ctx, e.leaseID)
	if err != nil {
		return err
	}

	go func() {
		for {
			if _, ok := <-ch; !ok {
				e.leaseLiving = false
				log.Println("failed keepalive session")
				if e.onKeepaliveFailure != nil {
					e.onKeepaliveFailure()
				}
				break
			}
		}
	}()

	return nil
}

func (e *etcdClient) getKeyByNamespace(prefix string) string {
	if len(e.namespace) == 0 {
		return prefix
	}
	return fmt.Sprintf("%s%s", e.namespace, prefix)
}

func (e *etcdClient) Close() error {
	if e.cli == nil {
		return nil
	}
	return e.cli.Close()
}

func (e *etcdClient) Client() *clientv3.Client {
	return e.cli
}

func (e *etcdClient) SessionLiving() bool {
	return e.leaseLiving
}

func (e *etcdClient) RestartSession() error {
	if e.leaseLiving {
		return fmt.Errorf("session is living")
	}
	return e.startSession()
}

func (e *etcdClient) PutKv(ctx context.Context, key, value string, session bool) error {
	gCtx, cancelFunc := context.WithTimeout(ctx, e.requestTimeout)
	defer cancelFunc()
	key = e.getKeyByNamespace(key)
	if session {
		if _, err := e.cli.Put(gCtx, key, value, clientv3.WithLease(e.leaseID)); err != nil {
			return err
		}
		return nil
	}

	if _, err := e.cli.Put(gCtx, key, value); err != nil {
		return err
	}

	return nil
}

func (e *etcdClient) GetKv(ctx context.Context, key string) ([]byte, error) {
	gCtx, cancelFunc := context.WithTimeout(ctx, e.requestTimeout)
	defer cancelFunc()
	key = e.getKeyByNamespace(key)
	getResp, err := e.cli.Get(gCtx, key)
	if err != nil {
		return nil, err
	}

	if len(getResp.Kvs) == 0 {
		return nil, fmt.Errorf("key isn't exist")
	}

	return getResp.Kvs[0].Value, nil
}

func (e *etcdClient) DeleteKv(ctx context.Context, key string) ([]byte, error) {
	gCtx, cancelFunc := context.WithTimeout(ctx, e.requestTimeout)
	defer cancelFunc()

	key = e.getKeyByNamespace(key)
	dResp, err := e.cli.Delete(gCtx, key, clientv3.WithPrevKV())

	if err != nil {
		return nil, err
	}

	if dResp.Deleted == 1 {
		return dResp.PrevKvs[0].Value, nil
	}
	return nil, nil
}

func (e *etcdClient) Watch(ctx context.Context, prefix string, onCreate etcdOnCreateFunc, onModify etcdOnModifyFunc, onDelete etcdOnDeleteFunc) error {
	// if isn't exist
	if _, ok := e.watchers[prefix]; ok {
		return fmt.Errorf("watcher prefix %s already registed", prefix)
	}

	gctx, cancelFunc := context.WithCancel(ctx)
	watcher := clientv3.NewWatcher(e.cli)

	e.watchers[prefix] = &etcdWatchers{
		watcher:    watcher,
		cancelFunc: cancelFunc,
	}

	prefix = e.getKeyByNamespace(prefix)
	watchChan := e.cli.Watch(gctx, prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())
	go func() {
		for watchResponse := range watchChan {
			for _, event := range watchResponse.Events {
				keyByte := event.Kv.Key[len(e.namespace):]
				if event.PrevKv == nil {
					// trigger onCreate
					onCreate(gctx, keyByte, event.Kv.Value)
				} else {
					switch event.Type {
					case mvccpb.PUT:
						onModify(gctx, keyByte, event.PrevKv.Value, event.Kv.Value)
					case mvccpb.DELETE:
						if onDelete != nil {
							onDelete(gctx, keyByte)
						}
					}
				}
			}
		}
		log.Printf("stopped watch %s", prefix)
	}()

	return nil
}

func (e *etcdClient) UnWatch(prefix string) error {
	if watcher, ok := e.watchers[prefix]; ok {
		delete(e.watchers, prefix)
		return watcher.Cancel()
	}
	return nil
}

func (w *etcdWatchers) Cancel() error {
	w.cancelFunc()
	return w.watcher.Close()
}

func defaultKeepaliveFailure() {
	log.Println("failed keepalive lease etcd server")
}
