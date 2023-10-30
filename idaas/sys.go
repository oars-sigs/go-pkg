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

func (c *Client) Sys(param string) (map[string]interface{}, error) {
	var res SysResp
	err := req.ReqJSON(http.MethodGet, c.getUrl("/idaas-app/sys?param="+param), nil, &res, c.setAuthHeader(nil))
	if err != nil {
		return nil, err
	}
	if res.Error.Error() != nil {
		return nil, err
	}
	return res.Data, err
}
