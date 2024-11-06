package crud

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExportOption struct {
	Name    string
	Data    interface{}
	Headers []HeaderItem
	Extends []HeaderItem
	Ignores []string
}

type HeaderItem struct {
	Title     string
	Filed     string
	Order     int
	Fns       map[string]string
	Width     int
	AutoMerge bool
	Hide      bool
	Children  []HeaderItem
	W         int
	D         int
	L         int
}

func GenHeaders(hs []HeaderItem) [][]HeaderItem {
	var dres []HeaderItem
	for _, h := range hs {
		if h.Hide {
			continue
		}
		dres = append(dres, h)
	}
	tree := HeaderItem{
		Children: dres,
	}
	genHeaderTree(&tree, -1)
	res := make([][]HeaderItem, tree.D)
	rangeHeaderTree(dres, &res)
	return res
}

func rangeHeaderTree(hs []HeaderItem, res *[][]HeaderItem) {
	for _, h := range hs {
		rangeHeaderTree(h.Children, res)
		if (*res)[h.L-1] == nil {
			(*res)[h.L-1] = make([]HeaderItem, 0)
		}
		(*res)[h.L-1] = append((*res)[h.L-1], h)
	}
}

func genHeaderTree(h *HeaderItem, l int) {
	h.L = l + 1
	if len(h.Children) == 0 {
		h.W = 0
		h.D = 0
		return
	}
	h.D = -1
	//fmt.Println(h.Title, len(h.Children))
	for i := range h.Children {
		genHeaderTree(&h.Children[i], h.L)
		if h.Children[i].W == 0 {
			h.W += 1
		} else {
			h.W += h.Children[i].W
		}
		if h.Children[i].D+1 > h.D {
			h.D = h.Children[i].D + 1
		}
	}
	//fmt.Println(h.Title, h.W)
}

type CommonModelExport interface {
	ExportORM(db *gorm.DB, c any, g *gin.Context) (*gorm.DB, any, *ExportOption, error)
}

type ExportListOption struct {
	Sheets map[string]CommonModelExport
	Name   string
}

type CommonModelExportList interface {
	ListExportORM(g *gin.Context) *ExportListOption
}

func GetExportHeaders(v reflect.Type, hs *[]HeaderItem) {
	for n := 0; n < v.NumField(); n++ {
		f := v.Field(n)
		if f.Type.Kind() == reflect.Struct {
			GetExportHeaders(f.Type, hs)
		}
		exportTag := f.Tag.Get("export")
		if exportTag == "" {
			continue
		}
		items := strings.Split(exportTag, ";")
		fns := make(map[string]string)
		for _, item := range items {
			keys := strings.Split(item, ":")
			if len(keys) > 1 {
				fns[keys[0]] = keys[1]
				continue
			}
			ss := strings.Split(item, ",")
			h := HeaderItem{Title: ss[0], Filed: f.Name, Fns: fns}
			if len(ss) > 1 {
				h.Order, _ = strconv.Atoi(ss[1])
			}
			*hs = append(*hs, h)
		}

	}
}

func GetNextIndex(s map[int]bool, i int) int {
	if _, ok := s[i]; ok {
		i = GetNextIndex(s, i+1)
	}
	return i
}

func GetEcelAxis(row int, columnCount int) string {
	var column = GetColumnIndex(columnCount)
	return fmt.Sprintf("%s%d", column, row)
}

// 获取excel的列索引
var columnIndexList = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func GetColumnIndex(num int) string {
	num--
	var column = columnIndexList[num%26]
	for num = num / 26; num > 0; num = num/26 - 1 {
		column = columnIndexList[(num-1)%26] + column
	}
	return column
}
