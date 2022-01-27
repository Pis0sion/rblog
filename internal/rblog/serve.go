package rblog

import (
	"context"
	"fmt"
	"github.com/Pis0sion/rblog/internal/pkg/serve"
	pb "github.com/Pis0sion/rblog/internal/proto/v1"
	"github.com/Pis0sion/rblog/internal/rblog/cfg"
	"github.com/Pis0sion/rblog/internal/rblog/dto"
	"github.com/Pis0sion/rblog/internal/rblog/dto/mysql"
	"github.com/Pis0sion/rblog/internal/rblog/opts"
	"github.com/Pis0sion/rblog/internal/rblog/route"
	"github.com/Pis0sion/rblog/internal/rblog/route/api/v1/cache"
	"github.com/Pis0sion/rblog/pkg/db"
	"google.golang.org/grpc"
)

type ApplicationServe struct {
	grpcServe *grpcServe
	httpServe *serve.GenericServe
}

type PrepareApplicationServe struct {
	*ApplicationServe
}

type ExtraComponents struct {
	grpcOptions  *opts.GrpcOpts
	mysqlOptions *opts.MysqlOpts
	redisOptions *opts.RedisOpts
}

type CompleteExtraComponents struct {
	*ExtraComponents
}

func CreateApplicationServe(configure *cfg.Configure) (*ApplicationServe, error) {
	var grpcServe *grpcServe

	genericConfigure, err := buildGenericServe(configure)

	if err != nil {
		return nil, err
	}

	extraConfigure, _ := buildExtraComponents(configure)

	httpServe := genericConfigure.Complete().New()

	if grpcServe, err = extraConfigure.Complete().New(); err != nil {
		return nil, err
	}

	return &ApplicationServe{httpServe: httpServe, grpcServe: grpcServe}, nil
}

func (app *ApplicationServe) PrepareRun() *PrepareApplicationServe {

	route.InitializeRouters(app.httpServe.Engine)
	return &PrepareApplicationServe{app}
}

func (serve *PrepareApplicationServe) Run() error {
	serve.grpcServe.Run()
	return serve.httpServe.Run()
}

func (e *ExtraComponents) Complete() *CompleteExtraComponents {

	if e.grpcOptions.Address == "" {
		e.grpcOptions.Address = "127.0.0.1"
	}

	if e.grpcOptions.Port == 0 {
		e.grpcOptions.Port = 8081
	}

	return &CompleteExtraComponents{e}
}

func (c *CompleteExtraComponents) New() (*grpcServe, error) {
	// initialize mysql
	agent, err := mysql.GetDatabaseFactoryEntity(c.mysqlOptions)
	if err != nil {
		return nil, err
	}
	dto.SetClient(agent)

	// initialize redis
	initializeRedis(c.redisOptions)

	// initialize grpc
	gServe := grpc.NewServer()

	// pb
	pb.RegisterCacheServer(gServe, cache.NewCache())


	return &grpcServe{
		address: fmt.Sprintf("%s:%d", c.grpcOptions.Address, c.grpcOptions.Port),
		Server:  gServe,
	}, nil
}

func buildGenericServe(configure *cfg.Configure) (genericConfigure *serve.GenericConfigure, err error) {
	genericConfigure = serve.NewGenericConfigure()

	if err = configure.SrvOpts.ApplyTo(genericConfigure); err != nil {
		return nil, err
	}

	return
}

func buildExtraComponents(configure *cfg.Configure) (extraConfigure *ExtraComponents, err error) {
	return &ExtraComponents{
		grpcOptions:  configure.GrpcOpts,
		mysqlOptions: configure.MysqlOpts,
		redisOptions: configure.RedisOpts,
	}, nil
}

func initializeRedis(opts *opts.RedisOpts) {

	redisOptions := &db.RedisOptions{
		Host:          opts.Host,
		Port:          opts.Port,
		Address:       opts.Address,
		Username:      opts.Username,
		Password:      opts.Password,
		Database:      opts.Database,
		MasterName:    opts.MasterName,
		MinIdleConns:  opts.MinIdleConns,
		MaxActive:     opts.MaxActive,
		EnableCluster: opts.EnableCluster,
		Timeout:       opts.Timeout,
	}

	go db.Connect2Redis(context.Background(), redisOptions)
}
