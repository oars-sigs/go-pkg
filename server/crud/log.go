package crud

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
)

type OperationLog struct {
	CommonModel
	Resource      string `json:"resource" gorm:"column:resource;type:varchar(255);size:255;comment:资源"`
	ResourceName  string `json:"resourceName" gorm:"column:resource_name;type:varchar(255);size:255;comment:资源实例"`
	ResourceTitle string `json:"resourceTitle" gorm:"column:resource_title;type:varchar(255);size:255;comment:资源实例标题"`
	Action        string `json:"action" gorm:"column:action;type:varchar(255);size:255;comment:操作"`
	Content       string `json:"content" gorm:"column:content;comment:操作内容"`
}

type CommonModelTitle interface {
	GetTitle() string
}

type OperationLogService interface {
	Create(l any) error
}

func (c *BaseInfoController) genOperationLog(l *OperationLog, oldRes, curRes any) {
	if fn, ok := curRes.(CommonModelTitle); ok {
		l.ResourceTitle = fn.GetTitle()
	}
	if l.Action == UpdateKind || l.Action == CreateKind {
		l.Content = operationLogContent(oldRes, curRes)
	}
	err := c.opt.OperationLogSrv.Create(l)
	if err != nil {
		logrus.Error(err)
	}
}

func operationLogContent(oldRes, curRes any) string {
	oldFileds := getOperationFileds(oldRes)
	curFileds := getOperationFileds(curRes)
	var contents []string
	for k, v := range curFileds {
		if v1, ok := oldFileds[k]; ok {
			if v1 == v {
				continue
			}
			contents = append(contents, fmt.Sprintf("%s:%s->%s", k, v1, v))
		}
		if oldFileds == nil {
			contents = append(contents, fmt.Sprintf("%s:%s", k, v))
		}
	}
	return strings.Join(contents, ";")
}

func getOperationFileds(res any) map[string]string {
	if res == nil {
		return nil
	}
	typeObj := reflect.TypeOf(res).Elem()
	valueObj := reflect.ValueOf(res).Elem()
	fitemts := make(map[string]string)
	for i := 0; i < typeObj.NumField(); i++ {
		item := typeObj.Field(i)
		if item.Name == "Updated" || item.Name == "SearchText" {
			continue
		}
		comment := getTagComment(item.Tag.Get("gorm"))
		if comment != "" {
			fitemts[comment] = fmt.Sprint(valueObj.FieldByName(item.Name).Interface())
		}
	}
	return fitemts
}

func getTagComment(s string) string {
	for _, t := range strings.Split(s, ";") {
		tt := strings.Split(t, ":")
		if len(tt) > 1 && tt[0] == "comment" {
			return tt[1]
		}
	}
	return ""
}
