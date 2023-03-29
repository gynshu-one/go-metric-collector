package storage

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/gynshu-one/go-metric-collector/internal/tools"
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
)

type AgentInterface interface {
	FindMetricByName(name string) (Metrics, bool)
	RandomValue()
	AddPollCount()
	ReadRuntime()
	ApplyToAll(f ApplyToAll, exclude ...string)
	PrintAll()
	GetAll() []Metrics
}

// ApplyToAll applies a function to all metrics in MemStorage
type ApplyToAll func(Metrics)

func InitAgentStorage() AgentInterface {
	return &MemStorage{
		Collection: &sync.Map{},
	}
}

func (M *MemStorage) RandomValue() {
	act, load := M.Collection.LoadOrStore("RandomValue", Metrics{
		ID:    "RandomValue",
		MType: "gauge",
		Value: tools.Float64Ptr(rand.Float64()),
	})
	if load {
		*act.(Metrics).Value = rand.Float64()
		M.Collection.Store("RandomValue", act)
	}
}

// AddPollCount adds 1 to the PollCount metric if not presented creates it
func (M *MemStorage) AddPollCount() {
	act, load := M.Collection.LoadOrStore("PollCount", Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: tools.Int64Ptr(1),
	})
	if load {
		*act.(Metrics).Delta += 1
		M.Collection.Store("PollCount", act)
	}
}

// ReadRuntime reads all values of runtime.MemStats and stores it in MemStorage
func (M *MemStorage) ReadRuntime() {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	input := reflect.ValueOf(memStats).Elem()
	for i := 0; i < input.NumField(); i++ {
		switch input.Field(i).Kind() {
		case reflect.Uint64, reflect.Uint32:
			value := float64(input.Field(i).Uint())
			m := Metrics{
				ID:    input.Type().Field(i).Name,
				MType: "gauge",
				Value: &value,
			}
			M.Collection.Store(m.ID, m)
		case reflect.Float64:
			value := input.Field(i).Float()
			m := Metrics{
				ID:    input.Type().Field(i).Name,
				MType: "gauge",
				Value: &value,
			}
			M.Collection.Store(m.ID, m)
		}
	}
}

// ApplyToAll applies the function f to all metrics in the MemStorage
// Default exclusion is "PauseNs", "PauseEnd", "EnableGC", "DebugGC", "BySize"
// Additional exclusion can be passed as a variadic argument
// This function does not change the value of any metric
func (M *MemStorage) ApplyToAll(f ApplyToAll, exclude ...string) {
	var defaultExclusion = []string{"PauseNs", "PauseEnd", "EnableGC", "DebugGC", "BySize"}
	defaultExclusion = append(defaultExclusion, exclude...)
	M.Collection.Range(func(key, value interface{}) bool {
		if !tools.Contains(defaultExclusion, key.(string)) {
			f(value.(Metrics))
		}
		return true
	})
}

// PrintAll prints all metrics in MemStorage
func (M *MemStorage) PrintAll() {
	color.Green("Metric collected so far:\n" +
		"----------------------------------------\n")
	M.ApplyToAll(func(m Metrics) {
		jsn, err := json.MarshalIndent(m, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(jsn))
	})
	color.Green("----------------------------------------\n")
}

func (M *MemStorage) GetAll() []Metrics {
	var metrics []Metrics
	M.ApplyToAll(func(m Metrics) {
		metrics = append(metrics, m)
	})
	return metrics
}
