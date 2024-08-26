package former

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"pkg.oars.vip/go-pkg/constant"
	"pkg.oars.vip/go-pkg/perr"
	"pkg.oars.vip/go-pkg/req"
	"pkg.oars.vip/go-pkg/server/base"
)

func New() {

}

type Client struct {
	Addr  string
	Token string
}

func (c *Client) getUrl(uri string) string {
	return c.Addr + uri
}

func (c *Client) setAuth(uid string, h map[string]string) map[string]string {
	if h == nil {
		h = make(map[string]string)
	}
	if c.Token != "" {
		h[constant.SessionHeader] = "Bearer " + c.Token
	}
	h[constant.ProxyUserIDHeader] = uid
	return h
}

type GetResp struct {
	base.Response
	Data *FlowData `json:"data"`
}

func (c *Client) GetModel(uid string, data *BusData) (*FlowData, error) {
	var resp GetResp
	err := req.ReqJSON(http.MethodPost, c.getUrl("/oars-former/api/v1/user/flowmodels"), data, &resp, c.setAuth(uid, nil))
	if err != nil {
		return nil, err
	}
	if resp.Error.Error() != nil {
		return nil, resp.Error.Error()
	}
	return resp.Data, nil
}

type CreateResp struct {
	base.Response
	Data *BusData `json:"data"`
}

func (c *Client) Create(uid string, data *BusData) (*BusData, error) {
	var resp CreateResp
	err := req.ReqJSON(http.MethodPost, c.getUrl("/oars-former/api/v1/user/flows"), data, &resp, c.setAuth(uid, nil))
	if err != nil {
		return nil, err
	}
	if resp.Error.Error() != nil {
		return nil, resp.Error.Error()
	}
	return resp.Data, nil
}

func (c *Client) Hook(h func(*Hook) error) func(g *gin.Context) {
	return func(g *gin.Context) {
		var p Hook
		err := g.ShouldBindJSON(&p)
		if err != nil {
			g.Writer.WriteHeader(500)
			return
		}
		h(&p)
	}
}

func (c *Client) RefreshUsers(uid, id string) error {
	var resp base.Response
	err := req.ReqJSON(http.MethodPost, c.getUrl("/oars-former/api/v1/admin/flow/refreshusers/"+id), nil, &resp, c.setAuth(uid, nil))
	if err != nil {
		return err
	}
	if resp.Error.Error() != nil {
		return resp.Error.Error()
	}
	return nil
}

func (c *Client) NodeUsers(uid, busId, mark string) ([]ActionUser, error) {
	var resp GetResp
	err := req.ReqJSON(http.MethodGet, c.getUrl("/oars-former/api/v1/user/flows/"+busId), nil, &resp, c.setAuth(uid, nil))
	if err != nil {
		return nil, err
	}
	if resp.Error.Error() != nil {
		return nil, resp.Error.Error()
	}
	for _, act := range resp.Data.Actions {
		if act.Mark == mark {
			return act.Users, nil
		}
	}
	return nil, perr.New("不存在节点" + mark)
}

func (c *Client) Approve(uid string, data *BusTask) error {
	var resp base.DataResponse
	err := req.ReqJSON(http.MethodPost, c.getUrl("/oars-former/api/v1/user/flowapprove"), data, &resp, c.setAuth(uid, nil))
	if err != nil {
		return err
	}
	if resp.Error.Error() != nil {
		return resp.Error.Error()
	}
	return nil
}

type curTasksResp struct {
	base.DataResponse
	Data map[string]BusTask `json:"data"`
}

func (c *Client) CurTasks(uid string, ids []string) (map[string]BusTask, error) {
	var resp curTasksResp
	err := req.ReqJSON(http.MethodPost, c.getUrl("/oars-former/api/v1/user/flowapprove/curtask"), ids, &resp, c.setAuth(uid, nil))
	if err != nil {
		return nil, err
	}
	if resp.Error.Error() != nil {
		return nil, resp.Error.Error()
	}
	return resp.Data, nil
}
