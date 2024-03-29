package filebase

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/url"
	"os"
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

	SessionToken string
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
	if c.cfg.AppID != "" {
		headers[constant.ProxyAppIDHeader] = c.cfg.AppID
		headers[constant.SessionHeader] = "App " + base64.StdEncoding.EncodeToString([]byte(c.cfg.AppID+":"+c.cfg.AppSecret))
	}
	if c.cfg.SessionToken != "" {
		headers[constant.SessionHeader] = "Session " + c.cfg.SessionToken
	}
	return headers
}

func (c *Client) PutURL(namespace, parent, name, ext, digest string, size, expireSecond int64) (string, error) {
	expireTime := time.Now().Unix() + expireSecond
	ustr := fmt.Sprintf("%s/filebase/api/v1/%s/files?%s=%s&%s=%s&digest=%s&size=%d&parent=%s&name=%s&type=%s&expireTime=%d",
		c.cfg.Address, namespace, constant.OarsAuthKind, constant.OarsHmacSignatureKind, constant.ProxyAppIDHeader, c.cfg.AppID,
		digest, size, parent, name, ext, expireTime)
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

func (c *Client) GetURL(namespace, id string, expireSecond int64) (string, error) {
	expireTime := time.Now().Unix() + expireSecond
	ustr := fmt.Sprintf("%s/filebase/api/v1/%s/files/%s?%s=%s&%s=%s&expireTime=%d", c.cfg.Address, namespace, id,
		constant.OarsAuthKind, constant.OarsHmacSignatureKind, constant.ProxyAppIDHeader, c.cfg.AppID, expireTime)
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

func (c *Client) FileMd5(path string) (string, error) {
	fs, err := os.Open(path)
	if err != nil {
		return "", err
	}
	info, err := fs.Stat()
	if err != nil {
		return "", err
	}
	filesize := info.Size()
	const filechunk = 4 * 1 << 20
	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))
	hash := md5.New()
	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)
		fs.Read(buf)
		io.WriteString(hash, string(buf))
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (c *Client) FPut(path, namespace, parent, name, ext string) (*FileMetadata, error) {
	digest, err := c.FileMd5(path)
	if err != nil {
		return nil, err
	}
	fs, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	info, err := fs.Stat()
	if err != nil {
		return nil, err
	}
	size := info.Size()
	return c.Put(fs, namespace, parent, name, ext, digest, size)
}

type PutResp struct {
	base.DataResponse
	Data *FileMetadata `json:"data"`
}

func (c *Client) Put(body io.Reader, namespace, parent, name, ext, digest string, size int64) (*FileMetadata, error) {
	ustr := fmt.Sprintf("%s/filebase/api/v1/%s/files?digest=%s&size=%d&parent=%s&name=%s&type=%s",
		c.cfg.Address, namespace, digest, size, parent, name, ext)
	resp, err := req.Req("POST", ustr, body, c.setAuthHeader(nil))
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

type FileInfo struct {
	MD5  string `json:"md5"`
	Size int64  `json:"size"`
	Name string `json:"name"`
}

func (c *Client) FileInfo(path string) (*FileInfo, error) {
	fs, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	info, err := fs.Stat()
	if err != nil {
		return nil, err
	}
	filesize := info.Size()
	const filechunk = 4 * 1 << 20
	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))
	hash := md5.New()
	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)
		fs.Read(buf)
		io.WriteString(hash, string(buf))
	}
	fi := FileInfo{
		Name: info.Name(),
		Size: info.Size(),
		MD5:  hex.EncodeToString(hash.Sum(nil)),
	}
	return &fi, nil
}

func (c *Client) CreateFile(f *FileMetadata) (*FileMetadata, error) {
	ustr := fmt.Sprintf("%s/filebase/api/v1/%s/metadatas", c.cfg.Address, f.Namespace)
	fmt.Println(ustr)
	var out PutResp
	err := req.ReqJSON("POST", ustr, f, &out, c.setAuthHeader(nil))
	if err != nil {
		return nil, err
	}
	return out.Data, nil
}
