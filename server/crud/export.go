package crud

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
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

func ExportExcel(sheets []ExportOption) *excelize.File {
	xlsx := excelize.NewFile()
	isfirst := true
	for _, sheet := range sheets {
		xlsx.NewSheet(sheet.Name)
		if isfirst && sheet.Name != "Sheet1" {
			xlsx.DeleteSheet("Sheet1")
			isfirst = false
		}

		sliceValue := reflect.ValueOf(sheet.Data)
		var headers [][]HeaderItem
		var allHeaders []HeaderItem
		style, _ := xlsx.NewStyle(&excelize.Style{
			Border: []excelize.Border{
				{
					Style: 1,
					Type:  "left",
					Color: "0000000",
				},
				{
					Style: 1,
					Type:  "right",
					Color: "0000000",
				},
				{
					Style: 1,
					Type:  "top",
					Color: "0000000",
				},
				{
					Style: 1,
					Type:  "bottom",
					Color: "0000000",
				},
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
		})
		type posCell struct {
			v string
			p int
		}
		mergeCell := make(map[int][]posCell)
		for i := 0; i < sliceValue.Len(); i++ {
			if len(sheet.Headers) == 0 {
				sheet.Headers = make([]HeaderItem, 0)
				GetExportHeaders(sliceValue.Index(i).Type(), &sheet.Headers)
				if len(sheet.Extends) > 0 {
					sheet.Headers = append(sheet.Headers, sheet.Extends...)
				}
				sort.Slice(sheet.Headers, func(i, j int) bool {
					return sheet.Headers[i].Order < sheet.Headers[j].Order
				})

				if len(sheet.Ignores) > 0 {
					hs := make([]HeaderItem, 0)
					for _, h := range sheet.Headers {
						ignore := false
						for _, s := range sheet.Ignores {
							if s == h.Filed {
								ignore = true
								break
							}
						}
						if !ignore {
							hs = append(hs, h)
						}
					}
					sheet.Headers = hs
				}
			}
			if len(headers) == 0 {
				headers = GenHeaders(sheet.Headers)
			}
			headerRow := len(headers)
			if len(allHeaders) == 0 {
				var vallHeaders []HeaderItem
				for _, hs := range headers {
					for _, h := range hs {
						if h.W == 0 {
							vallHeaders = append(vallHeaders, h)
						}
					}
				}
				allHeaders = vallHeaders
				indexspan := make(map[int]bool)
				allHeaders = make([]HeaderItem, len(vallHeaders))
				for _, hs := range headers {
					w := 0
					for _, h := range hs {
						w = GetNextIndex(indexspan, w)
						if h.W == 0 {
							allHeaders[w] = h
							indexspan[w] = true
							w++
						} else {
							w += h.W
						}
					}
				}
			}

			for j, h := range allHeaders {
				var data any
				if v, ok := sliceValue.Index(i).Interface().(map[string]interface{}); ok {
					data = v[h.Filed]

					if data != nil && reflect.TypeOf(data).Kind() == reflect.Ptr {
						vdata := reflect.ValueOf(data)
						if !vdata.IsNil() {
							data = vdata.Elem().Interface()
						} else {
							data = ""
						}
					}
				} else {
					v := sliceValue.Index(i).FieldByName(h.Filed)
					if !v.IsValid() {
						continue
					}
					if v.Kind() == reflect.Ptr {
						if !v.IsNil() {
							data = v.Elem().Interface()
						} else {
							data = ""
						}
					} else {
						data = v.Interface()
					}
				}

				if fn, ok := h.Fns["format"]; ok {
					switch fn {
					case "toDate":
						data = time.UnixMilli(data.(int64)).Format("2006-01-02")
					case "toDatetime":
						data = time.UnixMilli(data.(int64)).Format("2006-01-02 15:04:05")
					}
				}
				if h.AutoMerge {
					if len(mergeCell[j]) == 0 {
						mergeCell[j] = append(mergeCell[j], posCell{
							p: i + headerRow,
							v: fmt.Sprint(data),
						})
					} else if mergeCell[j][len(mergeCell[j])-1].v != fmt.Sprint(data) {
						mergeCell[j] = append(mergeCell[j], posCell{
							p: i + headerRow,
							v: fmt.Sprint(data),
						})
					}
				}
				xlsx.SetCellValue(sheet.Name, GetEcelAxis(i+headerRow+1, j+1), data)
				xlsx.SetCellStyle(sheet.Name, GetEcelAxis(i+headerRow+1, j+1), GetEcelAxis(i+headerRow+1, j+1), style)
			}
		}

		//表头
		indexspan := make(map[int]bool)
		for i, row := range headers {
			j := 0
			cindexspan := make(map[int]bool)
			for _, col := range row {
				j = GetNextIndex(indexspan, j)
				if col.D == 0 && col.L < len(headers) {
					xlsx.MergeCell(sheet.Name, GetEcelAxis(i+1, j+1),
						GetEcelAxis(i+1+(len(headers)-col.L), j+1))
					cindexspan[j] = true
					xlsx.SetCellStyle(sheet.Name, GetEcelAxis(i+1, j+1),
						GetEcelAxis(i+1+(len(headers)-col.L), j+1), style)
				}
				xlsx.SetCellValue(sheet.Name, GetEcelAxis(i+1, j+1), col.Title)
				if col.W > 0 {
					xlsx.MergeCell(sheet.Name, GetEcelAxis(i+1, j+1), GetEcelAxis(i+1, j+col.W))
					xlsx.SetCellStyle(sheet.Name, GetEcelAxis(i+1, j+1), GetEcelAxis(i+1, j+col.W), style)
					j += col.W
				} else {
					xlsx.SetCellStyle(sheet.Name, GetEcelAxis(i+1, j+1), GetEcelAxis(i+1, j+1), style)
					j++
				}

				if col.Width > 0 {
					xlsx.SetColWidth(sheet.Name, GetColumnIndex(i+1),
						GetColumnIndex(i+1), float64(col.Width))
				}
				for j, m := range mergeCell {
					for i, p := range m {
						if i == len(m)-1 {
							if sliceValue.Len()+len(headers)-p.p > 2 {
								xlsx.MergeCell(sheet.Name, GetEcelAxis(p.p+1, j+1), GetEcelAxis(sliceValue.Len()+len(headers), j+1))
							}
							continue
						}
						if m[i+1].p-p.p > 1 {
							xlsx.MergeCell(sheet.Name, GetEcelAxis(p.p+1, j+1), GetEcelAxis(m[i+1].p, j+1))
						}
					}
				}
			}
			for k, v := range cindexspan {
				indexspan[k] = v
			}
		}
	}
	return xlsx
}

func WriteExcelToCtx(xlsx *excelize.File, filename string, g *gin.Context) {
	g.Writer.Header().Set("Content-Type", "application/octet-stream")
	disposition := fmt.Sprintf("attachment; filename=\"%s.xlsx\"", filename)
	g.Writer.Header().Set("Content-Disposition", disposition)
	_ = xlsx.Write(g.Writer)
}
