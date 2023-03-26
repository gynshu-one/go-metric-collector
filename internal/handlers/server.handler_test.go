package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gynshu-one/go-metric-collector/internal/storage"
	"github.com/gynshu-one/go-metric-collector/internal/tools"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/", HTMLAllMetrics)
	router.GET("/live/", Live)

	router.POST("/value/", ValueJson)
	router.POST("/update/", UpdateMetricsJson)

	router.GET("/value/:metric_type/:metric_name", Value)
	router.POST("/update/:metric_type/:metric_name/:metric_value", UpdateMetrics)
	return router
}

func TestLive(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/live/", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Server is live")
}

func TestValue(t *testing.T) {
	gaugeMetric := storage.Metrics{
		ID:    "TestGauge",
		MType: "gauge",
		Value: tools.Float64Ptr(2.0),
	}
	storage.Memory.UpdateMetric(gaugeMetric)

	router := setupRouter()
	jsonData, _ := json.Marshal(gaugeMetric)
	req1 := httptest.NewRequest("POST", "/update/", bytes.NewBuffer(jsonData))
	resp1 := httptest.NewRecorder()
	router.ServeHTTP(resp1, req1)
	req, _ := http.NewRequest("GET", "/value/gauge/TestGauge", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "2.000", resp.Body.String())

	req2, _ := http.NewRequest("GET", "/value/gauge/NonExistent", nil)
	resp2 := httptest.NewRecorder()
	router.ServeHTTP(resp2, req2)

	assert.Equal(t, http.StatusNotFound, resp2.Code)
	assert.Contains(t, resp2.Body.String(), "Type or name not found")
}

func TestValueJson(t *testing.T) {
	router := setupRouter()

	gaugeMetric := storage.Metrics{
		ID:    "TestGauge",
		MType: "gauge",
		Value: tools.Float64Ptr(2.0),
	}
	jsonData, _ := json.Marshal(gaugeMetric)
	req1 := httptest.NewRequest("POST", "/update/", bytes.NewBuffer(jsonData))
	resp1 := httptest.NewRecorder()
	router.ServeHTTP(resp1, req1)
	req, _ := http.NewRequest("POST", "/value/", bytes.NewBuffer(jsonData))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "TestGauge")
	assert.Contains(t, resp.Body.String(), "2")

	req2, _ := http.NewRequest("POST", "/value/", bytes.NewBuffer([]byte("invalid_json")))
	resp2 := httptest.NewRecorder()
	router.ServeHTTP(resp2, req2)

	assert.Equal(t, http.StatusNotFound, resp2.Code)
	assert.Contains(t, resp2.Body.String(), "Invalid metric")
}

func TestUpdateMetricsJson(t *testing.T) {
	router := setupRouter()

	gaugeMetric := storage.Metrics{
		ID:    "TestGauge",
		MType: "gauge",
		Value: tools.Float64Ptr(3.0),
	}
	jsonData, _ := json.Marshal(gaugeMetric)
	req, _ := http.NewRequest("POST", "/update/", bytes.NewBuffer(jsonData))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "TestGauge")
	assert.Contains(t, resp.Body.String(), "3")

	req2, _ := http.NewRequest("POST", "/update/", bytes.NewBuffer([]byte("invalid_json")))
	resp2 := httptest.NewRecorder()
	router.ServeHTTP(resp2, req2)

	assert.Equal(t, http.StatusBadRequest, resp2.Code)
	assert.Contains(t, resp2.Body.String(), "Invalid metric")
}

func TestUpdateMetrics(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("POST", "/update/gauge/TestGauge/4", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "TestGauge")
	assert.Contains(t, resp.Body.String(), "4")

	req2, _ := http.NewRequest("POST", "/update/gauge/NewGauge/5", nil)
	resp2 := httptest.NewRecorder()
	router.ServeHTTP(resp2, req2)

	assert.Equal(t, http.StatusOK, resp2.Code)
	assert.Contains(t, resp2.Body.String(), "NewGauge")
	assert.Contains(t, resp2.Body.String(), "5")
}
