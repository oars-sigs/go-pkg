package static

import (
	"compress/gzip"
	"embed"
	"fmt"
	"io/ioutil"
	"mime"
	"strings"

	"github.com/gin-gonic/gin"
)

func WrapHandler(dir, contextPath string, static embed.FS) func(*gin.Context) {
	return func(g *gin.Context) {
		path := g.Request.URL.Path
		s := strings.Split(path, ".")
		contentType := mime.TypeByExtension(fmt.Sprintf(".%s", s[len(s)-1]))
		fs, err := static.Open(dir + path)
		if err != nil {
			fs, err = static.Open(dir + contextPath + "/index.html")
			contentType = mime.TypeByExtension(".html")
			if err != nil {
				g.AbortWithStatus(404)
			}
		}
		if !strings.Contains(g.GetHeader("Accept-Encoding"), "gzip") {
			data, err := ioutil.ReadAll(fs)
			if err != nil {
				g.AbortWithStatus(500)
				return
			}
			g.Data(200, contentType, data)
			return
		}
		gr, err := gzip.NewReader(fs)
		if err != nil {
			g.AbortWithStatus(500)
			return
		}
		data, err := ioutil.ReadAll(gr)
		if err != nil {
			g.AbortWithStatus(500)
			return
		}
		g.Header("Content-Encoding", "gzip")
		g.Data(200, contentType, data)
	}
}
