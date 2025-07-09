package crud

import (
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/xuri/excelize/v2"
	"pkg.oars.vip/go-pkg/filebase"
)

type ImportOption struct {
	Sheet     string         `json:"sheet"`
	OffsetRow int            `json:"offsetRow"`
	Column    []ImportColumn `json:"column"`
}

type ImportColumn struct {
	Label    string          `json:"label"`
	Prop     string          `json:"prop"`
	Type     string          `json:"type"`
	TableDic *ImportTableDic `json:"tableDic"`
	DataDic  []ImportDataDic `json:"dataDic"`
}

type ImportTableDic struct {
	Table string `json:"table"`
	Label string `json:"label"`
	Value string `json:"value"`
}

type ImportDataDic struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func ParseImport(g *gin.Context, opt *ImportOption, res interface{}) error {
	fs, err := g.FormFile("file")
	if err != nil {
		return err
	}
	path := "/tmp/" + uuid.NewString()
	err = g.SaveUploadedFile(fs, path)
	if err != nil {
		return err
	}
	defer os.RemoveAll(path)
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}
	if opt.Sheet == "" {
		opt.Sheet = f.GetSheetName(0)
	}
	defer f.Close()
	rows, err := f.GetRows(opt.Sheet)
	if err != nil {
		return err
	}

	var hs = make(map[int]ImportColumn)
	for i, row := range rows {
		if i == opt.OffsetRow {
			for j, r := range row {
				for n, col := range opt.Column {
					if col.Label == r {
						hs[j] = opt.Column[n]
					}
				}
			}
			break
		}
	}
	var data []map[string]interface{}
	for i, row := range rows {
		if i <= opt.OffsetRow {
			continue
		}
		item := make(map[string]interface{})
		for j, r := range row {
			if h, ok := hs[j]; ok {
				if h.Type == "int" {
					s, err := strconv.Atoi(r)
					if err != nil {
						return err
					}
					item[h.Prop] = s
				} else {
					item[h.Prop] = r
				}
			}
		}
		data = append(data, item)
	}
	mapstructure.Decode(data, res)
	return nil

}

func GetImportTplHeader(fb *filebase.Client, fileId string, opt *ImportOption) ([]string, error) {
	fs, err := fb.Get(fileId)
	if err != nil {
		return nil, err
	}
	defer fs.Close()
	f, err := excelize.OpenReader(fs)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	rows, err := f.GetRows(opt.Sheet)
	if err != nil {
		return nil, err
	}
	for i, row := range rows {
		if i == opt.OffsetRow {
			return row, nil
		}
	}
	return nil, err
}

func GetImportHeaders(v reflect.Type, hs *[]ImportColumn) {
	for n := 0; n < v.NumField(); n++ {
		f := v.Field(n)
		if f.Type.Kind() == reflect.Struct {
			GetImportHeaders(f.Type, hs)
		}
		importTag := f.Tag.Get("import")
		if importTag == "" {
			continue
		}
		jsonTag := f.Tag.Get("json")
		if jsonTag == "" {
			continue
		}
		items := strings.Split(importTag, ";")
		h := ImportColumn{}
		fns := make(map[string]string)
		for _, item := range items {
			keys := strings.Split(item, ":")
			if len(keys) > 1 {
				fns[keys[0]] = keys[1]
				continue
			}
			h = ImportColumn{Label: item, Prop: jsonTag}
		}
		for k, fn := range fns {
			if k == "type" {
				h.Type = fn
			}
		}
		*hs = append(*hs, h)

	}
}
