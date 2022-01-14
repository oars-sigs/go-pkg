package filebase

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"time"

	"pkg.oars.vip/go-pkg/constant"
	"pkg.oars.vip/go-pkg/req"
	"pkg.oars.vip/go-pkg/server/base"
)

type Config struct {
	//Address filebase 地址
	Address string
	//AppID 应用id
	AppID string
	//AppSecret 应用secret
	AppSecret string
}

//Client filebase client
type Client struct {
	cfg *Config
}

func New(cfg *Config) *Client {
	return &Client{cfg}
}

func (c *Client) setAuthHeader(headers map[string]string) map[string]string {
	if headers == nil {
		headers = make(map[string]string)
	}
	headers[constant.ProxyAppIDHeader] = c.cfg.AppID
	headers[constant.SessionHeader] = "App " + base64.StdEncoding.EncodeToString([]byte(c.cfg.AppID+":"+c.cfg.AppSecret))
	return headers
}

func (c *Client) PutURL(parent, name, ext, digest string, size, expireSecond int64) (string, error) {
	expireTime := time.Now().Unix() + expireSecond
	ustr := fmt.Sprintf("%s/filebase/api/v1/app/files?%s=%s&digest=%s&size=%d&parent=%s&name=%s&type=%s&expireTime=%d",
		c.cfg.Address, constant.OarsAuthKind, constant.OarsHmacSignatureKind, digest, size, parent, name, ext, expireTime)
	u, err := url.Parse(ustr)
	if err != nil {
		return "", err
	}
	qs := u.Query()
	s := SignURL(u, c.cfg.AppSecret)
	qs.Set(constant.SignatureKey, s)
	u.RawQuery = qs.Encode()
	return u.String(), nil
}

func (c *Client) Get(id string) (io.ReadCloser, error) {
	ustr := fmt.Sprintf("%s/filebase/api/v1/app/files/%s", c.cfg.Address, id)
	return req.Req("GET", ustr, nil, c.setAuthHeader(nil))
}

func (c *Client) Md5(r io.Reader) string {
	md5h := md5.New()
	io.Copy(md5h, r)
	return fmt.Sprintf("%x", md5h.Sum([]byte("")))
}

type PutResp struct {
	base.DataResponse
	Data *FileMetadata `json:"data"`
}

func (c *Client) Put(body io.Reader, parent, name, ext string, size int64) (*FileMetadata, error) {
	var buf bytes.Buffer
	tr := io.TeeReader(body, &buf)
	digest := c.Md5(&buf)
	ustr := fmt.Sprintf("%s/filebase/api/v1/app/files?digest=%s&size=%d&parent=%s&name=%s&type=%s",
		c.cfg.Address, digest, size, parent, name, ext)
	resp, err := req.Req("POST", ustr, tr, c.setAuthHeader(nil))
	if err != nil {
		return nil, err
	}
	defer resp.Close()
	data, err := ioutil.ReadAll(resp)
	if err != nil {
		return nil, err
	}
	var res PutResp
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	if res.Error.Error() != nil {
		return nil, res.Error.Error()
	}
	return res.Data, nil
}
