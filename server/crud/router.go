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
	r.POST("/:resource/import", crudc.Import)
	r.POST("/:resource/createinbatches", crudc.CreateInBatches)
	r.POST("/:resource/former/:id/:mark", crudc.CreateFormer)
	//
	r.GET("/presource/:presource/:pid/:resource/:id", crudc.Get)
	r.GET("/presource/:presource/:pid/:resource/:id/:action", crudc.Get)
	r.GET("/presource/:presource/:pid/:resource", crudc.List)
	r.POST("/presource/:presource/:pid/:resource", crudc.Create)
	r.PUT("/presource/:presource/:pid/:resource", crudc.Put)
	r.PUT("/presource/:presource/:pid/:resource/:id", crudc.Update)
	r.PUT("/presource/:presource/:pid/:resource/:id/:action", crudc.Update)
	r.DELETE("/presource/:presource/:pid/:resource/:id", crudc.Delete)
	r.GET("/presource/:presource/:pid/:resource/export", crudc.Export)
	r.POST("/presource/:presource/:pid/:resource/import", crudc.Import)
	r.POST("/presource/:presource/:pid/:resource/createinbatches", crudc.CreateInBatches)
	r.POST("/presource/:presource/:pid/:resource/former/:id/:mark", crudc.CreateFormer)
	if crudc.opt.Former != nil {
		r.POST("/flow/hook", crudc.opt.Former.Hook(crudc.FlowHook))
	}

}
