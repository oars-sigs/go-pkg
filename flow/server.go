package flow

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type ServerAction struct {
	Port   int                    `yaml:"port"`
	Routes []ServerRoute          `yaml:"routes"`
	Values map[string]interface{} `yaml:"values"`
}

type ServerRoute struct {
	Method   string `yaml:"method"`
	Path     string `yaml:"path"`
	Playbook string `yaml:"playbook"`
	Tasks    []Task `yaml:"tasks"`
	File     string `yaml:"file"`
}

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method

		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
	}
}

func (a *ServerAction) Do(conf *Config, params interface{}) (interface{}, error) {
	p := params.(ServerAction)
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(Cors())

	for _, route := range p.Routes {
		h := func(ctx *gin.Context) {
			gctx := make(map[string]interface{})
			qs := make(map[string][]string)
			for k, v := range ctx.Request.URL.Query() {
				qs[k] = v
			}
			gctx["query"] = qs
			body, _ := ioutil.ReadAll(ctx.Request.Body)
			gctx["body"] = body
			gctx["header"] = ctx.Request.Header
			if route.File != "" {
				data, err := ioutil.ReadFile(route.File)
				if err != nil {
					ctx.AbortWithStatus(500)
					return
				}
				var ts []Task
				err = yaml.Unmarshal(data, &ts)
				if err != nil {
					ctx.AbortWithStatus(500)
					return
				}
				route.Tasks = ts
			}
			if p.Values == nil {
				p.Values = make(map[string]interface{})
			}
			gvars := NewGvars(&Vars{
				Ctx:    gctx,
				Values: p.Values,
			})
			err := NewPlaybook(route.Tasks, gvars).Run(conf)
			if err != nil {
				ctx.AbortWithStatus(500)
				return
			}
			resp, ok := gvars.GetVar("ctx.resp")
			if ok {
				ctx.JSON(200, resp)
				return
			}
		}
		switch route.Method {
		case "GET":
			r.GET(route.Path, h)
		case "POST":
			r.POST(route.Path, h)
		case "DELETE":
			r.DELETE(route.Path, h)
		case "PUT":
			r.PUT(route.Path, h)
		case "PATCH":
			r.PATCH(route.Path, h)
		case "ANY":
			r.Any(route.Path, h)
		}
	}
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", p.Port),
		Handler:        r,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println("start server ", fmt.Sprintf(":%d", p.Port))
	err := s.ListenAndServe()
	return nil, err
}

func (a *ServerAction) Params() interface{} {
	return ServerAction{}
}

func (a *ServerAction) Scheme() string {
	return ""
}
