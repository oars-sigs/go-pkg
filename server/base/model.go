package base

import (
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
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

type page struct {
	Page
	List interface{} `json:"list"`
}

type Error struct {
	RequestId string `json:"requestId,omitempty"`
	Code      int32  `json:"code,omitempty"`
	Msg       string `json:"msg,omitempty"`
	Status    string `json:"status,omitempty"`
	Detail    string `json:"detail,omitempty"`
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
