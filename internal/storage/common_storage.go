package storage

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gynshu-one/go-metric-collector/internal/tools"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
)

var Gauge = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

var Counter = []string{
	"PollCount",
	"RandomValue",
}

// MemStorage is a struct that stores all metrics
// It Should be initialized using InitStorage() before using
// because it has a predefined set of metrics
type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

// ApplyToAll applies a function to all metrics in MemStorage
type ApplyToAll func(tp string, name string, value string)

func InitStorage() *MemStorage {
	Mem := &MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
	for _, v := range Gauge {
		Mem.Gauge[v] = 0
	}
	for _, v := range Counter {
		Mem.Counter[v] = 0
	}
	return Mem
}

// AddMetric adds single metrics to MemStorage
func (M *MemStorage) AddMetric(tp string, name string, value string) error {
	switch tp {
	case "gauge":
		ui, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		M.Gauge[name] = ui
	case "counter":
		ui, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		if name == "PollCount" {
			M.Counter[name] += ui
		} else {
			M.Counter[name] = ui
		}
	}
	return nil
}

// ReadRuntime reads all values of runtime.MemStats and stores it in MemStorage
func (M *MemStorage) ReadRuntime(memStats *runtime.MemStats) {
	v := reflect.ValueOf(memStats).Elem()
	for i := 0; i < v.NumField(); i++ {
		if _, ok := M.Gauge[v.Type().Field(i).Name]; ok {
			if v.Field(i).Kind() == reflect.Uint64 {
				M.Gauge[v.Type().Field(i).Name] = float64(v.Field(i).Uint())
			} else if v.Field(i).Kind() == reflect.Float64 {
				M.Gauge[v.Type().Field(i).Name] = v.Field(i).Float()
			}
		}
	}
	M.Counter["PollCount"]++
	M.Counter["RandomValue"] = rand.Int63()
}

// ApplyToAll applies the function f to all metrics in the MemStorage
// Default exclusion is "PauseNs", "PauseEnd", "EnableGC", "DebugGC", "BySize"
// Additional exclusion can be passed as a variadic argument
// This function does not change the value of any metric
func (M *MemStorage) ApplyToAll(f ApplyToAll, exclude ...string) {
	var defaultExclusion = []string{"PauseNs", "PauseEnd", "EnableGC", "DebugGC", "BySize"}
	defaultExclusion = append(defaultExclusion, exclude...)
	for k, v := range M.Gauge {
		if !tools.Contains(defaultExclusion, k) {
			f("gauge", k, fmt.Sprint(v))
		}
	}
	for k, v := range M.Counter {
		if !tools.Contains(defaultExclusion, k) {
			f("counter", k, fmt.Sprint(v))
		}
	}
}

// PrintAll prints all metrics in MemStorage
func (M *MemStorage) PrintAll() {
	color.Green("Metric collected so far:\n" +
		"----------------------------------------\n")
	M.ApplyToAll(func(tp string, name string, value string) {
		color.Cyan("Type %s Name %s Value %s", tp, name, value)
	})
	color.Green("----------------------------------------\n")
}

// CheckIfNameExists checks if a metric name exists in MemStorage
func (M *MemStorage) CheckIfNameExists(name string) bool {
	var found bool
	M.ApplyToAll(func(tp string, n string, value string) {
		if n == name {
			found = true
		}
	})
	return found
}
