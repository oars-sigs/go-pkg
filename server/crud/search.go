package crud

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

func BuildSeniorSearch(data any, db *gorm.DB, text string) *gorm.DB {
	typeObj := reflect.TypeOf(data).Elem()
	json2f := make(map[string]string)
	getdbTag(typeObj, json2f)
	return buildSeniorSearch(db, text, json2f)
}

func getdbTag(typeObj reflect.Type, json2f map[string]string) {
	for i := 0; i < typeObj.NumField(); i++ {
		json2db(typeObj.Field(i).Tag.Get("json"), typeObj.Field(i).Tag.Get("gorm"), json2f)
		if typeObj.Field(i).Type.Kind() == reflect.Struct {
			getdbTag(typeObj.Field(i).Type, json2f)
		}
	}
}

type SeniorSearchRule struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    any    `json:"value"`
}

func buildSeniorSearch(db *gorm.DB, text string, json2f map[string]string) *gorm.DB {
	var rules []SeniorSearchRule
	json.Unmarshal([]byte(text), &rules)
	sfs := make(map[string]string)
	selects := db.Statement.Selects
	for _, s := range selects {
		ll := strings.Split(s, "as")
		if len(ll) > 1 {
			sfs[strings.TrimSpace(ll[1])] = strings.TrimSpace(ll[0])

		}
	}
	for _, rule := range rules {
		k, ok := json2f[rule.Field]
		if !ok {
			continue
		}
		if f, ok := sfs[k]; ok {
			db = db.Where(fmt.Sprintf("%s %s ?", f, getOp(rule.Value, rule.Operator)), getV(rule.Value, rule.Operator))
		} else if len(db.Statement.Joins) > 0 {
			db = db.Where(fmt.Sprintf("%s.%s %s ?", db.Statement.Table, k, getOp(rule.Value, rule.Operator)), getV(rule.Value, rule.Operator))
		} else {
			db = db.Where(fmt.Sprintf("%s %s ?", k, getOp(rule.Value, rule.Operator)), getV(rule.Value, rule.Operator))
		}
	}
	return db
}

func getOp(s any, op string) string {
	switch op {
	case "==":
		if reflect.TypeOf(s).Kind() == reflect.Slice {
			return "in"
		}
		return "="
	case "contains":
		return "LIKE"
	}
	return "="
}

func getV(s any, op string) any {
	if op == "contains" {
		return fmt.Sprintf("%%%s%%", s)
	}
	return s
}
