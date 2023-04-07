package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gynshu-one/go-metric-collector/internal/configs"
	"github.com/gynshu-one/go-metric-collector/internal/storage"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ServerHandler struct {
	Memory storage.MemActions
}

func NewServerHandler() *ServerHandler {
	hand := &ServerHandler{
		Memory: storage.NewMemStorage(),
	}
	if configs.CFG.Restore {
		hand.Memory.Restore()
	}
	if configs.CFG.StoreInterval != 0 {
		ticker := time.NewTicker(configs.CFG.StoreInterval)
		go func() {
			for {
				t := <-ticker.C
				hand.Memory.Dump()
				fmt.Println("Saved to file at", t)
			}
		}()
	}
	return hand
}
func (s *ServerHandler) Live(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Server is live"})
}
func (s *ServerHandler) ValueJSON(ctx *gin.Context) {
	var m storage.Metrics
	body := ctx.Request.Body
	defer body.Close()
	err := json.NewDecoder(body).Decode(&m)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric"})
		return
	}
	err = getPreCheck(&m)
	if err != nil {
		handleCustomError(ctx, err)
		return
	}
	val := s.Memory.Get(&m)
	if val == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": storage.MetricNotFound})
		return
	}
	if configs.CFG.Key != "" {
		val.CalculateAndWriteHash()
	}
	ctx.JSON(http.StatusOK, val)
}
func (s *ServerHandler) Value(ctx *gin.Context) {
	m := storage.Metrics{
		ID:    ctx.Param("metric_name"),
		MType: ctx.Param("metric_type"),
	}
	err := getPreCheck(&m)
	if err != nil {
		handleCustomError(ctx, err)
		return
	}
	val := s.Memory.Get(&m)
	if val == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": storage.MetricNotFound})
		return
	}
	if val.Value != nil {
		floatVal := *val.Value
		floatStr := strconv.FormatFloat(floatVal, 'f', 3, 64)
		ctx.Data(http.StatusOK, "text/plain", []byte(floatStr))
	} else if val.Delta != nil {
		intDelta := *val.Delta
		intStr := strconv.FormatInt(intDelta, 10)
		ctx.Data(http.StatusOK, "text/plain", []byte(intStr))
	}

}
func (s *ServerHandler) UpdateMetricsJSON(ctx *gin.Context) {
	var m storage.Metrics
	body := ctx.Request.Body
	defer body.Close()
	err := json.NewDecoder(body).Decode(&m)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric"})
		return
	}
	err = setPreCheck(&m)
	if err != nil {
		handleCustomError(ctx, err)
		return
	}
	val := s.Memory.Set(&m)
	if val == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": storage.NameTypeMismatch})
		return
	}
	if configs.CFG.StoreInterval == 0 {
		go s.Memory.Dump()
	}
	ctx.JSON(http.StatusOK, val)
}

func (s *ServerHandler) UpdateMetrics(ctx *gin.Context) {
	m := storage.Metrics{
		ID:    ctx.Param("metric_name"),
		MType: ctx.Param("metric_type"),
	}
	metricValue := ctx.Param("metric_value")
	switch m.MType {
	case storage.GaugeType:
		val, err_ := strconv.ParseFloat(metricValue, 64)
		if err_ != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric value, should be a number"})
			return
		}
		m.Value = &val
	case storage.CounterType:
		val, err_ := strconv.ParseInt(metricValue, 10, 64)
		if err_ != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric value, should be a number"})
			return
		}
		m.Delta = &val
	default:
		m.Delta = nil
	}
	err := setPreCheck(&m)
	if err != nil {
		handleCustomError(ctx, err)
		return
	}
	val := s.Memory.Set(&m)
	if val == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": storage.NameTypeMismatch})
		return
	}
	if configs.CFG.StoreInterval == 0 {
		go s.Memory.Dump()
	}
	ctx.JSON(http.StatusOK, val)
}

func (s *ServerHandler) HTMLAllMetrics(ctx *gin.Context) {
	body := generateHTMLTable(s.Memory)
	// Sort the table by type, name, so it's easier to read when page updates
	sort.Strings(body)
	var sb strings.Builder
	sb.WriteString("<html><head><title>Metrics</title></head><body><table><tbody><tr><th>Type</th><th>Name</th><th>Value</th></tr>")
	for _, v := range body {
		sb.WriteString(v)
	}
	sb.WriteString("</tbody></table></body></html>")
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(sb.String()))
}
