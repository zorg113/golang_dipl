package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/app"
	"github.com/zorg113/golang_dipl/atibruteforce/internal/config"
)

func main() {
	confParam := flag.String("conf", "path to configuration file", "a string")
	flag.Parse()
	if *confParam == "" {
		fmt.Println("No config file provided")
		return
	}
	fmt.Println("Init Config from file")
	conf, err := config.NewConfig(*confParam)
	if err != nil {
		fmt.Println("cant't initialize config")
		return
	}
	if key := os.Getenv("ADMIN_API_KEY"); key != "" {
		conf.Admin.APIKey = key
	}

	if conf.Admin.APIKey == "" {
		fmt.Println("ADMIN_API_KEY is not set")
	}

	fmt.Println("init logger")
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	app := app.NewAntiBruteForceApp(&logger, &conf)
	app.StartAppAPI()
}
