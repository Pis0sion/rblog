package rblog

import (
	"context"
	"github.com/Pis0sion/rblog/internal/pkg/serve"
	"github.com/Pis0sion/rblog/internal/rblog/cfg"
	"github.com/Pis0sion/rblog/internal/rblog/dto"
	"github.com/Pis0sion/rblog/internal/rblog/dto/mysql"
	"github.com/Pis0sion/rblog/internal/rblog/opts"
	"github.com/Pis0sion/rblog/internal/rblog/route"
	"github.com/Pis0sion/rblog/pkg/db"
)

type ApplicationServe struct {
	httpServe *serve.GenericServe
}

type PrepareApplicationServe struct {
	*ApplicationServe
}

type ExtraComponents struct {
	mysqlOptions *opts.MysqlOpts
	redisOptions *opts.RedisOpts
}

type CompleteExtraComponents struct {
	*ExtraComponents
}

func CreateApplicationServe(configure *cfg.Configure) (*ApplicationServe, error) {

	genericConfigure, err := buildGenericServe(configure)

	if err != nil {
		return nil, err
	}

	extraConfigure, _ := buildExtraComponents(configure)

	httpServe := genericConfigure.Complete().New()

	if err = extraConfigure.Complete().New(); err != nil {
		return nil, err
	}

	return &ApplicationServe{httpServe: httpServe}, nil
}

func (app *ApplicationServe) PrepareRun() *PrepareApplicationServe {

	route.InitializeRouters(app.httpServe.Engine)

	return &PrepareApplicationServe{app}
}

func (serve *PrepareApplicationServe) Run() error {
	return serve.httpServe.Run()
}

func (e *ExtraComponents) Complete() *CompleteExtraComponents {
	return &CompleteExtraComponents{e}
}

func (c *CompleteExtraComponents) New() error {
	// initialize mysql
	agent, err := mysql.GetDatabaseFactoryEntity(c.mysqlOptions)
	if err != nil {
		return err
	}
	dto.SetClient(agent)

	// initialize redis
	initializeRedis(c.redisOptions)

	return nil
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
