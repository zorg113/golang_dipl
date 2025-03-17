package app

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/config"
)

type AntiBruteForceApp struct {
	router           interface{}
	grpcBlackList    interface{}
	grpcWhiteList    interface{}
	grpcBucket       interface{}
	grpcAutorization interface{}
	comLineInterface interface{}
	clintDb          interface{}
	logger           *zerolog.Logger
	config           *config.Config
}

func NewAntiBruteForceApp(config *config.Config) *AntiBruteForceApp {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	return &AntiBruteForceApp{logger: &logger, config: config}
}

func (a *AntiBruteForceApp) StartAppApi() {

}
