package agent

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	config "github.com/gynshu-one/go-metric-collector/internal/config/agent"
	"github.com/gynshu-one/go-metric-collector/internal/domain/entity"
	"github.com/gynshu-one/go-metric-collector/internal/domain/service"
	"github.com/gynshu-one/go-metric-collector/internal/tools"
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"time"
)

var client = resty.New()

type handler struct {
	memory service.MemStorage
}

type Handler interface {
	Start()
}

func NewAgent(storage service.MemStorage) *handler {
	return &handler{
		memory: storage,
	}
}

// Start polls runtime Metrics and reports them to the server by calling Report()
func (h *handler) Start() {
	pollCount := 0
	go func() {
		for {
			pollCount++
			h.readRuntime()
			// Sleep for poll interval
			time.Sleep(config.GetConfig().Agent.PollInterval)
		}
	}()
	for {
		time.Sleep(config.GetConfig().Agent.ReportInterval)
		go func() {
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
			h.report()
			pollCount = 0
		}()
	}

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
}
func (h *handler) report() {
	if config.GetConfig().Key != "" {
		h.memory.ApplyToAll(func(m *entity.Metrics) {
			m.Hash = m.CalculateHash(config.GetConfig().Key)
		})
	}
	err := h.bulkReport()
	if err == nil {
		return
	}
	h.memory.ApplyToAll(makeReport)
}

// MakeReport makes a report to the server
// Notice that serverAddr must include the protocol
func makeReport(m *entity.Metrics) {
	var err error
	jsonData, err := json.Marshal(&m)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(jsonData).
		Post(config.GetConfig().Server.Address + "/update/")

	if err != nil {
		return
	}
	defer resp.RawBody().Close()
}

func (h *handler) bulkReport() error {
	m := h.memory.GetAll()
	if len(m) == 0 {
		return nil
	}
	var err error
	jsonData, err := json.Marshal(&m)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(jsonData).
		Post(config.GetConfig().Server.Address + "/updates/")
	if err != nil {
		fmt.Printf("Error: %v", err)
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("response: %v", resp)
	}
	return nil
}
