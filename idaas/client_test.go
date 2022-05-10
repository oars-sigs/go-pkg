package idaas

import (
	"encoding/json"
	"testing"
)

var cfg = &Config{
	Address:   "https://idaas.oars-vm.hashwing.cn",
	AppID:     "b197a05e-3e4b-4ec4-b1de-c53915a22dca",
	AppSecret: "tu9iVZ6Hg5AkUjTAMXb6F8uBZZpuYXou",
}

func TestDepts(t *testing.T) {
	c := New(cfg)
	depts, err := c.Depts("", false, false)
	if err != nil {
		t.Error(err)
		return
	}
	data, _ := json.Marshal(depts)
	t.Log(string(data))
}

func TestUsers(t *testing.T) {
	c := New(cfg)
	users, err := c.Users([]string{"test"}, false)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(users)
}
