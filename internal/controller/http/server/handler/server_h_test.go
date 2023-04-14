package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gynshu-one/go-metric-collector/internal/domain/entity"
	"github.com/gynshu-one/go-metric-collector/internal/domain/service"
	usecase "github.com/gynshu-one/go-metric-collector/internal/domain/usecase/storage"
	"github.com/gynshu-one/go-metric-collector/internal/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
)

func setupRouter() (*gin.Engine, *handler) {
	// Then init files
	gin.SetMode(gin.TestMode)
	h := NewServerHandler(usecase.NewServerUseCase(service.NewMemService(&sync.Map{}), nil), nil)
	r := gin.Default()
	r.GET("/live", h.Live)
	r.GET("/value/:metric_type/:metric_name", h.Value)
	r.POST("/value/", h.ValueJSON)
	r.POST("/update/", h.UpdateMetricsJSON)
	r.POST("/update/:metric_type/:metric_name/:metric_value", h.UpdateMetric)
	r.GET("/html_all_metrics", h.HTMLAllMetrics)

	return r, h
}

var (
	router, serverHandler = setupRouter()
)

func TestNewServerHandler(t *testing.T) {
	h := NewServerHandler(usecase.NewServerUseCase(service.NewMemService(&sync.Map{}), nil), nil)
	assert.NotNil(t, h)
	assert.NotNil(t, h.storage)
}

func TestLive(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/live", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestValueJSON(t *testing.T) {
	metric := entity.Metrics{
		ID:    "TestMetric",
		MType: entity.GaugeType,
		Value: tools.Float64Ptr(55.0),
	}
	serverHandler.storage.Set(&metric)

	metricJSON, err := json.Marshal(metric)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/value/", bytes.NewBuffer(metricJSON))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var respMetric entity.Metrics
	err = json.NewDecoder(resp.Body).Decode(&respMetric)
	require.NoError(t, err)
	assert.Equal(t, metric.ID, respMetric.ID)
	assert.Equal(t, metric.MType, respMetric.MType)
	assert.Equal(t, *metric.Value, *respMetric.Value)
}

func TestValue(t *testing.T) {

	// Set a value first
	metric := entity.Metrics{
		ID:    "TestMetric",
		MType: entity.GaugeType,
		Value: tools.Float64Ptr(55),
	}
	serverHandler.storage.Set(&metric)

	req := httptest.NewRequest(http.MethodGet, "/value/gauge/TestMetric", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "55", resp.Body.String())
}

func TestUpdateMetricsJSON(t *testing.T) {
	TesCases := []struct {
		name   string
		arg    *entity.Metrics
		status int
	}{
		{
			name: "TestGauge",
			arg: &entity.Metrics{
				ID:    "TestGauge",
				MType: entity.GaugeType,
				Value: tools.Float64Ptr(55.0),
			},
			status: http.StatusOK,
		},
		{
			name: "TestCounter",
			arg: &entity.Metrics{
				ID:    "TestCounter",
				MType: entity.CounterType,
				Delta: tools.Int64Ptr(55),
			},
			status: http.StatusOK,
		},
		{
			name: "InvalidMetricType",
			arg: &entity.Metrics{
				ID:    "InvalidMetricType",
				MType: "aa",
				Value: tools.Float64Ptr(55.0),
			},
			status: http.StatusNotImplemented,
		},
		{
			name: "InvalidTypeAndValue",
			arg: &entity.Metrics{
				ID:    "InvalidTypeAndValue",
				MType: entity.CounterType,
				Value: tools.Float64Ptr(55.0),
			},
			status: http.StatusBadRequest,
		},
		{
			name: "InvalidTypeAndValue",
			arg: &entity.Metrics{
				ID:    "InvalidTypeAndValue",
				MType: entity.GaugeType,
				Delta: tools.Int64Ptr(55),
			},
			status: http.StatusBadRequest,
		},
	}
	for _, tc := range TesCases {
		t.Run(tc.name, func(t *testing.T) {
			metricJSON, err := json.Marshal(tc.arg)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewBuffer(metricJSON))
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)
			assert.Equal(t, tc.status, resp.Code)

			updatedMetric := serverHandler.storage.Get(tc.arg)
			if tc.status != http.StatusOK {
				assert.Nil(t, updatedMetric)
				return
			} else {
				require.NotNil(t, updatedMetric)
				assert.Equal(t, *tc.arg, *updatedMetric)
			}
		})
	}
}

func TestUpdateMetrics(t *testing.T) {
	TesCases := []struct {
		name   string
		arg    *entity.Metrics
		status int
	}{
		{
			name: "TestGauge",
			arg: &entity.Metrics{
				ID:    "TestGauge2",
				MType: entity.GaugeType,
				Value: tools.Float64Ptr(55.000),
			},
			status: http.StatusOK,
		},
		{
			name: "TestCounter",
			arg: &entity.Metrics{
				ID:    "TestCounter2",
				MType: entity.CounterType,
				Delta: tools.Int64Ptr(55),
			},
			status: http.StatusOK,
		},
		{
			name: "InvalidMetricType",
			arg: &entity.Metrics{
				ID:    "InvalidMetricType2",
				MType: "aa",
				Value: tools.Float64Ptr(55.000),
			},
			status: http.StatusNotImplemented,
		},
		{
			name: "InvalidTypeAndValue",
			arg: &entity.Metrics{
				ID:    "InvalidTypeAndValue2",
				MType: entity.CounterType,
				Value: tools.Float64Ptr(55.000),
			},
			status: http.StatusBadRequest,
		},
	}
	for _, tc := range TesCases {
		t.Run(tc.name, func(t *testing.T) {
			val := ""
			if tc.arg.Delta != nil {
				val = strconv.FormatInt(*tc.arg.Delta, 10)
			} else {
				val = strconv.FormatFloat(*tc.arg.Value, 'f', 3, 64)
			}
			url := fmt.Sprintf("/update/%s/%s/%s", tc.arg.MType, tc.arg.ID, val)
			req := httptest.NewRequest(http.MethodPost, url, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)
			assert.Equal(t, tc.status, resp.Code)

			updatedMetric := serverHandler.storage.Get(tc.arg)
			if tc.status != http.StatusOK {
				assert.Nil(t, updatedMetric)
				return
			} else {
				require.NotNil(t, updatedMetric)
				assert.Equal(t, *tc.arg, *updatedMetric)
			}
		})
	}
}

func TestHTMLAllMetrics(t *testing.T) {
	// Set some values first
	metric1 := entity.Metrics{
		ID:    "TestMetric1",
		MType: entity.GaugeType,
		Value: tools.Float64Ptr(55.0),
	}
	serverHandler.storage.Set(&metric1)

	metric2 := entity.Metrics{
		ID:    "TestMetric2",
		MType: entity.CounterType,
		Delta: tools.Int64Ptr(1),
	}
	serverHandler.storage.Set(&metric2)

	req := httptest.NewRequest(http.MethodGet, "/html_all_metrics", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	htmlResponse := resp.Body.String()
	assert.Contains(t, htmlResponse, "<th>Type</th><th>Name</th><th>Value</th>")
	assert.Contains(t, htmlResponse, "<td>gauge</td><td>TestMetric1</td><td>55.000000</td>")
	assert.Contains(t, htmlResponse, "<td>counter</td><td>TestMetric2</td><td>1</td>")
}
