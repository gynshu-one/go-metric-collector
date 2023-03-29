package storage

import (
	"encoding/json"
	"fmt"
	"github.com/gynshu-one/go-metric-collector/internal/configs"
	"log"
	"os"
	"sync"
	"time"
)

var Memory ServerInterface

type ServerInterface interface {
	StoreEverythingToFile() error
	LoadEverythingFromFile() error
	FindMetricByName(name string) (Metrics, bool)
	CheckMetricType(tp string) bool
	UpdateMetric(m Metrics) Metrics
	ValidateValue(m Metrics) bool
	ValidateTypeAndValue(m Metrics) bool
	CheckIfNameExists(name string) bool
	GenerateHTMLTable() []string
	GetAll() []Metrics
	ResetAll()
}

func InitServerStorage() ServerInterface {
	mem := &MemStorage{
		Collection: &sync.Map{},
	}
	if configs.CFG.Restore {
		err := mem.LoadEverythingFromFile()
		if err != nil {
			log.Fatal(err)
		}
	}
	if configs.CFG.StoreInterval != 0 {
		ticker := time.NewTicker(configs.CFG.StoreInterval)
		go func() {
			for {
				t := <-ticker.C
				err := mem.StoreEverythingToFile()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("Saved to file at", t)
			}
		}()
	}
	return mem

}

// ResetAll resets all metrics. Used for testing
func (M *MemStorage) ResetAll() {
	M.Collection = &sync.Map{}
}
func (M *MemStorage) StoreEverythingToFile() error {
	allMetrics := M.GetAll()
	// save to jsonData file
	jsonData, err := json.Marshal(allMetrics)
	if err != nil {
		return err
	}
	// save to file
	err = os.WriteFile(configs.CFG.StoreFile, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
	//path := configs.CFG.StoreFile
}
func (M *MemStorage) LoadEverythingFromFile() error {
	file, err := os.OpenFile(configs.CFG.StoreFile, os.O_RDONLY, 0666)
	if err != nil {
		log.Printf("Nothing to resore from storage file: %v", err)
		return nil
	}
	defer file.Close()
	var metrics []Metrics
	err = json.NewDecoder(file).Decode(&metrics)
	if err != nil {
		log.Printf("Error decoding json may be file is empty: %v", err)
		return nil
	}
	for _, m := range metrics {
		M.Collection.Store(m.ID, m)
	}
	metrics = nil
	return nil
}

func (M *MemStorage) ValidateTypeAndValue(m Metrics) bool {
	if (m.MType == "gauge" && m.Value != nil) || (m.MType == "counter" && m.Delta != nil) {
		return true
	}
	return false
}
func (M *MemStorage) ValidateValue(m Metrics) bool {
	if m.Value == nil && m.Delta == nil {
		return false
	}
	return true
}

// FindMetricByName finds a metric by name and returns its value
// If the metric is not found, it returns false
func (M *MemStorage) FindMetricByName(name string) (Metrics, bool) {
	m, ok := M.Collection.Load(name)
	if !ok {
		return Metrics{}, false
	}
	return m.(Metrics), true
}

// CheckMetricType checks if the metric type is presented in MemStorage
func (M *MemStorage) CheckMetricType(tp string) bool {
	switch tp {
	case "gauge", "counter":
		return true
	default:
		return false
	}
}

// UpdateMetric adds single metrics to MemStorage
func (M *MemStorage) UpdateMetric(m Metrics) Metrics {
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
				M.Collection.Store(m.ID, act)
			}
		}
	}
	value, ok := M.Collection.Load(m.ID)
	if !ok {
		return Metrics{}
	}
	if configs.CFG.StoreInterval == 0 {
		go func() {
			err := M.StoreEverythingToFile()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Saved to file asynchronously")
		}()
	}
	return value.(Metrics)
}

// CheckIfNameExists checks if a metric name exists in MemStorage
func (M *MemStorage) CheckIfNameExists(name string) bool {
	_, ok := M.Collection.Load(name)
	return ok
}

func (M *MemStorage) GenerateHTMLTable() []string {
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
