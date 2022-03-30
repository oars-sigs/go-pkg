package idaas

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"pkg.oars.vip/go-pkg/constant"
	"pkg.oars.vip/go-pkg/perr"
	"pkg.oars.vip/go-pkg/req"
	"pkg.oars.vip/go-pkg/server/base"
)

type Config struct {
	//Address idaas 地址
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

func (c *Client) getUrl(uri string) string {
	return c.cfg.Address + uri
}

type DeptsResp struct {
	*base.DataResponse
	Data []Department `json:"data"`
}

func (c *Client) Depts(depts interface{}, tree, useBindPool bool) ([]Department, error) {
	if deptId, ok := depts.(string); ok {
		urlstr := c.getUrl(fmt.Sprintf("/idaas/api/departments?deptId=%s&tree=%v&useBindPool=%v", deptId, tree, useBindPool))
		var depts DeptsResp
		err := req.ReqJSON("GET", urlstr, nil, &depts, c.setAuthHeader(nil))
		if err != nil {
			return nil, err
		}
		if depts.Error.Error() != nil {
			return nil, err
		}
		return depts.Data, err
	}
	if deptIds, ok := depts.([]string); ok {
		urlstr := c.getUrl(fmt.Sprintf("/idaas/api/departmentslist?useBindPool=%v", useBindPool))
		var depts DeptsResp
		err := req.ReqJSON("POST", urlstr, deptIds, &depts, c.setAuthHeader(nil))
		if err != nil {
			return nil, err
		}
		if depts.Error.Error() != nil {
			return nil, err
		}
		return depts.Data, err
	}
	return nil, errors.New("depts must string or []string")
}

type DeptResp struct {
	*base.DataResponse
	Data *Department `json:"data"`
}

func (c *Client) Dept(deptId string, useBindPool bool) (*Department, error) {
	urlstr := c.getUrl(fmt.Sprintf("/idaas/api/departments/%s?useBindPool=%v", deptId, useBindPool))
	var dept DeptResp
	err := req.ReqJSON("GET", urlstr, nil, &dept, c.setAuthHeader(nil))
	if err != nil {
		return nil, err
	}
	if dept.Error.Error() != nil {
		return nil, err
	}
	return dept.Data, err
}

type UsersResp struct {
	*base.DataResponse
	Data []UserInfo `json:"data"`
}

func (c *Client) Users(userIds []string, useBindPool bool) ([]UserInfo, error) {
	urlstr := c.getUrl(fmt.Sprintf("/idaas/api/userslist?useBindPool=%v", useBindPool))
	var users UsersResp
	err := req.ReqJSON("POST", urlstr, userIds, &users, c.setAuthHeader(nil))
	if err != nil {
		return nil, err
	}
	if users.Error.Error() != nil {
		return nil, err
	}
	return users.Data, err
}

type UserResp struct {
	*base.DataResponse
	Data *UserInfo `json:"data"`
}

func (c *Client) User(userId string, useBindPool bool) (*UserInfo, error) {
	urlstr := c.getUrl(fmt.Sprintf("/idaas/api/users/%s?useBindPool=%v", userId, useBindPool))
	var user UserResp
	err := req.ReqJSON("GET", urlstr, nil, &user, c.setAuthHeader(nil))
	if err != nil {
		return nil, err
	}
	if user.Error.Error() != nil {
		return nil, err
	}
	return user.Data, err
}

func (c *Client) UserId(g *gin.Context) (string, error) {
	uid := g.GetHeader(constant.ProxyUserIDHeader)
	if uid == "" {
		return "", perr.ErrUnauthorized
	}
	return uid, nil
}

func (c *Client) Me(g *gin.Context) (*UserInfo, error) {
	uid, err := c.UserId(g)
	if err != nil {
		return nil, err
	}
	return c.User(uid, false)
}

type TokenClaims struct {
	jwt.StandardClaims
}

func (c *Client) Auth(g *gin.Context) (*TokenInfo, error) {
	tokenStr := g.GetHeader(constant.ProxyUserTokenHeader)
	s, _ := base64.URLEncoding.DecodeString(c.cfg.AppSecret)
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %s", token.Header["alg"])
		}
		return s, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		var u TokenInfo
		err := json.Unmarshal([]byte(claims.Subject), &u)
		if err != nil {
			return nil, err
		}
		return &u, nil
	}
	return nil, errors.New("must auth")
}
