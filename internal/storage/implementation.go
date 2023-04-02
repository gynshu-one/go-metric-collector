package storage

import (
	"encoding/json"
	"github.com/gynshu-one/go-metric-collector/internal/configs"
	"github.com/gynshu-one/go-metric-collector/internal/tools"
	"log"
	"os"
	"sync"
)

func NewMemStorage() MemActions {
	return &MemStorage{
		repo: sync.Map{},
	}
}
func (M *MemStorage) Get(m *Metrics) *Metrics {
	found, ok := M.repo.Load(m.ID)
	if !ok {
		return nil
	}
	return found.(*Metrics)
}

func (M *MemStorage) Set(m *Metrics) *Metrics {
	found, ok := M.repo.LoadOrStore(m.ID, m)
	if ok {
		if found.(*Metrics).MType != m.MType {
			log.Printf("name and type you have sent mismatch with the one in the storage: %s", m.ID)
			return nil
		}
		switch m.MType {
		case GaugeType:
			*found.(*Metrics).Value = *m.Value
		case CounterType:
			*found.(*Metrics).Delta += *m.Delta
		}
		M.repo.Store(m.ID, found.(*Metrics))
	}
	return found.(*Metrics)
}

// ApplyToAll applies the function f to all metrics in the MemStorage
// Default exclusion is "PauseNs", "PauseEnd", "EnableGC", "DebugGC", "BySize"
// Additional exclusion can be passed as a variadic argument
// This function does not change the value of any metric
func (M *MemStorage) ApplyToAll(f ApplyToAll, exclude ...string) {
	var defaultExclusion = []string{"PauseNs", "PauseEnd", "EnableGC", "DebugGC", "BySize"}
	defaultExclusion = append(defaultExclusion, exclude...)
	M.repo.Range(func(key, value interface{}) bool {
		if !tools.Contains(defaultExclusion, key.(string)) {
			f(value.(*Metrics))
		}
		return true
	})
}

func (M *MemStorage) Dump() {
	allMetrics := make([]*Metrics, 0)
	M.repo.Range(func(key, value interface{}) bool {
		allMetrics = append(allMetrics, value.(*Metrics))
		return true
	})
	// save to jsonData file
	jsonData, err := json.Marshal(allMetrics)
	if err != nil {
		log.Fatal(err)
	}
	// save to file
	err = os.WriteFile(configs.CFG.StoreFile, jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//path := configs.CFG.StoreFile
}
func (M *MemStorage) Restore() {
	file, err := os.OpenFile(configs.CFG.StoreFile, os.O_RDONLY, 0666)
	if err != nil {
		log.Printf("Nothing to resore from storage file: %v", err)
		return
	}
	defer file.Close()
	var metrics []*Metrics
	err = json.NewDecoder(file).Decode(&metrics)
	if err != nil {
		log.Printf("Error decoding json may be file is empty: %v", err)
		return
	}
	for _, m := range metrics {
		M.Set(m)
	}
	metrics = nil
}
