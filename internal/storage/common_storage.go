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

// FindMetricByName finds a metric by name and returns its value
// If the metric is not found, it returns an error
func (M *MemStorage) FindMetricByName(tp string, name string) (string, error) {
	// this is general approach, in case we want to add more metric types
	f := reflect.ValueOf(M).Elem().FieldByName(tp)
	if !f.IsValid() {
		return "", fmt.Errorf("metric type is not presented%s", tp)
	}
	metric := f.MapIndex(reflect.ValueOf(name))
	if !metric.IsValid() {
		return "", fmt.Errorf("metric name is not presented%s", name)
	} else {
		return fmt.Sprintf("%v", metric), nil
	}
}

// CheckMetricType checks if the metric type is presented in MemStorage
func (M *MemStorage) CheckMetricType(tp string) bool {
	f := reflect.ValueOf(M).Elem().FieldByName(tp)
	return f.IsValid()
}

// AddMetric adds single metrics to MemStorage
func (M *MemStorage) AddMetric(tp string, name string, value string) error {
	switch tp {
	case "Gauge":
		ui, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		ui = tools.ToFixed(ui, 3)
		M.Gauge[name] = ui
	case "Counter":
		ui, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		if _, ok := M.Counter[name]; !ok {
			M.Counter[name] = ui
		} else {
			M.Counter[name] += ui
		}
	}
	return nil
}

// ReadRuntime reads all values of runtime.MemStats and stores it in MemStorage
func (M *MemStorage) ReadRuntime(memStats *runtime.MemStats) {
	input := reflect.ValueOf(memStats).Elem()
	for i := 0; i < input.NumField(); i++ {
		switch input.Field(i).Kind() {
		case reflect.Uint64:
			value := float64(input.Field(i).Uint())
			M.Gauge[input.Type().Field(i).Name] = value
		case reflect.Float64:
			value := input.Field(i).Float()
			M.Gauge[input.Type().Field(i).Name] = value
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
	v := reflect.ValueOf(M).Elem()
	for i := 0; i < v.NumField(); i++ {
		metricType := v.Type().Field(i).Name
		metric := v.Field(i)
		for _, key := range metric.MapKeys() {
			if !tools.Contains(defaultExclusion, key.String()) {
				f(metricType, key.String(), fmt.Sprintf("%v", metric.MapIndex(key)))
			}
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

func (M *MemStorage) GenerateHTMLTable() []string {
	var table []string
	M.ApplyToAll(func(tp string, name string, value string) {
		table = append(table, fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td></tr>", tp, name, value))
	})
	return table
}
