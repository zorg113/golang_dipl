package app

import (
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/config"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/cli"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/grpcapi"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/httpapi"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/controller/httpapi/handlers"
	"github.com/zorg113/golang_dipl/atibruteforce/model/service"
	"github.com/zorg113/golang_dipl/atibruteforce/store/adapters"
	"github.com/zorg113/golang_dipl/atibruteforce/store/client"
)

type AntiBruteForceApp struct {
	router                  *httpapi.HttpApiRouter
	grpcBlackListServer     *grpcapi.BlackListServer
	grpcWhiteListServer     *grpcapi.WhiteListServer
	grpcBucketServer        *grpcapi.BucketServer
	grpcAuthorizationServer *grpcapi.AuthorizationServer
	cli                     *cli.CommandLineInterface
	dbClient                *client.PostgresSql
	logger                  *zerolog.Logger
	config                  *config.Config
}

func NewAntiBruteForceApp(logger *zerolog.Logger, config *config.Config) *AntiBruteForceApp {
	dbClient := client.NewPostgresSql(logger, config)
	err := dbClient.Open()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to postgres")
	}
	blackListStor := adapters.NewBlackListStorage(dbClient)
	blackListService := service.NewBlackList(blackListStor, logger)
	blackListHandler := handlers.NewBlackList(blackListService, logger)
	blacklistGrpc := grpcapi.NewBlackListServer(blackListService, logger)

	whiteListStor := adapters.NewWhiteListStorage(dbClient)
	whiteListService := service.NewWhiteList(whiteListStor, logger)
	whiteListHandler := handlers.NewWhiteList(whiteListService, logger)
	whiteListGrpc := grpcapi.NewWhiteListServer(whiteListService, logger)

	authorizarionService := service.NewAuthorization(blackListService, whiteListService, config, logger)
	authrizationHandler := handlers.NewAuthorization(authorizarionService, logger)
	bucketHandler := handlers.NewBucket(authorizarionService, logger)
	bucketGrpc := grpcapi.NewBucketServer(authorizarionService, logger)
	authorizationGrpc := grpcapi.NewAuthorization(authorizarionService, logger)

	router := httpapi.NewRouter(authrizationHandler, blackListHandler, whiteListHandler, bucketHandler, logger)

	cli := cli.NewCommandLineInterface(authorizarionService, blackListService, whiteListService)

	return &AntiBruteForceApp{
		router:                  router,
		grpcBlackListServer:     blacklistGrpc,
		grpcWhiteListServer:     whiteListGrpc,
		grpcBucketServer:        bucketGrpc,
		grpcAuthorizationServer: authorizationGrpc,
		cli:                     cli,
		dbClient:                dbClient,
		logger:                  logger,
		config:                  config,
	}
}

func (a *AntiBruteForceApp) StartAppApi() {
	c := make(chan os.Signal, 1)
	go a.cli.Run(c)
	switch a.config.Server.ServerType {
	case "grpc":
		a.logger.Info().Msg("Init grpc server")
		grpcServer := grpcapi.NewServerGRPC(a.grpcAuthorizationServer,
			a.grpcBlackListServer,
			a.grpcWhiteListServer,
			a.grpcBucketServer,
			a.config,
			a.logger)
		go grpcServer.Shutdown(c)
		err := grpcServer.Start()
		if err != nil {
			a.logger.Fatal().Err(err).Msg("failed to start grpc server")
		}
		err = a.dbClient.Close()
		if err != nil {
			a.logger.Fatal().Err(err).Msg("failed to close db connection")
		}
	case "http":
		a.logger.Info().Msg("Init http server")
		a.router.InitRouters()

		server := httpapi.NewHttpApiServer(a.router.GetRouter(), a.config, a.logger)
		go server.ShutdowService(c)
		err := server.Start()
		if err != nil {
			if err == http.ErrServerClosed {
				a.logger.Error().Err(err)
				err = a.dbClient.Close()
				if err != nil {
					a.logger.Error().Err(err).Msg("failed to close db connection")
				}
				return
			}
			a.logger.Fatal().Err(err).Msg("failed to start http server")
		}

	}
}
