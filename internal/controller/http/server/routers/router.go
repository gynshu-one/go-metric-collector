package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/gynshu-one/go-metric-collector/internal/controller/http/server/handler"
)

func MetricsRoute(router *gin.Engine, handler handler.Handler) {
	router.GET("/", handler.HTMLAllMetrics)
	router.GET("/live/", handler.Live)

	router.POST("/value/", handler.ValueJSON)
	router.POST("/update/", handler.UpdateMetricJSON)
	router.POST("/updates/", handler.UpdateMetricsJSON)

	router.GET("/value/:metric_type/:metric_name", handler.Value)
	router.POST("/update/:metric_type/:metric_name/:metric_value", handler.UpdateMetric)

	router.GET("/ping", handler.PingDB)
}
