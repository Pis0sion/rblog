package opts

import (
	"github.com/Pis0sion/rblog/pkg/vips"
	"time"
)

type Opts struct {
	SrvOpts   *SrvOpts
	AppOpts   *AppOpts
	MysqlOpts *MysqlOpts
	RedisOpts *RedisOpts
	GrpcOpts  *GrpcOpts
}

func NewOptions() (*Opts, error) {

	var (
		srvOpts   *SrvOpts
		appOpts   *AppOpts
		mysqlOpts *MysqlOpts
		redisOpts *RedisOpts
		grpcOpts  *GrpcOpts
	)

	vp, err := vips.NewViper("config", "config", "yaml")

	if err != nil {
		return nil, err
	}

	if err = vp.ReadSection("Server", &srvOpts); err != nil {
		return nil, err
	}

	if err = vp.ReadSection("App", &appOpts); err != nil {
		return nil, err
	}

	if err = vp.ReadSection("DataBase", &mysqlOpts); err != nil {
		return nil, err
	}

	if err = vp.ReadSection("Redis", &redisOpts); err != nil {
		return nil, err
	}

	if err = vp.ReadSection("Grpc", &grpcOpts); err != nil {
		return nil, err
	}

	srvOpts.ReadTimeout *= time.Second
	srvOpts.WriteTimeout *= time.Second
	mysqlOpts.MaxConnectionLifeTime *= time.Minute

	return &Opts{
		SrvOpts:   srvOpts,
		AppOpts:   appOpts,
		MysqlOpts: mysqlOpts,
		RedisOpts: redisOpts,
		GrpcOpts:  grpcOpts,
	}, nil
}
