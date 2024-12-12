package crud

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(d, &res)
	if err != nil {
		return err
	}
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
