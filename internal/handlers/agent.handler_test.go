package handlers

import (
	"fmt"
	"github.com/gynshu-one/go-metric-collector/internal/storage"
	"github.com/gynshu-one/go-metric-collector/internal/tools"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewAgent(t *testing.T) {
	pollInterval := 5 * time.Second
	reportInterval := 10 * time.Second
	serverAddr := "http://localhost:8080"

	agent := NewAgent(pollInterval, reportInterval, serverAddr)

	assert.NotNil(t, agent)
	assert.Equal(t, pollInterval, agent.PollInterval)
	assert.Equal(t, reportInterval, agent.ReportInterval)
	assert.Equal(t, serverAddr, agent.ServerAddr)
}

func TestAgent_Poll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	}))
	defer server.Close()

	pollInterval := 1 * time.Second
	reportInterval := 2 * time.Second

	agent := NewAgent(pollInterval, reportInterval, server.URL)

	go agent.Poll()

	time.Sleep(5 * time.Second)

	allMetrics := agent.Metrics.GetAll()
	assert.NotEmpty(t, allMetrics)
}

func TestAgent_MakeReport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	}))
	defer server.Close()

	agent := NewAgent(1*time.Second, 2*time.Second, server.URL)

	testMetric := storage.Metrics{
		ID:    "TestCounter",
		MType: "counter",
		Delta: tools.Int64Ptr(1),
	}

	agent.MakeReport(testMetric)
}
