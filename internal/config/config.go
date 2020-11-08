package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

const (
	dirConfig  = "./config.yaml"
	packageErr = "ERROR internal/config/config.go::" // package standart error
)

// Config project
type Config struct {
	Token   string `yaml:"token"`
	GroupID int    `yaml:"group_id"`
}

// StartLog to init bot
func (c *Config) StartLog() {
	log.Printf("token = %s", c.Token)
	log.Printf("group_id = %d", c.GroupID)
}

type packErr struct {
	funcName string
	comment  string
	err      error
}

var (
	stdErr packErr
)

func (p *packErr) addComment(comm string) {
	p.comment = comm
}

func errInit(pe packErr) {
	log.Fatalf("%s%s ==> %s <== COMMENT: %s",
		packageErr,
		pe.funcName,
		pe.err.Error(),
		pe.comment,
	)
}

// GetProjectConfig from dir config
// if err exit with status 1
func GetProjectConfig() (config Config) {
	stdFuncErr := func(err error) packErr {
		return packErr{
			funcName: "GetProjectConfig",
			comment:  "no comment",
			err:      err,
		}
	}
	// read file
	file, err := ioutil.ReadFile(dirConfig)
	if err != nil {
		stdErr = stdFuncErr(err)
		stdErr.addComment("err read file")
		errInit(stdErr)
	}
	// unmarshal file
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		stdErr = stdFuncErr(err)
		stdErr.addComment("err unmarshal file")
		errInit(stdErr)
	}

	return config
}
