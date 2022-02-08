package db

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"sync"
)

type etcdClient struct {
	cli *clientv3.Client
}

type etcdOption struct {
	endpoints []string
	username  string
	password  string
}

var etcdFactory *etcdClient
var once sync.Once

func GetEtcdFactory(opt *etcdOption) (*etcdClient, error) {
	if opt == nil && etcdFactory == nil {
		return nil, fmt.Errorf("get etcd store faild")
	}

	var err error
	once.Do(func() {

		var cli *clientv3.Client

		e := &etcdClient{}

		cli, err = clientv3.New(clientv3.Config{
			Endpoints: opt.endpoints,
			Username:  opt.username,
			Password:  opt.password,
			DialOptions: []grpc.DialOption{
				grpc.WithBlock(),
			},
		})

		if err != nil {
			return
		}

		e.cli = cli

	})

	return nil, nil
}
