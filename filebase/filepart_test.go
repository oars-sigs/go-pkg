package filebase

import (
	"fmt"
	"testing"
)

func TestMultipart(t *testing.T) {
	cfg := &Config{
		Address:   "https://filebase.oars-vm.hashwing.cn",
		AppID:     "97593a90-cff9-4040-a2f9-a7c56f69fb81",
		AppSecret: "hdhS3MYxpiKCssJX06LuzcGc9vFjbPLQ",
	}
	cli := New(cfg)

	filePath := "D:/vm/iso/TinyCore-10.0.iso"
	fi, err := cli.FileInfo(filePath)
	if err != nil {
		t.Error(err)
		return
	}
	res, err := cli.CreateMultipart("app", "root", fi.MD5)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(*res.Data)
	fm := &FileMetadata{
		Namespace: "app",
		Parent:    "root",
		Name:      fi.Name,
		Size:      fi.Size,
		Kind:      FileKind,
		Digest:    fi.MD5,
		DirPath:   "测试/a/b",
	}
	if res.Data.Status != ExistStatus {
		num, err := cli.PutMultipart(filePath, "app", res.Data.ID, 0)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(num)

		err = cli.MergeMultiPart("app", res.Data.ID, fm)
		if err != nil {
			t.Error(err)
			return
		}
	}
	nf, err := cli.CreateFile(fm)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(*nf)
}
