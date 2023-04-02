package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gynshu-one/go-metric-collector/internal/configs"
	"github.com/gynshu-one/go-metric-collector/internal/storage"
	"github.com/gynshu-one/go-metric-collector/internal/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func setupRouter() (*gin.Engine, *ServerHandler) {
	configs.CFG.Restore = false
	configs.CFG.StoreInterval = 0
	configs.CFG.StoreFile = "/tmp/test_metrics.json"
	configs.CFG.Address = "localhost:8080"
	// Then init files
	configs.CFG.InitFiles()
	gin.SetMode(gin.TestMode)
	h := NewServerHandler()
	r := gin.Default()
	r.GET("/live", h.Live)
	r.GET("/value/:metric_type/:metric_name", h.Value)
	r.POST("/value/", h.ValueJSON)
	r.POST("/update/", h.UpdateMetricsJSON)
	r.POST("/update/:metric_type/:metric_name/:metric_value", h.UpdateMetrics)
	r.GET("/html_all_metrics", h.HTMLAllMetrics)

	return r, h
}

var (
	router, serverHandler = setupRouter()
)

func TestNewServerHandler(t *testing.T) {
	handler := NewServerHandler()
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.Memory)
}

func TestLive(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/live", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestValueJSON(t *testing.T) {
	metric := storage.Metrics{
		ID:    "TestMetric",
		MType: storage.GaugeType,
		Value: tools.Float64Ptr(55.0),
	}
	serverHandler.Memory.Set(&metric)

	metricJSON, err := json.Marshal(metric)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/value/", bytes.NewBuffer(metricJSON))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var respMetric storage.Metrics
	err = json.NewDecoder(resp.Body).Decode(&respMetric)
	require.NoError(t, err)
	assert.Equal(t, metric.ID, respMetric.ID)
	assert.Equal(t, metric.MType, respMetric.MType)
	assert.Equal(t, *metric.Value, *respMetric.Value)
}

func TestValue(t *testing.T) {

	// Set a value first
	metric := storage.Metrics{
		ID:    "TestMetric",
		MType: storage.GaugeType,
		Value: tools.Float64Ptr(55.0),
	}
	serverHandler.Memory.Set(&metric)

	req := httptest.NewRequest(http.MethodGet, "/value/gauge/TestMetric", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "55.000", resp.Body.String())
}

func TestUpdateMetricsJSON(t *testing.T) {
	TesCases := []struct {
		name   string
		arg    *storage.Metrics
		status int
	}{
		{
			name: "TestGauge",
			arg: &storage.Metrics{
				ID:    "TestGauge",
				MType: storage.GaugeType,
				Value: tools.Float64Ptr(55.0),
			},
			status: http.StatusOK,
		},
		{
			name: "TestCounter",
			arg: &storage.Metrics{
				ID:    "TestCounter",
				MType: storage.CounterType,
				Delta: tools.Int64Ptr(55),
			},
			status: http.StatusOK,
		},
		{
			name: "InvalidMetricType",
			arg: &storage.Metrics{
				ID:    "InvalidMetricType",
				MType: "aa",
				Value: tools.Float64Ptr(55.0),
			},
			status: http.StatusNotImplemented,
		},
		{
			name: "InvalidTypeAndValue",
			arg: &storage.Metrics{
				ID:    "InvalidTypeAndValue",
				MType: storage.CounterType,
				Value: tools.Float64Ptr(55.0),
			},
			status: http.StatusBadRequest,
		},
		{
			name: "InvalidTypeAndValue",
			arg: &storage.Metrics{
				ID:    "InvalidTypeAndValue",
				MType: storage.GaugeType,
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

			updatedMetric := serverHandler.Memory.Get(tc.arg)
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
		arg    *storage.Metrics
		status int
	}{
		{
			name: "TestGauge",
			arg: &storage.Metrics{
				ID:    "TestGauge2",
				MType: storage.GaugeType,
				Value: tools.Float64Ptr(55.0),
			},
			status: http.StatusOK,
		},
		{
			name: "TestCounter",
			arg: &storage.Metrics{
				ID:    "TestCounter2",
				MType: storage.CounterType,
				Delta: tools.Int64Ptr(55),
			},
			status: http.StatusOK,
		},
		{
			name: "InvalidMetricType",
			arg: &storage.Metrics{
				ID:    "InvalidMetricType2",
				MType: "aa",
				Value: tools.Float64Ptr(55.0),
			},
			status: http.StatusNotImplemented,
		},
		{
			name: "InvalidTypeAndValue",
			arg: &storage.Metrics{
				ID:    "InvalidTypeAndValue2",
				MType: storage.CounterType,
				Value: tools.Float64Ptr(55.0),
			},
			status: http.StatusBadRequest,
		},
	}
	for _, tc := range TesCases {
		t.Run(tc.name, func(t *testing.T) {
			val := ""
			if tc.arg.Value != nil {
				val = strconv.FormatFloat(*tc.arg.Value, 'f', 3, 64)
			}
			if tc.arg.Delta != nil {
				val = strconv.FormatInt(*tc.arg.Delta, 10)
			}
			url := fmt.Sprintf("/update/%s/%s/%s", tc.arg.MType, tc.arg.ID, val)
			req := httptest.NewRequest(http.MethodPost, url, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)
			assert.Equal(t, tc.status, resp.Code)

			updatedMetric := serverHandler.Memory.Get(tc.arg)
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
	metric1 := storage.Metrics{
		ID:    "TestMetric1",
		MType: storage.GaugeType,
		Value: tools.Float64Ptr(55.0),
	}
	serverHandler.Memory.Set(&metric1)

	metric2 := storage.Metrics{
		ID:    "TestMetric2",
		MType: storage.CounterType,
		Delta: tools.Int64Ptr(1),
	}
	serverHandler.Memory.Set(&metric2)

	req := httptest.NewRequest(http.MethodGet, "/html_all_metrics", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	htmlResponse := resp.Body.String()
	assert.Contains(t, htmlResponse, "<th>Type</th><th>Name</th><th>Value</th>")
	assert.Contains(t, htmlResponse, "<td>gauge</td><td>TestMetric1</td><td>55.000000</td>")
	assert.Contains(t, htmlResponse, "<td>counter</td><td>TestMetric2</td><td>1</td>")
}
