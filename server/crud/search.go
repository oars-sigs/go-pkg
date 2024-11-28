package crud

import (
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

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
			db = db.Where(fmt.Sprintf("%s %s ?", f, getOp(rule.Operator)), getV(rule.Value, rule.Operator))
		} else if len(db.Statement.Joins) > 0 {
			db = db.Where(fmt.Sprintf("m.%s %s ?", k, getOp(rule.Operator)), getV(rule.Value, rule.Operator))
		} else {
			db = db.Where(fmt.Sprintf("%s %s ?", k, getOp(rule.Operator)), getV(rule.Value, rule.Operator))
		}
	}
	return db
}

func getOp(s string) string {
	switch s {
	case "==":
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
