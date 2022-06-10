package beacon

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"pkg.oars.vip/go-pkg/req"
)

type Config struct {
	AppId      string
	Secret     string
	BeaconAddr string
}

type Client struct {
	cfg *Config
}

func New(cfg *Config) *Client {
	return &Client{cfg}
}

type MsgBody struct {
	ToType      string      `json:"toType"`
	MessageType string      `json:"messageType"`
	Text        interface{} `json:"text"`
	Textcard    TextcardMsg `json:"textcard"`
	Receivers   []string    `json:"receivers"`
}

type TextMsg struct {
	Content string `json:"content"`
}

type TextcardMsg struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Btntxt      string `json:"btntxt"`
}

type resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (c *Client) Send(data *MsgBody) error {
	var res resp
	token := base64.StdEncoding.EncodeToString([]byte(c.cfg.AppId + ":" + c.cfg.Secret))
	err := req.ReqJSON(http.MethodPost, c.cfg.BeaconAddr+"/beacon/app/api/v1/messages", data, &res, map[string]string{"Authorization": "App " + token})
	if err != nil {
		return err
	}
	if res.Code != 10000 {
		return fmt.Errorf("error code: %d, Msg: %s", res.Code, res.Msg)
	}
	return nil
}
