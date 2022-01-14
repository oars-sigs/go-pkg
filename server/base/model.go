package base

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Response struct {
	Error
}

type DataResponse struct {
	Response
	Data interface{} `json:"data"`
}

type Page struct {
	Total    int `json:"total"`
	PageNum  int `json:"pageNum" form:"pageNum"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

type page struct {
	Page
	List interface{} `json:"rows"`
}

type Error struct {
	RequestId string `json:"requestId,omitempty"`
	Code      int    `json:"code,omitempty"`
	Msg       string `json:"msg,omitempty"`
	Status    string `json:"status,omitempty"`
	Detail    string `json:"detail,omitempty"`
}

func (e Error) Error() error {
	if e.Code != SuccessCode && e.Code != 200 {
		if e.Detail != "" {
			return errors.New(e.Detail)
		}
		return errors.New(e.Msg)
	}
	return nil
}

const (
	TrafficKey = "X-Request-Id"
	LoggerKey  = "_go-admin-logger-request"
)

// GenerateMsgIDFromContext msgid
func GenerateMsgIDFromContext(c *gin.Context) string {
	requestId := c.GetHeader(TrafficKey)
	if requestId == "" {
		requestId = uuid.New().String()
		c.Header(TrafficKey, requestId)
	}
	return requestId
}
