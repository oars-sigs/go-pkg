package idaas

import (
	"encoding/json"
	"testing"
)

var cfg = &Config{
	Address:   "https://idaas.oars.gzsunrun.cn",
	AppID:     "871ac2e3-9451-48ab-8cae-74c35120a419",
	AppSecret: "Tl0bDXIrL6z6cOn53VgOa2XoxlqIk2B0",
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
