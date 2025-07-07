package crud

import "github.com/gin-gonic/gin"

func AddBaseRouter(r *gin.RouterGroup, crudc *BaseInfoController) {
	r.GET("/:resource/:id", crudc.Get)
	r.GET("/:resource/:id/:action", crudc.Get)
	r.GET("/:resource", crudc.List)
	r.POST("/:resource", crudc.Create)
	r.PUT("/:resource", crudc.Put)
	r.PUT("/:resource/:id", crudc.Update)
	r.PUT("/:resource/:id/:action", crudc.Update)
	r.DELETE("/:resource/:id", crudc.Delete)
	r.GET("/:resource/export", crudc.Export)
}
