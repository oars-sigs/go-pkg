package idaas

import (
	"net/http"

	"pkg.oars.vip/go-pkg/req"
	"pkg.oars.vip/go-pkg/server/base"
)

type SysResp struct {
	*base.DataResponse
	Data map[string]interface{} `json:"data"`
}

func (c *Client) Sys(param, client string) (map[string]interface{}, error) {
	var res SysResp
	uri := "/idaas-app/sys?client=" + client
	if param != "" {
		uri += "&params=" + param
	}
	err := req.ReqJSON(http.MethodGet, c.GetUrl(uri), nil, &res, c.SetAuthHeader(nil))
	if err != nil {
		return nil, err
	}
	if res.Error.Error() != nil {
		return nil, err
	}
	return res.Data, err
}

func (c *Client) CleanPoolCache() error {
	var res base.DataResponse
	uri := "/idaas-app/poolcache"
	err := req.ReqJSON(http.MethodDelete, c.GetUrl(uri), nil, &res, c.SetAuthHeader(nil))
	if err != nil {
		return err
	}
	if res.Error.Error() != nil {
		return err
	}
	return err
}

func (c *Client) BindThirdUser(tuser *ThirdUser) error {
	var res base.DataResponse
	uri := "/idaas-app/userbind"
	err := req.ReqJSON(http.MethodPost, c.GetUrl(uri), tuser, &res, c.SetAuthHeader(nil))
	if err != nil {
		return err
	}
	if res.Error.Error() != nil {
		return err
	}
	return err
}

type ListThirdUsersResp struct {
	*base.DataResponse
	Data []ThirdUser `json:"data"`
}

func (c *Client) ListThirdUsers(tuser *ThirdUser) ([]ThirdUser, error) {
	var res ListThirdUsersResp
	uri := "/idaas-app/thirdusers"
	err := req.ReqJSON(http.MethodPost, c.GetUrl(uri), tuser, &res, c.SetAuthHeader(nil))
	if err != nil {
		return nil, err
	}
	if res.Error.Error() != nil {
		return nil, err
	}
	return res.Data, err
}

type VerifyCaptchasResp struct {
	*base.DataResponse
	Data bool `json:"data"`
}

func (c *Client) VerifyCaptchas(b *VerifyCaptcha) (bool, error) {
	var res VerifyCaptchasResp
	uri := "/idaas-app/captchas"
	err := req.ReqJSON(http.MethodPost, c.GetUrl(uri), b, &res, c.SetAuthHeader(nil))
	if err != nil {
		return false, err
	}
	if res.Error.Error() != nil {
		return false, err
	}
	return res.Data, err
}
