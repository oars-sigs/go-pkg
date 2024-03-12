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
