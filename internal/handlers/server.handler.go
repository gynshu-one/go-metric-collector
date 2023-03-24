package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/gynshu-one/go-metric-collector/internal/storage"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"net/http"
	"sort"
)

func Live(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Server is live"})
}
func Value(ctx *gin.Context) {
	metricType := ctx.Param("metric_type")
	metricName := ctx.Param("metric_name")

	// This is because the metric finder is case-sensitive
	cs := cases.Title(language.English)
	metricType = cs.String(metricType)

	value, err := storage.Memory.FindMetricByName(metricType, metricName)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Type or name not found"})
		return
	}
	ctx.Data(http.StatusOK, "text/plain", []byte(value))
}
func UpdateMetrics(ctx *gin.Context) {
	metricType := ctx.Param("metric_type")
	metricName := ctx.Param("metric_name")
	metricValue := ctx.Param("metric_value")

	if metricType == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Please provide metric type"})
		return
	}
	if metricName == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Please provide metric name"})
		return
	}
	if metricValue == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Please provide metric value"})
		return
	}

	cs := cases.Title(language.English)
	metricType = cs.String(metricType)
	if !storage.Memory.CheckMetricType(metricType) {
		ctx.JSON(http.StatusNotImplemented, gin.H{"error": "Invalid metric type"})
		return
	}
	err := storage.Memory.AddMetric(metricType, metricName, metricValue)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric value, should be a number"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Metric updated", "type": metricType, "name": metricName, "value": metricValue})
}

func HTMLAllMetrics(ctx *gin.Context) {
	body := storage.Memory.GenerateHTMLTable()
	// Sort the table by type, name, so it's easier to read when page updates
	sort.Strings(body)
	html := "<html><head><title>Metrics</title></head><body><table><tbody><tr><th>Type</th><th>Name</th><th>Value</th></tr>"
	for _, v := range body {
		html += v
	}
	html += "</tbody></table></body></html>"
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
