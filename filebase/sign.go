package filebase

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/url"
	"sort"

	"pkg.oars.vip/go-pkg/constant"
)

func SignURL(u *url.URL, signKey string) string {
	qs := u.Query()
	keys := []string{}
	for k := range qs {
		if k != constant.SignatureKey {
			keys = append(keys, k)
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[j] > keys[i]
	})
	qstr := ""
	for _, k := range keys {
		qstr += k + qs.Get(k)
	}
	pstr := u.Path
	h := hmac.New(sha256.New, []byte(signKey))
	h.Write([]byte(pstr))
	h.Write([]byte(qstr))
	b := h.Sum(nil)
	return fmt.Sprintf("%x", b)
}
