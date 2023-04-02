package handlers

import (
	"github.com/gynshu-one/go-metric-collector/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestAgent(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		pollInterval   time.Duration
		reportInterval time.Duration
	}{
		{
			name:           "basic test",
			pollInterval:   500 * time.Millisecond,
			reportInterval: 2 * time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/update/", r.URL.Path)
				assert.Equal(t, "POST", r.Method)
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			agent := NewAgent(tc.pollInterval, tc.reportInterval, server.URL)
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				agent.Start()
			}()

			time.Sleep(5 * time.Second)

			pq := &storage.Metrics{
				ID:    "PollCount",
				MType: storage.CounterType,
			}
			rv := &storage.Metrics{
				ID:    "RandomValue",
				MType: storage.GaugeType,
			}
			pollCountMetric := agent.memory.Get(pq)
			assert.NotNil(t, pollCountMetric)
			assert.Equal(t, storage.CounterType, pollCountMetric.MType)
			assert.NotNil(t, pollCountMetric.Delta)
			assert.True(t, *pollCountMetric.Delta > 0)

			randomValueMetric := agent.memory.Get(rv)
			assert.NotNil(t, randomValueMetric)
			assert.Equal(t, storage.GaugeType, randomValueMetric.MType)
			assert.NotNil(t, randomValueMetric.Value)

			runtime.Gosched()
			wg.Done()
		})
	}
}
