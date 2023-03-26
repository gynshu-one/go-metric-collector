package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/gynshu-one/go-metric-collector/internal/handlers"
)

func MetricsRoute(router *gin.Engine) {
	router.GET("/", handlers.HTMLAllMetrics)
	router.GET("/live/", handlers.Live)

	router.POST("/value/", handlers.ValueJson)
	router.POST("/update/", handlers.UpdateMetricsJson)

	router.GET("/value/:metric_type/:metric_name", handlers.Value)
	router.POST("/update/:metric_type/:metric_name/:metric_value", handlers.UpdateMetrics)
}
