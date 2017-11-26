package main

import (
	"io/ioutil"
	"strconv"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type config struct {
	// log level available: "disable", "fatal", "error", "warn", "info", "debug"
	LogLevel string
	Port     uint16
	Commands []command
}

func (c *config) getListenAddress() string {
	return ":" + strconv.FormatInt(int64(c.Port), 10)
}

type command struct {
	Name string
	Cmd  []string
}

func readConfig(path string) (*config, error) {
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "ReadConfig/ReadFile")
	}
	var c config
	err = yaml.Unmarshal(fileContent, &c)
	if err != nil {
		return nil, errors.Wrap(err, "ReadConfig/Unmarshal")
	}
	return &c, nil
}
