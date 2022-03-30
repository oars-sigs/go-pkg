package dist

import (
	"embed"
	"fmt"
	"mime"
	"strings"

	"github.com/gin-gonic/gin"
)

func Dist(dir, contextPath string, static embed.FS) func(*gin.Context) {
	return func(g *gin.Context) {
		path := g.Request.URL.Path
		s := strings.Split(path, ".")
		contentType := mime.TypeByExtension(fmt.Sprintf(".%s", s[len(s)-1]))
		data, err := static.ReadFile(dir + path)
		if err != nil {
			data, err = static.ReadFile(dir + contextPath + "/index.html")
			contentType = mime.TypeByExtension(".html")
			if err != nil {
				g.AbortWithStatus(404)
			}
		}
		g.Data(200, contentType, data)
	}
}
