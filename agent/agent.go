package agent

import (
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"github.com/zhshch2002/goreq"
	"net/http"
	"time"
)

func init() {
	viper.SetDefault("server", "https://sender.xzhshch.com")
}

type Signaling struct {
	ID        string `json:"id"`
	Signaling string `json:"signaling"`
	Result    string `json:"result"`
}

type Agent struct {
	AccessKey, Name string
}

func (a *Agent) Run(fn func([]Signaling)) {
	data, err := goreq.Get(viper.GetString("server")+"/api/v1/ispro").
		AddParam("accesskey", a.AccessKey).
		Do().JSON()
	var pro bool
	if err == nil {
		pro = data.Get("res").Bool()
	}
	log.Err(err).Bool("pro", pro).Msg("get ispro")
	c := cron.New()
	fun := func() {
		signalings, err := a.PollingSignalings()
		if err != nil {
			log.Err(err).Msg("Agent PollingSignalings error")
		}
		fn(signalings)
	}
	if pro {
		_, err := c.AddFunc("@every 1m", fun)
		if err != nil {
			panic(err)
		}
	} else {
		_, err := c.AddFunc("@every 10m", fun)
		if err != nil {
			panic(err)
		}
	}
	c.Run()
}

func (a *Agent) PollingSignalings() ([]Signaling, error) {
	var res []Signaling
	resp, data, err := goreq.Get(viper.GetString("server")+"/api/v1/signaling").
		AddParam("accesskey", a.AccessKey).
		AddParam("agent", a.Name).
		Do().RespAndJSON()
	if err != nil {
		log.Err(err).Msg("polling signaling error")
		return nil, err
	}
	if resp.StatusCode != http.StatusOK || data.Get("code").Int() != 0 {
		log.Err(err).Int("status", resp.StatusCode).Str("resp", resp.Text).Interface("resp-header", resp.Header).Msg("polling signaling error")
		return nil, err
	}

	log.Info().Interface("signaling", data.Get("res").Value()).Msg(fmt.Sprint("get signaling"))

	data.Get("res").ForEach(func(k, v gjson.Result) bool {
		signaling := v.Get("signaling").String()
		if signaling != "" {
			res = append(res, Signaling{
				ID:        v.Get("id").String(),
				Signaling: v.Get("signaling").String(),
			})
		}
		return true
	})
	return res, nil
}

var SubmitMaxRetry = 2

func (a *Agent) SubmitSignalingResult(s Signaling) error {
	retry := 0
W:
	resp := goreq.Post(viper.GetString("server")+"/api/v1/result").
		AddParam("accesskey", a.AccessKey).
		AddParam("id", s.ID).
		SetMultipartBody(goreq.FormField{
			Name:  "result",
			Value: s.Result,
		}).
		Do()
	res, err := resp.JSON()
	if retry < SubmitMaxRetry && (err != nil || !res.Get("code").Exists()) {
		log.Err(resp.Err).
			Int("retry", retry).
			Str("id", s.ID).
			Str("response", resp.Text).
			Msg("submit result error retry")
		retry += 1
		time.Sleep(5 * time.Second)
		goto W
	}
	log.Err(resp.Err).
		Int("retry", retry).
		Str("id", s.ID).
		Str("result", s.Result).
		Str("response", resp.Text).
		Msg("submit result")
	return err
}

func SendMessage(accesskey, title, desp, link, channel string) error {
	data, err := goreq.Post(viper.GetString("server")+"/api/v1/message").
		AddParam("accesskey", accesskey).
		AddParam("title", title).
		AddParam("desp", desp).
		AddParam("link", link).
		AddParam("channel", channel).
		Do().JSON()
	if err == nil {
		if !data.Get("code").Exists() || data.Get("code").Int() != 0 {
			return errors.New(data.Get("msg").String())
		}
	}
	log.Err(err).
		Str("title", title).
		Msg("send message")
	return err
}
