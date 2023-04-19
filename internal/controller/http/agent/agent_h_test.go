package agent

import (
	config "github.com/gynshu-one/go-metric-collector/internal/config/agent"
	"github.com/gynshu-one/go-metric-collector/internal/domain/entity"
	"github.com/gynshu-one/go-metric-collector/internal/domain/service"
	"github.com/gynshu-one/go-metric-collector/internal/tools"
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
			reportInterval: 1 * time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/updates/", r.URL.Path)
				assert.Equal(t, "POST", r.Method)
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()
			config.GetConfig().Server.Address = server.URL
			newAgent := NewAgent(service.NewMemService(&sync.Map{}))
			runtime.Gosched()
			go func() {
				newAgent.Start()
			}()

			time.Sleep(12 * time.Second)

			pq := &entity.Metrics{
				ID:    "PollCount",
				MType: entity.CounterType,
				Delta: tools.Int64Ptr(1),
			}
			rv := &entity.Metrics{
				ID:    "RandomValue",
				MType: entity.GaugeType,
				Value: tools.Float64Ptr(1),
			}
			pollCountMetric := newAgent.memory.Get(pq)
			assert.NotNil(t, pollCountMetric)
			assert.Equal(t, entity.CounterType, pollCountMetric.MType)
			assert.NotNil(t, pollCountMetric.Delta)
			assert.True(t, *pollCountMetric.Delta > 0)

			randomValueMetric := newAgent.memory.Get(rv)
			assert.NotNil(t, randomValueMetric)
			assert.Equal(t, entity.GaugeType, randomValueMetric.MType)
			assert.NotNil(t, randomValueMetric.Value)
		})
	}
}
