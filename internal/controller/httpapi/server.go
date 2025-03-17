package httpapi

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/config"
)

type HttpApiServer struct {
	server  *http.Server
	handler http.Handler
	config  *config.Config
	logger  *zerolog.Logger
}

func NewHttpApiServer(handler http.Handler, config *config.Config, log *zerolog.Logger) *HttpApiServer {
	return &HttpApiServer{
		config:  config,
		handler: handler,
		logger:  log,
	}
}

func (srv *HttpApiServer) Start() error {
	srv.server = &http.Server{
		Addr:         srv.config.Listen.BindIP + ":" + srv.config.Listen.Port,
		Handler:      srv.handler,
		ReadTimeout:  time.Duration(srv.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(srv.config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(srv.config.Server.IdleTimeout) * time.Second,
	}
	srv.logger.Info().Msg("start http server")
	err := srv.server.ListenAndServe()
	return err
}

func (srv *HttpApiServer) ShutdowService(c chan os.Signal) {
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	sig := <-c
	srv.logger.Info().Msg("service is stopped by signal " + sig.String())
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := srv.server.Shutdown(ctx); err != nil {
		srv.logger.Error().Msg(err.Error())
	}
}
