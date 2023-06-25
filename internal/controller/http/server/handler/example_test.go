package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gynshu-one/go-metric-collector/internal/domain/entity"
	"github.com/gynshu-one/go-metric-collector/internal/tools"
	"net/http"
	"net/http/httptest"
)

func ExampleHandler_ValueJSON() {
	var testMetric = &entity.Metrics{
		ID:    "TestGauge",
		MType: entity.GaugeType,
		Value: tools.Float64Ptr(55.0),
	}

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(testMetric)

	// Set the metric to the storage first
	// as if it was already there
	serverHandler.storage.Set(testMetric)

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodPost, "/value/", body)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Output: 200 {"id":"TestGauge","type":"gauge","value":55}
	fmt.Println(resp.Code, resp.Body.String())
}

func ExampleHandler_Value() {
	var testMetric = &entity.Metrics{
		ID:    "TestGauge",
		MType: entity.GaugeType,
		Value: tools.Float64Ptr(55.0),
	}

	// Set the metric to the storage first
	// as if it was already there
	serverHandler.storage.Set(testMetric)

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodGet, "/value/gauge/TestGauge", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Output: 200 55
	fmt.Println(resp.Code, resp.Body.String())
}

func ExampleHandler_UpdateMetricsJSON() {
	// Create a new Metrics object
	var testMetric = &entity.Metrics{
		ID:    "TestGauge",
		MType: entity.GaugeType,
		Value: tools.Float64Ptr(55.0),
	}

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(testMetric)

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodPost, "/update/", body)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Output: 200 {"id":"TestGauge","type":"gauge","value":55}
	fmt.Println(resp.Code, resp.Body.String())
}

func ExampleHandler_UpdateMetric() {
	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodPost, "/update/gauge/TestGauge/55", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Output: 200 {"id":"TestGauge","type":"gauge","value":55}
	fmt.Println(resp.Code, resp.Body.String())
}

func ExampleHandler_BulkUpdateJSON() {
	// Create a new Metrics object
	var testMetric = &entity.Metrics{
		ID:    "TestGauge",
		MType: entity.GaugeType,
		Value: tools.Float64Ptr(55.0),
	}

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(testMetric)

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodPost, "/update/", body)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Output: 200 {"id":"TestGauge","type":"gauge","value":55}
	fmt.Println(resp.Code, resp.Body.String())
}
