package main

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/pkg/errors"
)

type Config struct {
	Port string //`yaml:"port"`
	Commands []Command
}

type Command struct {
	Name string //`yaml:"name"`
	Cmd []string //`yaml:"cmd"`
}

func readConfig(path string) (*Config, error) {
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "ReadConfig/ReadFile")
	}
	var c Config
	err = yaml.Unmarshal(fileContent, &c)
	if err != nil {
		return nil, errors.Wrap(err, "ReadConfig/Unmarshal")
	}
	return &c, nil
}