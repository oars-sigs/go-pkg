package filebase

import (
	"net/url"
	"testing"

	"pkg.oars.vip/go-pkg/constant"
)

func TestSign(t *testing.T) {
	u, _ := url.Parse("https://filebase.oars-vm.hashwing.cn/filebase/api/v1/app/files?X-Oars-App-Id=97593a90-cff9-4040-a2f9-a7c56f69fb81&X-Oars-Auth-Kind=Hmac&digest=13105dff25ba2a06892e9f5f4061ce13&expireTime=1642420762&name=test.txt&parent=")
	qs := u.Query()
	s := SignURL(u, "hdhS3MYxpiKCssJX06LuzcGc9vFjbPLQ")
	t.Log(s)
	qs.Set(constant.SignatureKey, s)
	u.RawQuery = qs.Encode()
	t.Log(u.String())
}
