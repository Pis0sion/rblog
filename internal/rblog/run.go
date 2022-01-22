package rblog

import "github.com/Pis0sion/rblog/internal/rblog/cfg"

func Run(configure *cfg.Configure) error {

	applicationServe, err := CreateApplicationServe(configure)

	if err != nil {
		return err
	}

	return applicationServe.PrepareRun().Run()
}
