package idaas

import (
	"pkg.oars.vip/go-pkg/req"
	"pkg.oars.vip/go-pkg/server/base"
)

// PermissionInitData 初始化权限数据
func (c *Client) PermissionInitData(data *InitPermissionData) error {
	urlstr := c.GetUrl("/idaas-app/permissions/initdata")
	var resp base.DataResponse
	err := req.ReqJSON("POST", urlstr, data, &resp, c.SetAuthHeader(nil))
	if err != nil {
		return err
	}
	if resp.Error.Error() != nil {
		return err
	}
	return nil
}

// PermissionPutRule 替换权限规则
func (c *Client) PermissionPutRule(data *PermissionRulePutParam) error {
	urlstr := c.GetUrl("/idaas-app/permissions/rules")
	var resp base.DataResponse
	err := req.ReqJSON("PUT", urlstr, data, &resp, c.SetAuthHeader(nil))
	if err != nil {
		return err
	}
	if resp.Error.Error() != nil {
		return err
	}
	return nil
}

// PermissionCreateRule 创建权限规则
func (c *Client) PermissionCreateRule(data *PermissionRolebindings) error {
	urlstr := c.GetUrl("/idaas-app/permissions/rules")
	var resp base.DataResponse
	err := req.ReqJSON("POST", urlstr, data, &resp, c.SetAuthHeader(nil))
	if err != nil {
		return err
	}
	if resp.Error.Error() != nil {
		return err
	}
	return nil
}

// PermissionDeleteRule 删除权限规则
func (c *Client) PermissionDeleteRule(data *PermissionRolebindings) error {
	urlstr := c.GetUrl("/idaas-app/permissions/rules")
	var resp base.DataResponse
	err := req.ReqJSON("DELETE", urlstr, data, &resp, c.SetAuthHeader(nil))
	if err != nil {
		return err
	}
	if resp.Error.Error() != nil {
		return err
	}
	return nil
}

type PermissionListRuleResp struct {
	*base.DataResponse
	Data []PermissionRolebindings `json:"data"`
}

// PermissionListRule 权限规则列表
func (c *Client) PermissionListRule(data *PermissionRolebindings) ([]PermissionRolebindings, error) {
	urlstr := c.GetUrl("/idaas-app/permissions/ruleslist?userDetail=true")
	var resp PermissionListRuleResp
	err := req.ReqJSON("POST", urlstr, data, &resp, c.SetAuthHeader(nil))
	if err != nil {
		return nil, err
	}
	if resp.Error.Error() != nil {
		return nil, err
	}
	return resp.Data, nil
}

type PermissionEnforceResp struct {
	*base.DataResponse
	Data bool `json:"data"`
}

// PermissionEnforce 权限检验
func (c *Client) PermissionEnforce(data EnforceParam) (bool, error) {
	urlstr := c.GetUrl("/idaas-app/permissions/enforce")
	var resp PermissionEnforceResp
	err := req.ReqJSON("POST", urlstr, data, &resp, c.SetAuthHeader(nil))
	if err != nil {
		return false, err
	}
	if resp.Error.Error() != nil {
		return false, err
	}
	return resp.Data, nil
}

type PermissionResourcesResp struct {
	*base.DataResponse
	Data *ResourceNames `json:"data"`
}

// PermissionResources 权限资源
func (c *Client) PermissionResources(data EnforceParam) (*ResourceNames, error) {
	urlstr := c.GetUrl("/idaas-app/permissions/resources")
	var resp PermissionResourcesResp
	err := req.ReqJSON("POST", urlstr, data, &resp, c.SetAuthHeader(nil))
	if err != nil {
		return nil, err
	}
	if resp.Error.Error() != nil {
		return nil, err
	}
	return resp.Data, nil
}
