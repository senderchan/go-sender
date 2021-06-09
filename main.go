package main

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/senderchan/go-sender/agent"
	"github.com/senderchan/go-sender/config"
	"github.com/senderchan/go-sender/executor"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io/ioutil"
	"os/exec"
	"strings"
)

var (
	version = "Dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if pflag.Arg(0) == "version" {
		fmt.Print("v")
		fmt.Println(version)
		fmt.Println(commit)
		fmt.Println(date)
		return
	} else if pflag.Arg(0) == "service" {
		_ = ioutil.WriteFile("/etc/systemd/system/sender.service", []byte(`[Unit]
Description=Sender
Documentation=https://github.com/zhshch2002/sender-agent
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
