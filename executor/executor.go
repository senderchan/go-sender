package executor

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/senderchan/go-sender/config"
	"github.com/senderchan/go-sender/tool"
	"github.com/spf13/viper"
	"github.com/zhshch2002/goreq"
	"os/exec"
	"regexp"
	"strconv"
)

var re = regexp.MustCompile(`%[0-9]+`)

func ExecuteAction(signaling string, cfg *config.Action) (res []byte) {
	defer func() {
		if err := recover(); err != nil {
			log.Error().Str("err", fmt.Sprint(err)).Msg("exec action panic")
			if len(res) == 0 {
				res = []byte(fmt.Sprint("exec action panic: ", err))
			}
		}
	}()

	name, arg := tool.ParseSignaling(signaling)
	log.Info().Str("signaling", signaling).Str("name", name).Interface("arg", arg).Msg("ExecuteAction")

	if cfg.Forward != "" {
		go func() {
			res := goreq.Get(cfg.Forward).AddParam("signaling", signaling).Do()
			log.Err(res.Err).Str("response", res.Text).Str("signaling", signaling).Msg("action forward")
		}()
	}

	if cfg.Cmd != "" {
		cmd := cfg.Cmd
		cmd = re.ReplaceAllStringFunc(cmd, func(s string) string {
			if len([]rune(s)) <= 1 {
				panic(errors.New("error placeholder " + s))
			}
			i, err := strconv.Atoi(s[1:])
			if err != nil {
				panic(err)
			}
			if i > len(arg) {
				panic(errors.New("parameter index over length"))
			}
			if i == 0 {
				return name
			}
			return arg[i-1]
		})

		c := exec.Command(viper.GetString("shell"), "-c", cmd)
		c.Dir = cfg.Dir
		var err error
		res, err = c.CombinedOutput()
		log.Err(err).
			Str("output", string(res)).
			Interface("signaling", signaling).
			Interface("config", cfg).
			Str("cmd", cmd).
			Msg("exec action cmd finish")
	} else if cfg.Script != "" {
		c := exec.Command(viper.GetString("shell"), append([]string{cfg.Script}, arg...)...)
		var err error
		res, err = c.CombinedOutput()
		log.Err(err).
			Str("output", string(res)).
			Interface("signaling", signaling).
			Interface("config", cfg).
			Msg("exec action script finish")
	}
	return res
}
