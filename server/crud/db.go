package crud

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mozillazg/go-pinyin"
	"gorm.io/gorm"
)

type ResourceTable interface {
	TableName() string
}

type BuildORMOption struct {
	Search     string
	SearchText string
}

func BuildListORM(data any, db *gorm.DB, opt *BuildORMOption) (*gorm.DB, bool) {
	typeObj := reflect.TypeOf(data).Elem()
	return buildORM(typeObj, db, opt)
}

func BuildGetORM(data any, db *gorm.DB) (*gorm.DB, bool) {
	typeObj := reflect.TypeOf(data).Elem()
	return buildORM(typeObj, db, nil)
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
		if tags.SearchText == "default" {
			searchText += fmt.Sprint(valueObj.FieldByName(item.Name).Interface())
		}
		if tags.SearchText == "pinyin" {
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

func buildORM(typeObj reflect.Type, db *gorm.DB, opt *BuildORMOption) (*gorm.DB, bool) {
	if opt == nil {
		opt = new(BuildORMOption)
	}
	selectFileds := make([]string, 0)
	searchFileds := make([]string, 0)
	var res *buildOrmRes
	db, res = buildORMItem(typeObj, db, &selectFileds, &searchFileds)
	//如果有关联表且没定义table
	if res.IsJoin && !res.IsTable {
		modelValue := reflect.New(typeObj)
		if tabler, ok := modelValue.Interface().(ResourceTable); ok {
			db = db.Table(tabler.TableName() + " as m")
		}
		selectFileds = append(selectFileds, "m.*")
	}
	if len(selectFileds) > 0 {
		db = db.Select(strings.Join(selectFileds, ","))
	}
	//搜索字段
	if opt.Search != "" && len(searchFileds) > 0 {

		var searchQs []string
		var ps []interface{}
		for _, s := range searchFileds {
			if res.IsJoin {
				fmt.Println(s, strings.Split(s, "."))
				if !strings.Contains(s, ".") {
					s = "m." + s
				}
			}
			searchQs = append(searchQs, s+" LIKE ?")
			ps = append(ps, `%`+opt.Search+`%`)
		}
		fmt.Println(strings.Join(searchQs, " OR "), ps)
		db = db.Where(strings.Join(searchQs, " OR "), ps...)
	}
	if opt.SearchText != "" {
		s := "search_text"
		if res.IsJoin {
			s = "m." + s
		}
		db = db.Where(s+" LIKE ?", opt.SearchText)
	}
	return db, res.Change
}

type buildOrmRes struct {
	IsJoin  bool
	IsTable bool
	Change  bool
}

func buildORMItem(typeObj reflect.Type, db *gorm.DB, selectFileds, searchFileds *[]string) (*gorm.DB, *buildOrmRes) {
	ok := false
	isJoin := false
	isTable := false
	db = db.Debug()
	for i := 0; i < typeObj.NumField(); i++ {
		tags := getTags(typeObj.Field(i).Tag.Get("gsql"))
		if tags.Table != "" {
			ok = true
			db = db.Table(tags.Table)
			isTable = true
		}
		for _, t := range tags.Joins {
			ok = true
			isJoin = true
			db = db.Joins(t)
		}
		for _, t := range tags.Select {
			ok = true
			*selectFileds = append(*selectFileds, t)
		}
		for _, t := range tags.Wheres {
			ok = true
			db = db.Where(t)
		}
		if tags.Search != "" {
			*searchFileds = append(*searchFileds, tags.Search)
		}
		if typeObj.Field(i).Type.Kind() == reflect.Struct {
			vdb, res := buildORMItem(typeObj.Field(i).Type, db, selectFileds, searchFileds)
			if res.IsJoin {
				isJoin = res.IsJoin
			}
			if res.IsTable {
				isTable = res.IsTable
			}
			if res.Change {
				ok = res.Change
			}
			db = vdb
		}
	}

	return db, &buildOrmRes{IsJoin: isJoin, IsTable: isTable, Change: ok}
}

type gSql struct {
	Joins      []string
	Wheres     []string
	Select     []string
	Table      string
	Search     string
	SearchText string
	ToMany     []string
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
		case "searchText":
			res.SearchText = keys[1]
		case "tomany":
			res.ToMany = append(res.ToMany, keys[1])
		case "search":
			res.Search = keys[1]
		}
	}
	return res
}
