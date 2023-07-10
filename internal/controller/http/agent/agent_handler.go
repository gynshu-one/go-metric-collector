// Package agent contains the implementation of the agent
// It is responsible for collecting metrics and reporting them to the server
// Basically it collects metrics from the runtime, uses reflection to get metrics by field name
// and reports them to the server
package agent

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	config "github.com/gynshu-one/go-metric-collector/internal/config/agent"
	"github.com/gynshu-one/go-metric-collector/internal/domain/entity"
	"github.com/gynshu-one/go-metric-collector/internal/domain/service"
	"github.com/gynshu-one/go-metric-collector/internal/tools"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v3/mem"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"time"
)

var client = resty.New()

type handler struct {
	mu      sync.Mutex
	memory  service.MemStorage
	workers service.WorkerPool
}

type Handler interface {
	Start()
}

func NewAgent(storage service.MemStorage) *handler {
	return &handler{
		memory:  storage,
		workers: service.NewWorkerPool(config.GetConfig().Agent.RateLimit),
	}
}

// Start polls runtime Metrics and reports them to the server by calling Report()
func (h *handler) Start() {
	pollCount := 0
	// Common metrics collection
	go func() {
		for {
			h.mu.Lock()
			pollCount++
			h.mu.Unlock()
			h.readRuntime()
			time.Sleep(config.GetConfig().Agent.PollInterval)
		}
	}()
	// Additional metrics collection
	go func() {
		for {
			h.mu.Lock()
			pollCount++
			h.mu.Unlock()
			h.readRuntime()
			// Sleep for poll interval
			time.Sleep(config.GetConfig().Agent.PollInterval)
		}
	}()
	for {
		bef, _ := cpu.Get()
		time.Sleep(config.GetConfig().Agent.ReportInterval)
		aft, _ := cpu.Get()
		total := float64(aft.Total-bef.Total) * 100
		h.memory.Set(&entity.Metrics{
			ID:    "CPUutilization1",
			MType: entity.GaugeType,
			Value: tools.Float64Ptr(float64(aft.System-bef.System) / total),
		})
		h.readAdditionalMetrics()
		go func() {
			h.mu.Lock()
			h.memory.Set(&entity.Metrics{
				ID:    "PollCount",
				MType: entity.CounterType,
				Delta: tools.Int64Ptr(int64(pollCount)),
			})
			h.memory.Set(&entity.Metrics{
				ID:    "RandomValue",
				MType: entity.GaugeType,
				Value: tools.Float64Ptr(rand.Float64()),
			})
			pollCount = 0
			h.mu.Unlock()
			h.report()
		}()
	}

}
func (h *handler) readAdditionalMetrics() {
	v, _ := mem.VirtualMemory()
	h.memory.Set(&entity.Metrics{
		ID:    "TotalMemory",
		MType: entity.GaugeType,
		Value: tools.Float64Ptr(float64(v.Total)),
	})
	h.memory.Set(&entity.Metrics{
		ID:    "FreeMemory",
		MType: entity.GaugeType,
		Value: tools.Float64Ptr(float64(v.Free)),
	})
}
func (h *handler) readRuntime() {
	// read runtime metrics
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	input := reflect.ValueOf(memStats).Elem()
	for i := 0; i < input.NumField(); i++ {
		switch input.Field(i).Kind() {
		case reflect.Uint64, reflect.Uint32:
			value := float64(input.Field(i).Uint())
			m := entity.Metrics{
				ID:    input.Type().Field(i).Name,
				MType: entity.GaugeType,
				Value: &value,
			}
			h.memory.Set(&m)
		case reflect.Float64:
			value := input.Field(i).Float()
			m := entity.Metrics{
				ID:    input.Type().Field(i).Name,
				MType: entity.GaugeType,
				Value: &value,
			}
			h.memory.Set(&m)
		}
	}
	log.Info().Msg("Runtime metrics read successfully")
}
func (h *handler) report() {
	if config.GetConfig().Key != "" {
		h.memory.ApplyToAll(func(m *entity.Metrics) {
			m.Hash = m.CalculateHash(config.GetConfig().Key)
		})
	}
	log.Debug().Msg("Trying to report metrics by bulk")
	h.workers.Push(&service.Task{
		ID: "bulkReport",
		Task: func() {
			err := h.bulkReport()
			if !errors.Is(err, entity.ErrBulkReport) {
				return
			}
		},
	})
	log.Debug().Msg("Bulk report unsuccessful, reporting metrics one by one")
	h.workers.Push(&service.Task{
		ID: "makeReport",
		Task: func() {
			h.makeReport()
		},
	})

}

// MakeReport makes a report to the server
// Notice that serverAddr must include the protocol
func (h *handler) makeReport() {
	for _, m := range h.memory.GetAll() {
		var err error
		jsonData, err := json.Marshal(m)
		if err != nil {
			log.Fatal().Err(err).Msg("Error marshalling metrics")
		}
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(jsonData).
			Post(config.GetConfig().Server.Address + "/update/")
		if err != nil {
			log.Error().Err(err).Msg("Error reporting metrics one by one")
			return
		}
		err = resp.RawBody().Close()
		if err != nil {
			log.Error().Err(err).Msg("Error closing Resty response body")
			return
		}
	}
}

func (h *handler) bulkReport() error {
	m := h.memory.GetAll()
	if len(m) == 0 {
		return nil
	}
	var err error
	jsonData, err := json.Marshal(&m)
	if err != nil {
		log.Fatal().Err(err).Msg("Error marshalling metrics")
	}
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(jsonData).
		Post(config.GetConfig().Server.Address + "/updates/")
	if err != nil {
		if resp.StatusCode() == 404 {
			log.Debug().Msgf("Path is unavailable: %v", resp)
			return entity.ErrBulkReport
		}
		log.Error().Err(err).Msg("Error reporting metrics by bulk")
		return err
	}

	return nil
}
