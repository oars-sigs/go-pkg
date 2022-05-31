package rick

import (
	"fmt"

	"pkg.oars.vip/go-pkg/req"
	"pkg.oars.vip/go-pkg/server/base"
)

type PutDocRequestParam struct {
	FileId      string   `json:"fileId"`
	FileName    string   `json:"fileName"`
	Tags        []string `json:"tags"`
	FileContent string   `json:"fileContent"`
	CreateTime  int64    `json:"createTime"`
	Uid         string   `json:"uid"`
}

type PutDocResp struct {
	base.DataResponse
	Data PutDocRequestParam `json:"data"`
}

func (c *Client) PutDoc(namespace string, param *PutDocRequestParam) error {
	url := fmt.Sprintf("%s/fulltextquery/api/v1/%s/file", c.cfg.Addr, namespace)
	var resp PutDocResp
	err := req.ReqJSON("PUT", url, param, &resp, nil)
	if err != nil {
		return err
	}
	if resp.Error.Error() != nil {
		return resp.Error.Error()
	}
	fmt.Println(resp.Data)
	return nil
}
