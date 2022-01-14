package filebase

import (
	"net/url"
	"testing"

	"pkg.oars.vip/go-pkg/constant"
)

func TestSign(t *testing.T) {
	u, _ := url.Parse("https://filebase.oars-vm.hashwing.cn/filebase/api/v1/app/files?X-Oars-Auth-Kind=Hmac&digest=d41d8cd98f00b204e9800998ecf8427e&expireTime=1642147450&name=test.txt&parent=&signature=83b9d9ee368e519593ff8b2f1ac640619438f8a2d299cf7ab5fe238dba31e3c1&size=4&type=txt")
	qs := u.Query()
	s := SignURL(u, "hdhS3MYxpiKCssJX06LuzcGc9vFjbPLQ")
	t.Log(s)
	qs.Set(constant.SignatureKey, s)
	u.RawQuery = qs.Encode()
	t.Log(u.String())
}
