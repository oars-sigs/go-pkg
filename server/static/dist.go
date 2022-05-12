package static

import (
	"bytes"
	"compress/gzip"
	"embed"
	"fmt"
	"io/ioutil"
	"mime"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var gzPool sync.Pool

func init() {
	gzPool = sync.Pool{
		New: func() interface{} {
			gz, err := gzip.NewWriterLevel(ioutil.Discard, gzip.DefaultCompression)
			if err != nil {
				panic(err)
			}
			return gz
		},
	}
}

func WrapHandler(dir, contextPath string, static embed.FS, defaultPath ...string) func(*gin.Context) {
	return func(g *gin.Context) {
		defaultContextPath := contextPath
		path := g.Request.URL.Path
		if len(defaultPath) > 0 {
			defaultContextPath = defaultPath[0]
			if defaultContextPath != contextPath {
				path = defaultContextPath + strings.TrimPrefix(path, contextPath)
			}
		}
		s := strings.Split(path, ".")
		contentType := mime.TypeByExtension(fmt.Sprintf(".%s", s[len(s)-1]))
		sdata, err := static.ReadFile(dir + path)
		if err != nil {
			sdata, err = static.ReadFile(dir + defaultContextPath + "/index.html")
			contentType = mime.TypeByExtension(".html")
			if err != nil {
				g.AbortWithStatus(404)
				return
			}
			if defaultContextPath != contextPath {
				sdata = bytes.ReplaceAll(sdata, []byte("<head>"), []byte("<head><script>window.OarsContextPath='"+contextPath+"'</script>"))
				sdata = bytes.ReplaceAll(sdata, []byte("href="+defaultContextPath), []byte("href="+contextPath))
				sdata = bytes.ReplaceAll(sdata, []byte("src="+defaultContextPath), []byte("src="+contextPath))
			}
		}
		if !strings.Contains(g.GetHeader("Accept-Encoding"), "gzip") {
			g.Data(200, contentType, sdata)
			return
		}
		g.Header("Content-Type", contentType)
		g.Header("Content-Encoding", "gzip")
		g.AbortWithStatus(200)
		gz := gzPool.Get().(*gzip.Writer)
		defer gzPool.Put(gz)
		defer gz.Reset(ioutil.Discard)
		gz.Reset(g.Writer)
		gz.Write(sdata)
		g.Writer = &gzipWriter{g.Writer, gz}
		defer func() {
			gz.Close()
			g.Header("Content-Length", fmt.Sprint(g.Writer.Size()))
		}()
	}
}

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipWriter) WriteString(s string) (int, error) {
	g.Header().Del("Content-Length")
	return g.writer.Write([]byte(s))
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	g.Header().Del("Content-Length")
	return g.writer.Write(data)
}

// Fix: https://github.com/mholt/caddy/issues/38
func (g *gzipWriter) WriteHeader(code int) {
	g.Header().Del("Content-Length")
	g.ResponseWriter.WriteHeader(code)
}
