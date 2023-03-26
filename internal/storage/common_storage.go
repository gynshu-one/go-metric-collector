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

//var Gauge = []string{
//	"Alloc",
//	"BuckHashSys",
//	"Frees",
//	"GCCPUFraction",
//	"GCSys",
//	"HeapAlloc",
//	"HeapIdle",
//	"HeapInuse",
//	"HeapObjects",
//	"HeapReleased",
//	"HeapSys",
//	"LastGC",
//	"Lookups",
//	"MCacheInuse",
//	"MCacheSys",
//	"MSpanInuse",
//	"MSpanSys",
//	"Mallocs",
//	"NextGC",
//	"NumForcedGC",
//	"NumGC",
//	"OtherSys",
//	"PauseTotalNs",
//	"StackInuse",
//	"StackSys",
//	"Sys",
//	"TotalAlloc",
//}
//
//var Counter = []string{
//	"PollCount",
//	"RandomValue",
//}

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

// MemStorage is a struct that stores all metrics
// It Should be initialized using InitStorage() before using
// because it has a predefined set of metrics
type MemStorage struct {
	Collection *sync.Map
}
type MemInterface interface {
	FindMetricByName(name string) (Metrics, bool)
	CheckMetricType(tp string) bool
	UpdateMetric(m Metrics) Metrics
	RandomValue()
	AddPollCount()
	ValidateValue(m Metrics) bool
	ValidateTypeAndValue(m Metrics) bool
	ReadRuntime()
	ApplyToAll(f ApplyToAll, exclude ...string)
	PrintAll()
	GetAll() []Metrics
	CheckIfNameExists(name string) bool
	GenerateHTMLTable() []string
}

// ApplyToAll applies a function to all metrics in MemStorage
type ApplyToAll func(Metrics)

func InitStorage() MemInterface {
	return MemStorage{
		Collection: &sync.Map{},
	}
}

func (M MemStorage) ValidateTypeAndValue(m Metrics) bool {
	if (m.MType == "gauge" && m.Value != nil) || (m.MType == "counter" && m.Delta != nil) {
		return true
	}
	return false
}
func (M MemStorage) ValidateValue(m Metrics) bool {
	if m.Value == nil && m.Delta == nil {
		return false
	}
	return true
}
func (M MemStorage) RandomValue() {
	act, load := M.Collection.LoadOrStore("RandomValue", Metrics{
		ID:    "RandomValue",
		MType: "gauge",
		Value: tools.Float64Ptr(rand.Float64()),
	})
	if load {
		*act.(Metrics).Value = rand.Float64()
		M.Collection.Swap("RandomValue", act)
	}
}

// AddPollCount adds 1 to the PollCount metric if not presented creates it
func (M MemStorage) AddPollCount() {
	act, load := M.Collection.LoadOrStore("PollCount", Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: tools.Int64Ptr(1),
	})
	if load {
		*act.(Metrics).Delta += 1
		M.Collection.Swap("PollCount", act)
	}
}

// FindMetricByName finds a metric by name and returns its value
// If the metric is not found, it returns false
func (M MemStorage) FindMetricByName(name string) (Metrics, bool) {
	m, ok := M.Collection.Load(name)
	if !ok {
		return Metrics{}, false
	}
	return m.(Metrics), true
}

// CheckMetricType checks if the metric type is presented in MemStorage
func (M MemStorage) CheckMetricType(tp string) bool {
	switch tp {
	case "gauge", "counter":
		return true
	default:
		return false
	}
}

// UpdateMetric adds single metrics to MemStorage
func (M MemStorage) UpdateMetric(m Metrics) Metrics {
	switch m.MType {
	case "gauge":
		if m.Value != nil {
			M.Collection.Store(m.ID, m)
		}
	case "counter":
		if m.Delta != nil {
			act, load := M.Collection.LoadOrStore(m.ID, m)
			if load {
				*act.(Metrics).Delta += *m.Delta
				M.Collection.Swap(m.ID, act)
			}
		}
	}
	value, ok := M.Collection.Load(m.ID)
	if !ok {
		return Metrics{}
	}
	return value.(Metrics)
}

// ReadRuntime reads all values of runtime.MemStats and stores it in MemStorage
func (M MemStorage) ReadRuntime() {
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
func (M MemStorage) ApplyToAll(f ApplyToAll, exclude ...string) {
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
func (M MemStorage) PrintAll() {
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

// CheckIfNameExists checks if a metric name exists in MemStorage
func (M MemStorage) CheckIfNameExists(name string) bool {
	_, ok := M.Collection.Load(name)
	return ok
}

func (M MemStorage) GenerateHTMLTable() []string {
	var table []string
	M.ApplyToAll(func(m Metrics) {
		val := ""
		if m.Value != nil {
			val = fmt.Sprintf("%f", *m.Value)
		}
		if m.Delta != nil {
			val = fmt.Sprintf("%d", *m.Delta)
		}
		table = append(table, fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td></tr>",
			m.MType, m.ID, val))
	})
	return table
}

func (M MemStorage) GetAll() []Metrics {
	var metrics []Metrics
	M.ApplyToAll(func(m Metrics) {
		metrics = append(metrics, m)
	})
	return metrics
}
