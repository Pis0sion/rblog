package db

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

type etcdClient struct {
	cli             *clientv3.Client
	requestTimeout  time.Duration
	leaseTTLTimeout int

	leaseID     clientv3.LeaseID
	leaseLiving bool

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
	gctx, cancelFunc := context.WithTimeout(ctx, e.requestTimeout)
	defer cancelFunc()

	if session {
		if _, err := e.cli.Put(gctx, key, value, clientv3.WithLease(e.leaseID)); err != nil {
			return err
		}
		return nil
	}

	if _, err := e.cli.Put(gctx, key, value); err != nil {
		return err
	}

	return nil
}

func (e *etcdClient) GetKv(ctx context.Context, key string) ([]byte, error) {
	gctx, cancelFunc := context.WithTimeout(ctx, e.requestTimeout)
	defer cancelFunc()

	getResp, err := e.cli.Get(gctx, key)
	if err != nil {
		return nil, err
	}

	if len(getResp.Kvs) == 0 {
		return nil, fmt.Errorf("key isn't exist")
	}

	return getResp.Kvs[0].Value, nil
}

func defaultKeepaliveFailure() {
	log.Println("failed keepalive lease etcd server")
}
