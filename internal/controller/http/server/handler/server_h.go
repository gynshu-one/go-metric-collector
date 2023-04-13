package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	config "github.com/gynshu-one/go-metric-collector/internal/config/server"
	"github.com/gynshu-one/go-metric-collector/internal/domain/entity"
	"github.com/gynshu-one/go-metric-collector/internal/domain/usecase/storage"
	"github.com/gynshu-one/go-metric-collector/pkg/client/postgres"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type handler struct {
	storage storage.ServerStorage
	dbConn  postgres.DBConn
}

type Handler interface {
	Live(ctx *gin.Context)
	ValueJSON(ctx *gin.Context)
	Value(ctx *gin.Context)
	UpdateMetricJSON(ctx *gin.Context)
	UpdateMetricsJSON(ctx *gin.Context)
	UpdateMetric(ctx *gin.Context)
	HTMLAllMetrics(ctx *gin.Context)
	PingDB(ctx *gin.Context)
}

func NewServerHandler(storage storage.ServerStorage, db postgres.DBConn) *handler {
	hand := &handler{
		storage: storage,
		dbConn:  db,
	}
	return hand
}
func (h *handler) Live(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "ServerStorage is live"})
}
func (h *handler) ValueJSON(ctx *gin.Context) {
	var m entity.Metrics
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
	val := h.storage.Get(&m)
	if val == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": entity.MetricNotFound})
		return
	}
	if config.GetConfig().Key != "" {
		val.CalculateAndWriteHash(config.GetConfig().Key)
	}
	ctx.JSON(http.StatusOK, val)
}
func (h *handler) Value(ctx *gin.Context) {
	m := entity.Metrics{
		ID:    ctx.Param("metric_name"),
		MType: ctx.Param("metric_type"),
	}
	err := getPreCheck(&m)
	if err != nil {
		handleCustomError(ctx, err)
		return
	}
	val := h.storage.Get(&m)
	if val == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": entity.MetricNotFound})
		return
	}
	if val.Value != nil {
		ctx.String(http.StatusOK, "%s",
			strconv.FormatFloat(
				*val.Value, 'f',
				h.storage.GetFltPrc(m.ID),
				64))
		return
	} else if val.Delta != nil {
		ctx.String(http.StatusOK, "%d", *val.Delta)
		return
	}

}
func (h *handler) UpdateMetric(ctx *gin.Context) {
	m := entity.Metrics{
		ID:    ctx.Param("metric_name"),
		MType: ctx.Param("metric_type"),
	}
	metricValue := ctx.Param("metric_value")
	switch m.MType {
	case entity.GaugeType:
		val, err_ := strconv.ParseFloat(metricValue, 64)
		if err_ != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric value, should be a number"})
			return
		}
		m.Value = &val
	case entity.CounterType:
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
	val := h.storage.Set(&m)
	if val == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": entity.NameTypeMismatch})
		return
	}
	h.storage.SetFltPrc(m.ID, metricValue)
	if config.GetConfig().Server.StoreInterval == 0 || config.GetConfig().Database.Address != "" {
		h.storage.Dump()
	}
	if config.GetConfig().Key != "" {
		val.CalculateAndWriteHash(config.GetConfig().Key)
	}
	ctx.JSON(http.StatusOK, val)
}
func (h *handler) UpdateMetricJSON(ctx *gin.Context) {
	var m entity.Metrics
	err := json.NewDecoder(ctx.Request.Body).Decode(&m)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": entity.InvalidMetric})
		return
	}
	err = setPreCheck(&m)
	if err != nil {
		handleCustomError(ctx, err)
		return
	}
	val := h.storage.Set(&m)
	if val == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": entity.NameTypeMismatch})
		return
	}
	if config.GetConfig().Server.StoreInterval == 0 || config.GetConfig().Database.Address != "" {
		h.storage.Dump()
	}
	if config.GetConfig().Key != "" {
		val.CalculateAndWriteHash(config.GetConfig().Key)
	}
	ctx.JSON(http.StatusOK, val)
}

func (h *handler) UpdateMetricsJSON(ctx *gin.Context) {
	var m []*entity.Metrics
	err := json.NewDecoder(ctx.Request.Body).Decode(&m)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": entity.InvalidMetric})
		return
	}
	if len(m) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": entity.EmptyMetric})
		return
	}
	fmt.Println("\n\nReceived metrics:")
	var mapMetrics = make(map[string]*entity.Metrics)
	for _, metric := range m {
		fmt.Println(metric)
		err = setPreCheck(metric)
		if err != nil {
			handleCustomError(ctx, err)
			return
		}
		val := h.storage.Set(metric)
		if val == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": entity.NameTypeMismatch})
			return
		}
		mapMetrics[metric.ID] = val
	}
	if config.GetConfig().Server.StoreInterval == 0 || config.GetConfig().Database.Address != "" {
		h.storage.Dump()
	}
	var metrics []*entity.Metrics
	fmt.Println("\n\nSending metrics:")
	for _, metric := range mapMetrics {
		fmt.Println(metric)
		if config.GetConfig().Key != "" {
			metric.CalculateAndWriteHash(config.GetConfig().Key)
		}
		metrics = append(metrics, metric)
	}
	ctx.JSON(http.StatusOK, metrics)
}

func (h *handler) HTMLAllMetrics(ctx *gin.Context) {
	body := generateHTMLTable(h.storage)
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

func (h *handler) PingDB(ctx *gin.Context) {
	c, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()
	err := h.dbConn.Ping(c)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "DBConn is live"})
}
