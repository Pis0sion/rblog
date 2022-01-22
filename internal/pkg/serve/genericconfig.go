package serve

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type GenericConfigure struct {
	Mode                   string
	InsecureServeConfigure *InsecureServeConfigure
}

type CompleteGenericConfigure struct {
	*GenericConfigure
}

type InsecureServeConfigure struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func (insecure InsecureServeConfigure) GetActualAddress() string {
	return fmt.Sprintf(":%s", insecure.Address)
}

func NewGenericConfigure() *GenericConfigure {
	return &GenericConfigure{
		Mode: gin.DebugMode,
	}
}

func (c *GenericConfigure) Complete() *CompleteGenericConfigure {
	return &CompleteGenericConfigure{c}
}

func (c CompleteGenericConfigure) New() *GenericServe {

	serve := &GenericServe{
		Engine:                 gin.New(),
		mode:                   c.Mode,
		InsecureServeConfigure: c.InsecureServeConfigure,
	}

	initializeGenericServe(serve)

	return serve
}
