package filebase

import (
	"fmt"

	"pkg.oars.vip/go-pkg/req"
	"pkg.oars.vip/go-pkg/server/base"
)

type ListNamespaceResp struct {
	base.DataResponse
	Data []Namespace `json:"data"`
}

func (c *Client) ListNamespace(userId string) ([]Namespace, error) {
	ustr := fmt.Sprintf("%s/filebase/app/api/v1/namespaces?userId=%s", c.cfg.Address, userId)
	var resp ListNamespaceResp
	err := req.ReqJSON("GET", ustr, nil, &resp, c.setAuthHeader(nil))
	if err != nil {
		return nil, err
	}
	if resp.Error.Error() != nil {
		return nil, err
	}
	return resp.Data, nil
}

type CreateNamespaceResp struct {
	base.DataResponse
	Data *Namespace
}

func (c *Client) CreateNamespace(ns *Namespace, userId string) (*Namespace, error) {
	ustr := fmt.Sprintf("%s/filebase/app/api/v1/namespaces?userId=%s", c.cfg.Address, userId)
	var resp CreateNamespaceResp
	err := req.ReqJSON("GET", ustr, ns, &resp, c.setAuthHeader(nil))
	if err != nil {
		return nil, err
	}
	if resp.Error.Error() != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (c *Client) UpdateNamespace(ns *Namespace) (*Namespace, error) {
	ustr := fmt.Sprintf("%s/filebase/app/api/v1/namespaces/%s", c.cfg.Address, ns.ID)
	var resp CreateNamespaceResp
	err := req.ReqJSON("GET", ustr, ns, &resp, c.setAuthHeader(nil))
	if err != nil {
		return nil, err
	}
	if resp.Error.Error() != nil {
		return nil, err
	}
	return resp.Data, nil
}
