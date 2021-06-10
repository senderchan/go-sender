package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/senderchan/go-sender/agent"
	"github.com/senderchan/go-sender/config"
	"github.com/senderchan/go-sender/executor"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	version = "Dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
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

	err = viper.UnmarshalKey("action", &config.Actions)
	if err != nil {
		panic(err)
	}
	for n := range config.Actions {
		config.Actions[n].Name = n
		if config.Actions[n].Dir == "" {
			config.Actions[n].Dir = "."
		}
		config.Actions[n].Dir = filepath.Join(config.Wd, config.Actions[n].Dir)
		if config.Actions[n].Script != "" {
			config.Actions[n].Script = filepath.Join(config.Wd, config.Actions[n].Script)
		}
	}

	if !viper.GetBool("log-json") {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: viper.GetBool("no-color")})
	}


	if pflag.Arg(0) == "version" {
		fmt.Print("v")
		fmt.Println(version)
		fmt.Println(commit)
		fmt.Println(date)
		return
	} else if pflag.Arg(0) == "service" {
		_ = ioutil.WriteFile("/etc/systemd/system/sender.service", []byte(`[Unit]
Description=Sender
Documentation=https://github.com/senderchan/go-sender
After=network.target network-online.target
Requires=network-online.target

[Service]
User=root
Group=root
ExecStart=/usr/bin/sender
Restart=on-success
TimeoutStopSec=5s
ProtectSystem=full

[Install]
WantedBy=multi-user.target`), 0644)
		exec.Command("systemctl","daemon-reload")
		return
	}

	if viper.GetString("name") == "" || viper.GetString("accesskey") == "" {
		log.Panic().Str("name", viper.GetString("name")).Str("accesskey", viper.GetString("accesskey")).Msg("accesskey or agent name is empty")
	}

	log.Info().Str("version", version).Msg("sender agent start")
	a := &agent.Agent{
		AccessKey: viper.GetString("accesskey"),
		Name:      viper.GetString("name"),
	}
	a.Run(func(signalings []agent.Signaling) {
		for _, s := range signalings {
			cfg, ok := config.Actions[strings.Split(s.Signaling, " ")[0]]
			if !ok {
				cfg, ok = config.Actions["default"]
			}
			if !ok {
				err := errors.New("unknown action")
				log.Err(err).Str("signaling", s.Signaling).Send()
				s.Result = fmt.Sprint("agent error: ", err)
				go a.SubmitSignalingResult(s)
			}
			s.Result = string(executor.ExecuteAction(s.Signaling, cfg))
			go a.SubmitSignalingResult(s)
		}
	})
}
