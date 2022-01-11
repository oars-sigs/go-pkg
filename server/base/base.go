package base

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"pkg.oars.vip/go-pkg/e"
)

//BaseController Controller
type BaseController struct {
}

//Health
func (c *BaseController) Health(g *gin.Context) {
	c.OK(g, "success")

}

// Error
func (c *BaseController) Error(g *gin.Context, code int, err error, msg ...string) {
	var res DataResponse
	if err != nil {
		logrus.Error(err)
		res.Detail = err.Error()
		if cerr, ok := err.(*e.Error); ok {
			res.Msg = cerr.Msg()
		}
	}
	if res.Msg == "" && len(msg) != 0 && msg[0] != "" {
		res.Msg = msg[0]
	}
	if res.Msg == "" {
		res.Msg = e.GetCodeMsg(code)
	}
	res.RequestId = GenerateMsgIDFromContext(g)
	res.Code = int32(code)
	g.Set("result", res)
	g.Set("status", code)
	g.AbortWithStatusJSON(http.StatusOK, res)

}

// OK
func (c *BaseController) OK(g *gin.Context, data interface{}, msg ...string) {
	var res DataResponse
	res.Data = data
	if len(msg) != 0 {
		res.Msg = msg[0]
	}
	res.RequestId = GenerateMsgIDFromContext(g)
	res.Code = e.SuccessCode
	g.Set("result", res)
	g.Set("status", http.StatusOK)
	g.AbortWithStatusJSON(http.StatusOK, res)
}

// PageOK
func (c *BaseController) PageOK(g *gin.Context, result interface{}, total, pageNum, pageSize int, msg string) {
	var res page
	res.List = result
	res.Total = total
	res.PageNum = pageNum
	res.PageSize = pageSize
	c.OK(g, res, msg)
}
