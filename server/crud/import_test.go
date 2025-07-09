package crud

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/mitchellh/mapstructure"
)

type user struct {
	Name string `json:"name" import:"姓名"`
	Age  int    `json:"age" import:"年龄;type:int"`
}

func getImport() any {
	return []user{}
}

func TestImportHeader(t *testing.T) {

	var s user

	var res = getImport()
	//u := res.([]user)
	res = []map[string]any{}
	parseJSON(&res)
	fmt.Println(res)

	hs := make([]ImportColumn, 0)

	GetImportHeaders(reflect.TypeOf(s), &hs)
	fmt.Println(hs)

}

func parseJSON(res any) {
	data := make([]map[string]any, 0)
	data = append(data, map[string]any{
		"name": "张三",
	})
	mapstructure.Decode(data, res)
}
