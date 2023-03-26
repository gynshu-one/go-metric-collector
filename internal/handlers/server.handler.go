package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gynshu-one/go-metric-collector/internal/storage"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func Live(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Server is live"})
}
func ValueJson(ctx *gin.Context) {
	// must get storage.Metrics is json body of request
	var m storage.Metrics
	body := ctx.Request.Body
	defer body.Close()
	err := json.NewDecoder(body).Decode(&m)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid metric"})
		return
	}
	m.MType = strings.ToLower(m.MType)
	ok := storage.Memory.CheckMetricType(m.MType)
	if !ok {
		ctx.JSON(http.StatusNotImplemented, gin.H{"error": "Invalid metric type"})
		return
	}
	// Add metric to storage
	newM, ok := storage.Memory.FindMetricByName(m.ID)
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Metric not found"})
		return
	}
	if !storage.Memory.ValidateTypeAndValue(newM) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Metric type and value mismatch"})
		return
	}
	ctx.JSON(http.StatusOK, newM)
}
func UpdateMetricsJson(ctx *gin.Context) {
	var m storage.Metrics
	body := ctx.Request.Body
	defer body.Close()
	err := json.NewDecoder(body).Decode(&m)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric"})
		return
	}
	m.MType = strings.ToLower(m.MType)
	ok := storage.Memory.CheckMetricType(m.MType)
	if !ok {
		ctx.JSON(http.StatusNotImplemented, gin.H{"error": "Invalid metric type"})
		return
	}
	// check if metric type and value are valid
	if !storage.Memory.ValidateTypeAndValue(m) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Metric type and value mismatch"})
		return
	}
	// check if metric value is valid
	if !storage.Memory.ValidateValue(m) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric value"})
		return
	}
	// Add metric to storage
	newM := storage.Memory.UpdateMetric(m)
	ctx.JSON(http.StatusOK, newM)
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

func Value(ctx *gin.Context) {
	metricType := ctx.Param("metric_type")
	metricName := ctx.Param("metric_name")

	// This is because the metric finder is case-sensitive
	metricType = strings.ToLower(metricType)

	m, ok := storage.Memory.FindMetricByName(metricName)
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Type or name not found"})
		return
	} else {
		if !storage.Memory.ValidateTypeAndValue(m) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Metric type and value mismatch"})
			return
		}
		if m.Value != nil {
			floatVal := *m.Value
			floatStr := strconv.FormatFloat(floatVal, 'f', -1, 64)
			ctx.Data(http.StatusOK, "text/plain", []byte(floatStr))
		} else if m.Delta != nil {
			intDelta := *m.Delta
			intStr := strconv.FormatInt(intDelta, 10)
			ctx.Data(http.StatusOK, "text/plain", []byte(intStr))
		}

	}
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

	metricType = strings.ToLower(metricType)
	if !storage.Memory.CheckMetricType(metricType) {
		ctx.JSON(http.StatusNotImplemented, gin.H{"error": "Invalid metric type"})
		return
	}
	m := storage.Metrics{
		MType: metricType,
		ID:    metricName,
	}

	switch metricType {
	case "counter":
		metricValueInt, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric value, should be an integer"})
			return
		}
		m.Delta = &metricValueInt

	case "gauge":
		metricValueFloat, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric value, should be a number"})
			return
		}
		m.Value = &metricValueFloat
	}
	newM := storage.Memory.UpdateMetric(m)
	ctx.JSON(http.StatusOK, newM)
}
