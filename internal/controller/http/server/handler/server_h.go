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
	"github.com/rs/zerolog/log"
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
	UpdateMetricsJSON(ctx *gin.Context)
	BulkUpdateJSON(ctx *gin.Context)
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
	var input entity.Metrics
	body := ctx.Request.Body
	defer body.Close()
	err := json.NewDecoder(body).Decode(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric"})
		return
	}
	log.Debug().Interface("Request ValueJson Input: %s", input)
	err = getPreCheck(&input)
	if err != nil {
		handleCustomError(ctx, err)
		return
	}
	output := h.storage.Get(&input)
	if output == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": entity.ErrMetricNotFound})
		return
	}
	output.CalculateHash(config.GetConfig().Key)
	log.Debug().Interface("Request ValueJson Output: %s", output)
	ctx.JSON(http.StatusOK, output)
}
func (h *handler) Value(ctx *gin.Context) {
	input := entity.Metrics{
		ID:    ctx.Param("metric_name"),
		MType: ctx.Param("metric_type"),
	}
	err := getPreCheck(&input)
	if err != nil {
		handleCustomError(ctx, err)
		return
	}
	output := h.storage.Get(&input)
	if output == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": entity.ErrMetricNotFound})
		return
	}
	if output.Value != nil {
		ctx.String(http.StatusOK, "%s",
			strconv.FormatFloat(
				*output.Value, 'f',
				h.storage.GetFltPrc(input.ID),
				64))
		return
	} else if output.Delta != nil {
		ctx.String(http.StatusOK, "%d", *output.Delta)
		return
	}

}
func (h *handler) UpdateMetricsJSON(ctx *gin.Context) {
	var input entity.Metrics
	err := json.NewDecoder(ctx.Request.Body).Decode(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": entity.ErrInvalidMetric})
		return
	}
	err = setPreCheck(&input)
	if err != nil {
		handleCustomError(ctx, err)
		return
	}
	output := h.storage.Set(&input)
	if output == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": entity.ErrNameTypeMismatch})
		return
	}
	if config.GetConfig().Server.StoreInterval == 0 || config.GetConfig().Database.Address != "" {
		h.storage.Dump()
	}
	output.CalculateHash(config.GetConfig().Key)
	ctx.JSON(http.StatusOK, output)
}
func (h *handler) UpdateMetric(ctx *gin.Context) {
	input := entity.Metrics{
		ID:    ctx.Param("metric_name"),
		MType: ctx.Param("metric_type"),
	}
	metricValue := ctx.Param("metric_value")
	switch input.MType {
	case entity.GaugeType:
		val, err_ := strconv.ParseFloat(metricValue, 64)
		if err_ != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric value, should be a number"})
			return
		}
		input.Value = &val
	case entity.CounterType:
		val, err_ := strconv.ParseInt(metricValue, 10, 64)
		if err_ != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric value, should be a number"})
			return
		}
		input.Delta = &val
	default:
		input.Delta = nil
	}
	err := setPreCheck(&input)
	if err != nil {
		handleCustomError(ctx, err)
		return
	}
	h.storage.SetFltPrc(input.ID, metricValue)
	output := h.storage.Set(&input)
	if output == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": entity.ErrNameTypeMismatch})
		return
	}
	h.storage.SetFltPrc(input.ID, metricValue)
	if config.GetConfig().Server.StoreInterval == 0 || config.GetConfig().Database.Address != "" {
		h.storage.Dump()
	}
	output.CalculateHash(config.GetConfig().Key)
	ctx.JSON(http.StatusOK, output)
}

func (h *handler) BulkUpdateJSON(ctx *gin.Context) {
	var input []*entity.Metrics
	fmt.Println(ctx.Request.Body)
	err := json.NewDecoder(ctx.Request.Body).Decode(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": entity.ErrInvalidMetric})
		return
	}
	//var inputMapper = make(map[string]*entity.Metrics)
	for i := range input {
		err = setPreCheck(input[i])
		if err != nil {
			log.Error().Err(err).Msg("Some of the input metrics are invalid")
			continue
		}
		val := h.storage.Set(input[i])
		if val == nil {
			log.Error().Err(err).Interface("Some of the input metrics are invalid %s", entity.ErrUnableToStore)
			continue
		}
	}
	if config.GetConfig().Server.StoreInterval == 0 || config.GetConfig().Database.Address != "" {
		h.storage.Dump()
	}

	// It would be reasonable to return the updated metrics, but looks like it's not required

	//var output []entity.Metrics
	//for i := range inputMapper {
	//	inputMapper[i].CalculateHash(config.GetConfig().Key)
	//	output = append(output, *inputMapper[i])
	//}

	ctx.Data(http.StatusOK, "application/json; charset=utf-8", []byte("{}"))
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
