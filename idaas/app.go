package idaas

import (
	"encoding/base64"
	"strings"

	"github.com/gin-gonic/gin"
	"pkg.oars.vip/go-pkg/constant"
	"pkg.oars.vip/go-pkg/perr"
)

func (c *Client) AppID(g *gin.Context) (string, error) {
	appId := g.GetHeader(constant.ProxyAppIDHeader)
	if appId != "" {
		return appId, nil
	}
	token := strings.TrimPrefix(g.GetHeader(constant.SessionHeader), "App ")
	b, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", perr.ErrUnauthorized.Set(err)
	}
	tokens := strings.Split(string(b), ":")
	if len(tokens) != 2 {
		return "", perr.ErrUnauthorized
	}
	return tokens[0], nil
}
