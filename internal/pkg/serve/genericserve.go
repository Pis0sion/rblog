package serve

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// GenericServe generic api server provider configure
// configure struct
type GenericServe struct {
	*gin.Engine
	mode                   string
	InsecureServeConfigure *InsecureServeConfigure
	insecureServe          *http.Server
}

func initializeGenericServe(serve *GenericServe) {

	serve.setServeMode()
	serve.installMiddlewares()
}

func (s GenericServe) setServeMode() {

	gin.SetMode(s.mode)
}

func (s GenericServe) installMiddlewares() {

	s.Engine.Use(gin.Logger())

}

func (s GenericServe) Run() error {

	s.insecureServe = &http.Server{
		Addr:         s.InsecureServeConfigure.GetActualAddress(),
		Handler:      s,
		ReadTimeout:  s.InsecureServeConfigure.ReadTimeout,
		WriteTimeout: s.InsecureServeConfigure.WriteTimeout,
	}

	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGKILL)

	go func() {

		if err := s.insecureServe.ListenAndServe(); err != nil {
			exitSignal <- syscall.SIGINT
			log.Fatalln(err)
		}
	}()

	<-exitSignal
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.insecureServe.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
