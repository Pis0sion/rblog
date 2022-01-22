package opts

import (
	"github.com/Pis0sion/rblog/internal/pkg/serve"
	"time"
)

type SrvOpts struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func (o SrvOpts) ApplyTo(configure *serve.GenericConfigure) error {

	configure.Mode = o.RunMode
	configure.InsecureServeConfigure = &serve.InsecureServeConfigure{
		Address:      o.HttpPort,
		ReadTimeout:  o.ReadTimeout,
		WriteTimeout: o.WriteTimeout,
	}

	return nil
}
