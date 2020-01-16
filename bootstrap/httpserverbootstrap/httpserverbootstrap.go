package httpserverbootstrap

import (
	"context"
	"errors"
	"github.com/zibilal/logwrapper"
	"net/http"
	"os"
	"os/signal"
	"simplims/api"
	v3 "simplims/api/v3"
	"simplims/appctx"
	"simplims/bootstrap/wrapper/gingonic"
	"time"
)

type HttpServerBootstrap struct {
	apiEngine api.ApiEngine
	address   string

	services      map[string]interface{}
	serverContext *appctx.AppContext
}

func NewHttpServerBootstrap() *HttpServerBootstrap {
	httpServer := new(HttpServerBootstrap)
	return httpServer
}

func (s *HttpServerBootstrap) registerVersions() error {
	err := s.apiEngine.RegisterVersion(v3.VersionOneApi())
	if err != nil {
		return err
	}

	return nil
}

func (s *HttpServerBootstrap) Init() error {
	s.serverContext = appctx.NewAppContext()

	file, err := os.Open(appctx.DefaultConfigFlagVal)
	if err != nil {
		return err
	}

	err = s.serverContext.LoadAppContext(file)

	if err != nil {
		return err
	}

	return nil
}

func (s *HttpServerBootstrap) Run() error {

	if s.serverContext == nil {
		return errors.New("bootstrap have not calling Init yet")
	}

	go func() {
		logwrapper.Info("Address: ", s.serverContext.Config.Address)
		s.apiEngine = gingonic.NewGonicEngine(s.serverContext.Config.Address)
		err := s.registerVersions()
		if err != nil {
			panic(err)
		}
		if err := s.apiEngine.Execute(); err != nil && err != http.ErrServerClosed {
			logwrapper.Fatal("unable to initiate server due to ", err.Error())
		}
	}()

	// wait for interrupt signal to gracefully shutdown the server
	// with timeout of 5 seconds
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// waits for 5 seconds
	parentContext := context.Background()
	ctx, cancel := context.WithTimeout(parentContext, 5*time.Second)
	defer cancel()

	if err := s.apiEngine.Shutdown(ctx); err != nil {
		logwrapper.Fatal("Server shutdown: ", err)
	}

	logwrapper.Info("Server exiting")

	return nil
}
