package filebase

import (
	"testing"
)

func TestClient(t *testing.T) {
	cfg := &Config{
		Address:   "https://filebase.oars-vm.hashwing.cn",
		AppID:     "97593a90-cff9-4040-a2f9-a7c56f69fb81",
		AppSecret: "hdhS3MYxpiKCssJX06LuzcGc9vFjbPLQ",
	}
	cli := New(cfg)

	path := "D:/vm/iso/ttylinux-pc_i686-16.1.iso"
	res, err := cli.FPut(path, "", "fedora-coreos-33.20210412.3.0-live.x86_64.iso", "iso")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res.ID, res.Digest)
	ustr, err := cli.PutURL("", "test.txt", "txt", res.Digest, res.Size, 3600)
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
	// rdata, err := ioutil.ReadAll(resp)
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// t.Log(string(rdata))
	durl, err := cli.GetURL(res.ID, 3600)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(durl)
}
