package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Listen struct {
	Type   string `yaml:"Type"`
	BindIP string `yaml:"BindIP"`
	Port   string `yaml:"Port"`
}

type Server struct {
	ReadTimeout  int    `yaml:"ReadTimeout"`
	WriteTimeout int    `yaml:"WriteTimeout"`
	IdleTimeout  int    `yaml:"IdleTimeout"`
	ServerType   string `yaml:"ServerType"`
}

type AppConfig struct {
	LogLevel string `yaml:"LogLevel"`
}

type DbData struct {
	Host       string `yaml:"Host"`
	DbName     string `yaml:"DbName"`
	Port       string `yaml:"Port"`
	DbUser     string `yaml:"DbUser"`
	DbPassword string `yaml:"DbPassword"`
	SslMode    string `yaml:"SslMode"`
}
type Bucket struct {
	IpLimit             int `yaml:"IpLimit"`
	LoginLimit          int `yaml:"LoginLimit"`
	PasswordLimit       int `yaml:"PasswordLimit"`
	ResetBucketInterval int `yaml:"ResetBucketInterval"`
}
type Config struct {
	Listen    Listen
	Server    Server
	AppConfig AppConfig
	DbData    DbData
	Bucket    Bucket
}

func WriteData(conf Config) {
	yamlFile, err := yaml.Marshal(&conf)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(yamlFile))

	f, err := os.Create("config.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(yamlFile)
	if err != nil {
		panic(err)
	}
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
