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
	db      postgres.DB
}

type Handler interface {
	Live(ctx *gin.Context)
	ValueJSON(ctx *gin.Context)
	Value(ctx *gin.Context)
	UpdateMetricsJSON(ctx *gin.Context)
	UpdateMetrics(ctx *gin.Context)
	HTMLAllMetrics(ctx *gin.Context)
	PingDB(ctx *gin.Context)
}

func NewServerHandler(storage storage.ServerStorage, db postgres.DB) *handler {
	hand := &handler{
		storage: storage,
		db:      db,
	}
	if config.GetConfig().Server.Restore {
		hand.storage.Restore()
	}
	if config.GetConfig().Server.StoreInterval != 0 {
		ticker := time.NewTicker(config.GetConfig().Server.StoreInterval)
		go func() {
			for {
				t := <-ticker.C
				hand.storage.Dump()
				fmt.Println("Saved to file at", t)
			}
		}()
	}
	return hand
}
func (h handler) Live(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "ServerStorage is live"})
}
func (h handler) ValueJSON(ctx *gin.Context) {
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
func (h handler) Value(ctx *gin.Context) {
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
		floatVal := *val.Value
		floatStr := strconv.FormatFloat(floatVal, 'f', 3, 64)
		ctx.Data(http.StatusOK, "text/plain", []byte(floatStr))
	} else if val.Delta != nil {
		intDelta := *val.Delta
		intStr := strconv.FormatInt(intDelta, 10)
		ctx.Data(http.StatusOK, "text/plain", []byte(intStr))
	}

}
func (h handler) UpdateMetricsJSON(ctx *gin.Context) {
	var m entity.Metrics
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
	val := h.storage.Set(&m)
	if val == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": entity.NameTypeMismatch})
		return
	}
	if config.GetConfig().Server.StoreInterval == 0 {
		go h.storage.Dump()
	}
	ctx.JSON(http.StatusOK, val)
}

func (h handler) UpdateMetrics(ctx *gin.Context) {
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
	if config.GetConfig().Server.StoreInterval == 0 {
		go h.storage.Dump()
	}
	ctx.JSON(http.StatusOK, val)
}

func (h handler) HTMLAllMetrics(ctx *gin.Context) {
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

func (h handler) PingDB(ctx *gin.Context) {
	c := context.Background()
	err := h.db.Ping(c)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "DB is live"})
}
