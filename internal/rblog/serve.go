package rblog

import (
	"github.com/Pis0sion/rblog/internal/pkg/serve"
	"github.com/Pis0sion/rblog/internal/rblog/cfg"
	"github.com/Pis0sion/rblog/internal/rblog/dto"
	"github.com/Pis0sion/rblog/internal/rblog/dto/mysql"
	"github.com/Pis0sion/rblog/internal/rblog/opts"
	"github.com/Pis0sion/rblog/internal/rblog/route"
)

type ApplicationServe struct {
	httpServe *serve.GenericServe
}

type PrepareApplicationServe struct {
	*ApplicationServe
}

type ExtraComponents struct {
	mysqlOptions *opts.MysqlOpts
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

	agent, err := mysql.GetDatabaseFactoryEntity(c.mysqlOptions)

	if err != nil {
		return err
	}

	dto.SetClient(agent)

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
	return &ExtraComponents{mysqlOptions: configure.MysqlOpts}, nil
}
