package main

import (
	"io/ioutil"
	"strconv"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type config struct {
	// log level available: "disable", "fatal", "error", "warn", "info", "debug"
	LogLevel   string
	Port       uint16
	Commands   []*command
	commandMap map[string]*command
}

func (c *config) getCommand(name string) (command, error) {
	cmd, ok := c.commandMap[name]
	if !ok {
		return command{}, errors.Errorf("command %+v not found", name)
	}
	return *cmd, nil
}

func (c *config) getListenAddress() string {
	return ":" + strconv.FormatInt(int64(c.Port), 10)
}

type command struct {
	Name  string
	Cmd   []string
	Stdin string
	Long  bool
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

	c.commandMap = make(map[string]*command)
	for _, cmd := range c.Commands {
		c.commandMap[cmd.Name] = cmd
	}

	return &c, nil
}
