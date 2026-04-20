package grpcapi

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/config"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/grpcapi/authorizationpb"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/grpcapi/blacklistpb"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/grpcapi/bucketpb"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/grpcapi/interceptor"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/grpcapi/whitelistpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ServerGRPC struct {
	authorizationServer authorizationpb.AuthorizationServer
	blacklistServer     blacklistpb.BlackListServer
	whitelistServer     whitelistpb.WhiteListServer
	bucketServer        bucketpb.BucketServer
	grpcServer          *grpc.Server
	config              *config.Config
	log                 *zerolog.Logger
}

func NewServerGRPC(authorizationServer authorizationpb.AuthorizationServer, blacklistServer blacklistpb.BlackListServer, whitelistServer whitelistpb.WhiteListServer, bucketServer bucketpb.BucketServer, config *config.Config, log *zerolog.Logger) *ServerGRPC { //nolint:lll
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			interceptor.AdminAuth(config.Admin.APIKey, log),
		),
		grpc.StreamInterceptor(
			interceptor.StreamAdminAuth(config.Admin.APIKey, log),
		),
	)
	return &ServerGRPC{
		authorizationServer: authorizationServer,
		blacklistServer:     blacklistServer,
		whitelistServer:     whitelistServer,
		bucketServer:        bucketServer,
		grpcServer:          grpcServer,
		config:              config,
		log:                 log,
	}
}

func (s *ServerGRPC) Start() error {
	s.log.Info().Msg("start grpc server")
	listener, err := net.Listen("tcp", s.config.Listen.BindIP+":"+s.config.Listen.Port)
	if err != nil {
		return err
	}
	authorizationpb.RegisterAuthorizationServer(s.grpcServer, s.authorizationServer)
	blacklistpb.RegisterBlackListServer(s.grpcServer, s.blacklistServer)
	whitelistpb.RegisterWhiteListServer(s.grpcServer, s.whitelistServer)
	bucketpb.RegisterBucketServer(s.grpcServer, s.bucketServer)
	if s.config.AppConfig.LogLevel == "debug" {
		reflection.Register(s.grpcServer)
		s.log.Warn().Msg("gRPC reflection enabled — disable in production")
	}
	err = s.grpcServer.Serve(listener)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServerGRPC) Shutdown(c chan os.Signal) {
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	sig := <-c
	s.log.Info().Msg("service is stopped by signal " + sig.String())
	s.grpcServer.GracefulStop()
}
