package config

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
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
	pflag.String("config", "", "path to config file")
	pflag.Bool("no-color", false, "log output with color")
	pflag.Bool("log-json", false, "log output in json")

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		panic(err)
	}
	pflag.Parse()

	viper.SetDefault("server", "https://sender.xzhshch.com")
	viper.SetDefault("shell", os.Getenv("SHELL"))

	if viper.GetString("config") != "" {
		f, err := os.ReadFile(viper.GetString("config"))
		if err != nil {
			panic(err)
		}
		viper.SetConfigType("yaml")
		err = viper.ReadConfig(bytes.NewReader(f))
		if err != nil {
			panic(err)
		}
	} else {
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
		viper.AddConfigPath("/etc/sender/agent")
		viper.AddConfigPath("$HOME/.sender/agent")
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			_ = ioutil.WriteFile("./config.yaml", []byte(`name: test
accesskey: 758e2e934f**************9ebafcb

action:
  hello:
    cmd: echo "hello"
  ls:
    cmd: ls -la
    dir: ~/
  aaa:
    script: ./run.sh
    dir: ~/
  send:
    cmd: echo "ok"
    forward: http://localhost:9000
  default:
    cmd: docker version`), 0644)
			fmt.Println(err)
			fmt.Println("A configuration file has been generated.")
			os.Exit(-1)
		}
		if err != nil {
			panic(err)
		}
	}

	err = viper.UnmarshalKey("action", &Actions)
	if err != nil {
		panic(err)
	}
	for n := range Actions {
		Actions[n].Name = n
		if Actions[n].Dir == "" {
			Actions[n].Dir = "."
		}
		Actions[n].Dir = filepath.Join(Wd, Actions[n].Dir)
		if Actions[n].Script != "" {
			Actions[n].Script = filepath.Join(Wd, Actions[n].Script)
		}
	}

	if !viper.GetBool("log-json") {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: viper.GetBool("no-color")})
	}

	Wd, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}
