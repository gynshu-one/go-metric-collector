package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/gynshu-one/go-metric-collector/internal/handlers"
)

func MetricsRoute(router *gin.Engine, handler *handlers.ServerHandler) {
	router.GET("/", handler.HTMLAllMetrics)
	router.GET("/live/", handler.Live)

	router.POST("/value/", handler.ValueJSON)
	router.POST("/update/", handler.UpdateMetricsJSON)

	router.GET("/value/:metric_type/:metric_name", handler.Value)
	router.POST("/update/:metric_type/:metric_name/:metric_value", handler.UpdateMetrics)

	router.GET("/ping", handler.PingDb)
}
