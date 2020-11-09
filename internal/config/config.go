package config

import (
	"io/ioutil"
	"msuxrun-bot/internal/logs"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	dirConfig = "./config.yaml"
)

// Config project
type Config struct {
	Token   string `yaml:"token"`
	GroupID int    `yaml:"group_id"`
}

func errLog(err error) {
	logs.ErrorLogger.Printf("%s", err)
	os.Exit(1)
}

// GetProjectConfig from dir config
// if err exit with status 1
func GetProjectConfig() (config Config) {
	// read file
	file, err := ioutil.ReadFile(dirConfig)
	if err != nil {
		errLog(err)
	}
	// unmarshal file
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		errLog(err)
	}

	logs.Succes("get config")
	return config
}
