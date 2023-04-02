package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gynshu-one/go-metric-collector/internal/storage"
	"github.com/gynshu-one/go-metric-collector/internal/tools"
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"time"
)

var client = resty.New()

type Agent struct {
	pollInterval   time.Duration
	reportInterval time.Duration
	serverAddr     string
	memory         storage.MemActions
}

func NewAgent(pollInterval, reportInterval time.Duration, serverAddr string) *Agent {
	return &Agent{
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		serverAddr:     serverAddr,
		memory:         storage.NewMemStorage(),
	}
}

// Start polls runtime Metrics and reports them to the server by calling Report()
func (a *Agent) Start() {
	pollCount := 0
	go func() {
		for {
			pollCount++
			a.readRuntime()
			// Sleep for poll interval
			time.Sleep(a.pollInterval)
		}
	}()
	for {
		time.Sleep(a.reportInterval)
		go func() {
			a.memory.Set(&storage.Metrics{
				ID:    "PollCount",
				MType: storage.CounterType,
				Delta: tools.Int64Ptr(int64(pollCount)),
			})
			a.memory.Set(&storage.Metrics{
				ID:    "RandomValue",
				MType: storage.GaugeType,
				Value: tools.Float64Ptr(rand.Float64()),
			})
			a.report()
			pollCount = 0
		}()
	}

}
func (a *Agent) readRuntime() {
	// read runtime metrics
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	input := reflect.ValueOf(memStats).Elem()
	for i := 0; i < input.NumField(); i++ {
		switch input.Field(i).Kind() {
		case reflect.Uint64, reflect.Uint32:
			value := float64(input.Field(i).Uint())
			m := storage.Metrics{
				ID:    input.Type().Field(i).Name,
				MType: storage.GaugeType,
				Value: &value,
			}
			a.memory.Set(&m)
		case reflect.Float64:
			value := input.Field(i).Float()
			m := storage.Metrics{
				ID:    input.Type().Field(i).Name,
				MType: storage.GaugeType,
				Value: &value,
			}
			a.memory.Set(&m)
		}
	}
}
func (a *Agent) report() {
	// check if the metric is presented in MemStorage
	a.memory.ApplyToAll(a.makeReport)
}

// MakeReport makes a report to the server
// Notice that serverAddr must include the protocol
func (a *Agent) makeReport(m *storage.Metrics) {
	jsonData, err := json.Marshal(&m)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(jsonData).
		Post(a.serverAddr + "/update/")

	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Response: %v", resp)
}
