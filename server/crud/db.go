package crud

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mozillazg/go-pinyin"
	"gorm.io/gorm"
)

func BuildListORM(data any, db *gorm.DB) (*gorm.DB, bool) {
	typeObj := reflect.TypeOf(data).Elem()
	return buildORM(typeObj, db)
}

func BuildGetORM(data any, db *gorm.DB) (*gorm.DB, bool) {
	typeObj := reflect.TypeOf(data).Elem()
	return buildORM(typeObj, db)
}

func BuildCreateGen(data any) {
	BuildGen(data, CreateKind)
}

func BuildUpdateGen(data any) {
	BuildGen(data, UpdateKind)
}

func BuildGen(data any, kind string) {
	typeObj := reflect.TypeOf(data).Elem()
	valueObj := reflect.ValueOf(data).Elem()
	searchText := ""
	for i := 0; i < typeObj.NumField(); i++ {
		item := typeObj.Field(i)
		tags := getTags(item.Tag.Get("gsql"))
		if tags.Search == "default" {
			searchText += fmt.Sprint(valueObj.FieldByName(item.Name).Interface())
		}
		if tags.Search == "pinyin" {
			pys := pinyin.LazyConvert(valueObj.FieldByName(item.Name).Interface().(string), nil)
			spy := ""
			qpy := ""
			for _, py := range pys {
				if len(py) == 0 {
					continue
				}
				spy += string(py[0])
				qpy += py
			}
			searchText += valueObj.FieldByName(item.Name).Interface().(string) + spy + qpy
		}
	}
	valueObj.FieldByName("SearchText").Set(reflect.ValueOf(searchText))
}

func buildORM(typeObj reflect.Type, db *gorm.DB) (*gorm.DB, bool) {
	ok := false
	var ss []string
	for i := 0; i < typeObj.NumField(); i++ {
		tags := getTags(typeObj.Field(i).Tag.Get("gsql"))
		if tags.Table != "" {
			ok = true
			db = db.Table(tags.Table)
		}
		for _, t := range tags.Joins {
			ok = true
			db = db.Joins(t)
		}
		for _, t := range tags.Select {
			ok = true
			ss = append(ss, t)
		}
		for _, t := range tags.Wheres {
			ok = true
			db = db.Where(t)
		}
		if typeObj.Field(i).Type.Kind() == reflect.Struct {
			buildORM(typeObj.Field(i).Type, db)
		}
	}
	if len(ss) > 0 {
		db = db.Select(strings.Join(ss, ","))
	}

	return db, ok
}

type gSql struct {
	Joins  []string
	Wheres []string
	Select []string
	Table  string
	Search string
	ToMany []string
}

func getTags(s string) *gSql {
	var res = &gSql{}
	tags := strings.Split(s, ";")
	for _, tag := range tags {
		keys := strings.Split(tag, ":")
		if len(keys) < 2 {
			keys = append(keys, "default")
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
		case "search":
			res.Search = keys[1]
		case "tomany":
			res.ToMany = append(res.ToMany, keys[1])
		}
	}
	return res
}
