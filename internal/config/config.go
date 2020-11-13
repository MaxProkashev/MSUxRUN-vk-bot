package config

import (
	"io/ioutil"
	"os"

	"github.com/MaxProkashev/MSUxRUN-vk-bot/internal/logs"
	_ "github.com/lib/pq" // pq driver for database/sql
	"gopkg.in/yaml.v3"
)

const (
	dirPr = "./config.yaml" // directory of the project configuration
)

// Config project template
type Config struct {
	Token string `yaml:"token"`
	DbURL string `yaml:"db_url"`

	CallH int `yaml:"call_h"`
	CallM int `yaml:"call_m"`
	CallS int `yaml:"call_s"`

	// preset message
	MessageSignUp   string `yaml:"message_sign_up"`
	MessageSignOut  string `yaml:"message_sign_out"`
	MessageDefault  string `yaml:"message_default"`
	MessageNonTrain string `yaml:"message_non_train"`
	MessageSpecFor0 string `yaml:"message_spec_for_0"`

	MessageNotice []string `yaml:"message_notice"`

	// keyboard
	CountTrain   int `yaml:"count_train"` // count training
	MainKeyboard []struct {
		Label  string `yaml:"label"`
		Coach  string `yaml:"coach"`
		NotDay string `yaml:"not_day"`
	} `yaml:"main_keyboard"`
	Schedule string `yaml:"schedule"` // расписание
	SchPhoto string `yaml:"sch_photo"`
	MyTrain  string `yaml:"my_train"`
}

// error unnamed func
var errLog = func(err error) {
	logs.Err("%s", err)
	os.Exit(1)
}

// GetProjectConfig from dir config
// if err exit with status 1
func GetProjectConfig() (c *Config) {
	// read file
	file, err := ioutil.ReadFile(dirPr)
	if err != nil {
		errLog(err)
	}
	// unmarshal file
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		errLog(err)
	}

	logs.Succes("get project configuration from %s", dirPr)

	return c
}
