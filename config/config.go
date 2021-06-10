package config

import (
	"os"
)

type Action struct {
	Name    string `yaml:"name,omitempty"`
	Cmd     string `yaml:"cmd"`
	Script  string `yaml:"script"`
	Forward string `yaml:"forward"`
	Dir     string `yaml:"dir"`
}

var Actions map[string]*Action
var Wd string

func init() {
	var err error
	Wd, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}
