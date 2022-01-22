package cfg

import (
	"github.com/Pis0sion/rblog/internal/rblog/opts"
)

var ApplicationParameters *opts.AppOpts

type Configure struct {
	*opts.Opts
}

func InitApplicationConfigure(options *opts.Opts) *Configure {

	ApplicationParameters = options.AppOpts
	return &Configure{options}
}
