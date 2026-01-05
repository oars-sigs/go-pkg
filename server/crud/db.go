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
	Search       string
	SearchText   string
	SeniorSearch string
	SortField    string
	Order        string
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
		tags := getTags(item.Tag.Get("gsql"), false)
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
	if valueObj.FieldByName("SearchText").IsValid() {
		valueObj.FieldByName("SearchText").Set(reflect.ValueOf(searchText))
	}
}

func buildORM(typeObj reflect.Type, db *gorm.DB, opt *BuildORMOption) (*gorm.DB, bool) {
	if opt == nil {
		opt = new(BuildORMOption)
	}
	selectFileds := make([]string, 0)
	searchFileds := make([]string, 0)
	json2f := make(map[string]string)
	var res *buildOrmRes
	db, res = buildORMItem(typeObj, db, &selectFileds, &searchFileds, json2f)
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
				if !strings.Contains(s, ".") {
					s = "m." + s
				}
			}
			searchQs = append(searchQs, s+" LIKE ?")
			ps = append(ps, `%`+opt.Search+`%`)
		}
		db = db.Where(strings.Join(searchQs, " OR "), ps...)
	}
	if opt.SearchText != "" {
		s := "search_text"
		if res.IsJoin {
			s = "m." + s
		}
		db = db.Where(s+" LIKE ?", opt.SearchText)
	}

	//高级搜索
	if opt.SeniorSearch != "" {
		db = buildSeniorSearch(db, opt.SeniorSearch, res.json2f)
	}
	if opt.SortField != "" {
		if opt.Order == "2" {
			db = db.Order("`" + res.json2f[opt.SortField] + "` desc")
		} else {
			db = db.Order("`" + res.json2f[opt.SortField] + "` asc")
		}
	}
	for _, v := range opt.Order {
		db = db.Order(v)
	}

	return db, res.Change
}

type buildOrmRes struct {
	IsJoin  bool
	IsTable bool
	Change  bool
	json2f  map[string]string
}

func buildORMItem(typeObj reflect.Type, db *gorm.DB, selectFileds, searchFileds *[]string, json2f map[string]string) (*gorm.DB, *buildOrmRes) {
	ok := false
	isJoin := false
	isTable := false
	del := existDeleteAt(typeObj)
	for i := 0; i < typeObj.NumField(); i++ {
		json2db(typeObj.Field(i).Tag.Get("json"), typeObj.Field(i).Tag.Get("gorm"), json2f)
		tags := getTags(typeObj.Field(i).Tag.Get("gsql"), del)
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
		for _, t := range tags.Order {
			db = db.Order(t)
		}
		if typeObj.Field(i).Type.Kind() == reflect.Struct {
			vdb, res := buildORMItem(typeObj.Field(i).Type, db, selectFileds, searchFileds, json2f)
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

	return db, &buildOrmRes{IsJoin: isJoin, IsTable: isTable, Change: ok, json2f: json2f}
}

func existDeleteAt(typeObj reflect.Type) bool {
	for i := 0; i < typeObj.NumField(); i++ {
		if typeObj.Field(i).Name == "DeletedAt" {
			return true
		}
		if typeObj.Field(i).Type.Kind() == reflect.Struct {
			if existDeleteAt(typeObj.Field(i).Type) {
				return true
			}
		}
	}
	return false
}

type gSql struct {
	Joins      []string
	Wheres     []string
	Select     []string
	Table      string
	Search     string
	SearchText string
	ToMany     []string
	Order      []string
}

func getTags(s string, del bool) *gSql {
	var res = &gSql{}
	tags := strings.Split(s, ";")
	for _, tag := range tags {
		keys := strings.Split(tag, ":")
		if len(keys) < 2 {
			keys = append(keys, "default")
		}
		switch keys[0] {
		case "join":
			d := keys[1]
			if !strings.Contains(d, "deleted_at") {
				ss := strings.Split(d, " as ")
				if len(ss) > 1 && del {
					d = d + " AND " + strings.Split(strings.TrimSpace(ss[1]), " ")[0] + ".deleted_at is NULL"
				}
			}
			res.Joins = append(res.Joins, d)
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
		case "order":
			res.Order = append(res.Order, keys[1])
		}
	}
	return res
}

func json2db(j, o string, data map[string]string) {
	if o == "" || j == "" {
		return
	}
	ss := strings.Split(o, ";")
	for _, s := range ss {
		kk := strings.Split(s, ":")
		if len(kk) > 1 {
			if kk[0] == "column" {
				data[j] = kk[1]
			}
		}
	}
}

func getDBTag(data any, json2f map[string]string) {
	typeObj := reflect.TypeOf(data).Elem()
	for i := 0; i < typeObj.NumField(); i++ {
		for k := range json2f {
			jk := typeObj.Field(i).Tag.Get("json")
			if k == jk {
				json2db(jk, typeObj.Field(i).Tag.Get("gorm"), json2f)
			}
		}
		if typeObj.Field(i).Type.Kind() == reflect.Struct {
			getdbTag(typeObj.Field(i).Type, json2f)
		}
	}
}
