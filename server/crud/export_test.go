package crud

import "testing"

func TestExportExcel(t *testing.T) {
	var sheets []ExportOption
	sheets = append(sheets, ExportOption{
		Name: "test",
		Data: []map[string]any{
			{
				"test": 1,
			},
		},
		Headers: []HeaderItem{
			{
				Title: "测试标题",
				Filed: "test",
			},
		},
	})
	fs := ExportExcel(sheets)
	fs.SaveAs("test.xlsx")
}
