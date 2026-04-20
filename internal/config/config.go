package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Listen struct {
	Type   string `yaml:"type"`
	BindIP string `yaml:"bindIp"`
	Port   string `yaml:"port"`
}

type Server struct {
	ReadTimeout  int    `yaml:"readTimeout"`
	WriteTimeout int    `yaml:"writeTimeout"`
	IdleTimeout  int    `yaml:"idleTimeout"`
	ServerType   string `yaml:"serverType"`
}

type AppConfig struct {
	LogLevel string `yaml:"logLevel"`
}

type DBData struct {
	Host       string `yaml:"host"`
	DBName     string `yaml:"dbName"`
	Port       string `yaml:"port"`
	DBUser     string `yaml:"dbUser"`
	DBPassword string `yaml:"dbPassword"`
	SslMode    string `yaml:"sslMode"`
}

type Bucket struct {
	IPLimit             int `yaml:"ipLimit"`
	LoginLimit          int `yaml:"loginLimit"`
	PasswordLimit       int `yaml:"passwordLimit"`
	ResetBucketInterval int `yaml:"resetBucketInterval"`
}

type AdminConf struct {
	APIKey string `yaml:"apiKey"`
}

type Config struct {
	Listen    Listen
	Server    Server
	AppConfig AppConfig
	DBData    DBData
	Bucket    Bucket
	Admin     AdminConf
}

func NewConfig(path string) (Config, error) {
	var conf Config
	file, err := os.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("can't open config file: %s", err.Error())
		return conf, err
	}
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		err = fmt.Errorf("can't unmarshall  config file: %s", err.Error())
		return conf, err
	}
	return conf, nil
}
