package crud

import (
	"reflect"
	"strings"

	"gorm.io/gorm"
)

func BuildListORM(data any, db *gorm.DB) (*gorm.DB, bool) {
	typeObj := reflect.TypeOf(data).Elem()
	return buildORM(typeObj, db)
}

func BuildGetORM(data any, db *gorm.DB) (*gorm.DB, bool) {
	typeObj := reflect.TypeOf(data)
	return buildORM(typeObj, db)
}

func buildORM(typeObj reflect.Type, db *gorm.DB) (*gorm.DB, bool) {
	ok := false
	for i := 0; i < typeObj.Len(); i++ {
		tags := getTags(typeObj.Field(i).Tag.Get("gsql"))
		if tags.Table != "" {
			ok = true
			db.Table(tags.Table)
		}
		for _, t := range tags.Joins {
			ok = true
			db.Joins(t)
		}
		for _, t := range tags.Select {
			ok = true
			db.Select(t)
		}
		for _, t := range tags.Wheres {
			ok = true
			db.Where(t)
		}
		if typeObj.Kind() == reflect.Struct {
			ok = true
			BuildListORM(typeObj, db)
		}
	}
	return db, ok
}

type gSql struct {
	Joins  []string
	Wheres []string
	Select []string
	Table  string
}

func getTags(s string) *gSql {
	var res = &gSql{}
	tags := strings.Split(s, ";")
	for _, tag := range tags {
		keys := strings.Split(tag, ":")
		if len(keys) != 2 {
			continue
		}
		switch keys[0] {
		case "join":
			res.Joins = append(res.Joins, keys[1])
		case "select":
			res.Select = append(res.Select, keys[1])
		case "where":
			res.Wheres = append(res.Wheres, keys[1])
		case "table":
			res.Table = keys[1]
		}
	}
	return res
}
