package base

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"pkg.oars.vip/go-pkg/perr"
)

const (
	SuccessCode = 10000
)

//BaseController Controller
type BaseController struct {
}

//Health
func (c *BaseController) Health(g *gin.Context) {
	c.OK(g, "success")

}

// Error
func (c *BaseController) Error(g *gin.Context, err error, msg ...string) {
	var res DataResponse
	if err == nil {
		err = perr.ErrUnkown
	}
	logrus.Error(err)
	res.Detail = err.Error()
	if cerr, ok := err.(*perr.Error); ok {
		res.Msg = cerr.Msg()
		res.Code = cerr.Code()
	} else {
		res.Msg = perr.ErrInternal.Msg()
		res.Code = perr.ErrInternal.Code()
	}
	if res.Msg == "" && len(msg) != 0 && msg[0] != "" {
		res.Msg = msg[0]
	}
	res.RequestId = GenerateMsgIDFromContext(g)
	g.Set("result", res)
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
	res.Code = SuccessCode
	g.Set("data", res)
	g.AbortWithStatusJSON(http.StatusOK, res)
}

// PageOK
func (c *BaseController) PageOK(g *gin.Context, result interface{}, total, pageNum, pageSize int, msg ...string) {
	var res page
	res.List = result
	res.Total = total
	res.PageNum = pageNum
	res.PageSize = pageSize
	c.OK(g, res, msg...)
}

func (c *BaseController) PageQuery(g *gin.Context) (*Page, error) {
	var res Page
	err := g.BindQuery(&res)
	if err == nil {
		if res.PageNum == 0 {
			res.PageNum = 1
		}
		if res.PageSize == 0 {
			res.PageSize = 10
		}
	}
	return &res, err
}
