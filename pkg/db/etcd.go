package db

import (
	"context"
	"errors"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

// rewrite etcd.io
type etcdOptions struct {
	endpoints       []string
	requestTimeout  time.Duration
	leaseTTLTimeout int64

	timeout   time.Duration
	namespace string
}

type etcdClient struct {
	cli             *clientv3.Client
	requestTimeout  time.Duration
	leaseTTLTimeout int64

	namespace   string
	leaseID     clientv3.LeaseID
	leaseLiving bool

	onKeepaliveFailure func()
}

var etcdFactory *etcdClient
var once sync.Once

// NewEtcdClient singleton client etcd v3
// custom etcd client
func NewEtcdClient(opt *etcdOptions, onKeepaliveFailure func()) (*etcdClient, error) {
	if opt == nil && etcdFactory == nil {
		return nil, errors.New("failed to get etcd store factory")
	}

	var err error
	once.Do(func() {

		var cli *clientv3.Client

		e := &etcdClient{}

		if onKeepaliveFailure == nil {
			onKeepaliveFailure = defaultOnKeepaliveFailure
		}

		e.onKeepaliveFailure = onKeepaliveFailure
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   opt.endpoints,
			DialTimeout: opt.timeout,

			DialOptions: []grpc.DialOption{
				grpc.WithBlock(),
			},
		})

		if err != nil {
			return
		}

		e.cli = cli
		e.requestTimeout = opt.requestTimeout
		e.leaseTTLTimeout = opt.leaseTTLTimeout
		e.namespace = opt.namespace

		if err = e.startSession(); err != nil {
			if err = e.close(); err != nil {
				log.Printf("etcdStore client close failed %s", err)
			}

			return
		}

		etcdFactory = e
	})

	if etcdFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get etcd store fatory, etcdFactory: %+v, error: %w", etcdFactory, err)
	}

	return etcdFactory, err
}

// startSession keepalive lease
func (c *etcdClient) startSession() error {
	// context
	ctx := context.TODO()
	resp, err := c.cli.Grant(ctx, c.leaseTTLTimeout)

	if err != nil {
		return errors.New("creates new lease failed")
	}

	c.leaseID = resp.ID
	ch, err := c.cli.KeepAlive(ctx, c.leaseID)

	if err != nil {
		return errors.New("keep alive failed")
	}

	go func() {

		for {
			if _, ok := <-ch; !ok {
				c.leaseLiving = false
				log.Println("fail to keepalive session")

				if c.onKeepaliveFailure != nil {
					c.onKeepaliveFailure()
				}

				break
			}
		}
	}()

	return nil
}

// close etcd client
func (c *etcdClient) close() error {

	if c.cli != nil {
		return c.cli.Close()
	}

	return nil
}

// default
func defaultOnKeepaliveFailure() {
	log.Println("etcdStore keepalive failed")
}
