package rblog

import (
	"github.com/Pis0sion/rblog/internal/rblog/cfg"
	"github.com/Pis0sion/rblog/internal/rblog/opts"
	"github.com/Pis0sion/rblog/lib/app"
	"log"
)

func NewApp() (application *app.App) {

	options, err := opts.NewOptions()

	if err != nil {
		log.Fatalln(err)
	}

	application = app.NewApplication("rblog",
		"rblog service",
		app.WithDescription("rblog description"),
		app.WithVersion("v1.0.0"),
		app.WithRunFunc(run(options)),
	)

	return
}

func run(opts *opts.Opts) app.RunFunc {

	return func(name string) error {

		configure := cfg.InitApplicationConfigure(opts)
		return Run(configure)
	}
}
