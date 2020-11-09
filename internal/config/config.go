package config

import (
	"database/sql"
	"io/ioutil"
	"msuxrun-bot/internal/logs"
	"os"

	_ "github.com/lib/pq" // pq driver for database/sql
	"gopkg.in/yaml.v3"
)

const (
	dirConfig = "./config.yaml"
)

// Config project
type Config struct {
	Token   string `yaml:"token"`
	PORT    string `yaml:"std_port"`
	DBURL   string `yaml:"db_url"`
	GroupID int    `yaml:"group_id"`

	MainKyeboard []struct {
		Row []struct {
			Label   string `yaml:"label"`
			Payload string `yaml:"payload"`
			Color   string `yaml:"color"`
		} `yaml:"row"`
	} `yaml:"main_kyeboard"`

	DB *sql.DB
}

// OpenDB heroku
func (c *Config) OpenDB() Config {
	db, err := sql.Open("postgres", c.DBURL)
	if err != nil {
		logs.ErrorLogger.Printf("could`t open db")
		os.Exit(1)
	}
	logs.Succes("open db")
	c.DB = db
	return *c
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
