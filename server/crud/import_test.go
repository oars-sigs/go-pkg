package crud

import (
	"fmt"
	"reflect"
	"testing"
)

func TestImportHeader(t *testing.T) {
	var s struct {
		Name string `json:"name" import:"姓名"`
		Age  int    `json:"age" import:"年龄;type:int"`
	}

	hs := make([]ImportColumn, 0)

	GetImportHeaders(reflect.TypeOf(s), &hs)
	fmt.Println(hs)

}
