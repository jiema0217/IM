package config

import (
	"IMProject/pkg/mysql"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

var Cfg BaseConfig

type BaseConfig struct {
	ServiceName string                       `yaml:"service_name"`
	LogLevel    int                          `yaml:"log_level"`
	RsaAK       string                       `yaml:"rsa_ak"`
	Port        int                          `yaml:"port"`
	Mysql       map[string]mysql.MySqlConfig `yaml:"mysql"`
}

func GetConfig(serverName string) {
	env := os.Getenv("im_env")
	if env == "" {
		panic("env is empty")
	}
	cfgFile, err := os.Open(fmt.Sprintf("./config/%s_%s.yaml", serverName, env))
	if err != nil {
		panic(err)
	}
	defer cfgFile.Close()
	decoder := yaml.NewDecoder(cfgFile)
	err = decoder.Decode(&Cfg)
	if err != nil {
		panic(err)
	}
}
