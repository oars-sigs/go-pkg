package filebase

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestClient(t *testing.T) {
	cfg := &Config{
		Address:   "https://filebase.oars-vm.hashwing.cn",
		AppID:     "97593a90-cff9-4040-a2f9-a7c56f69fb81",
		AppSecret: "hdhS3MYxpiKCssJX06LuzcGc9vFjbPLQ",
	}
	cli := New(cfg)
	data := strings.NewReader("test")
	res, err := cli.Put(data, "", "test.txt", "txt", 4)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res.ID, res.Digest)
	ustr, err := cli.PutURL("", "test.txt", "txt", res.Digest, int64(4), 3600)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ustr)

	resp, err := cli.Get(res.ID)
	if err != nil {
		t.Error(err)
		return
	}
	defer resp.Close()
	rdata, err := ioutil.ReadAll(resp)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(rdata))
	durl, err := cli.GetURL(res.ID, 3600)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(durl)
}
