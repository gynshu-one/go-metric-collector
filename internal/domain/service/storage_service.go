package service

import (
	"github.com/gynshu-one/go-metric-collector/internal/domain/entity"
	"github.com/gynshu-one/go-metric-collector/internal/tools"
	"log"
	"sync"
)

type MemStorage interface {
	Get(m *entity.Metrics) *entity.Metrics
	Set(m *entity.Metrics) *entity.Metrics
	ApplyToAll(f entity.ApplyToAll, exclude ...string)
}

type memService struct {
	repo *sync.Map
}

func NewMemService(repo *sync.Map) *memService {
	return &memService{repo: repo}
}

func (M memService) Get(m *entity.Metrics) *entity.Metrics {
	found, ok := M.repo.Load(m.ID)
	if !ok {
		return nil
	}
	return found.(*entity.Metrics)
}

func (M memService) Set(m *entity.Metrics) *entity.Metrics {
	found, ok := M.repo.LoadOrStore(m.ID, m)
	if ok {
		if found.(*entity.Metrics).MType != m.MType {
			log.Printf("name and type you have sent mismatch with the one in the storage: %s", m.ID)
			return nil
		}
		switch m.MType {
		case entity.GaugeType:
			*found.(*entity.Metrics).Value = *m.Value
			found.(*entity.Metrics).FloatPrecision = m.FloatPrecision
		case entity.CounterType:
			*found.(*entity.Metrics).Delta += *m.Delta
		}
		M.repo.Store(m.ID, found.(*entity.Metrics))
	}
	return found.(*entity.Metrics)
}

// ApplyToAll applies the function f to all metrics in the MemStorage
// Default exclusion is "PauseNs", "PauseEnd", "EnableGC", "DebugGC", "BySize"
// Additional exclusion can be passed as a variadic argument
// This function does not change the value of any metric
func (M memService) ApplyToAll(f entity.ApplyToAll, exclude ...string) {
	var defaultExclusion = []string{"PauseNs", "PauseEnd", "EnableGC", "DebugGC", "BySize"}
	defaultExclusion = append(defaultExclusion, exclude...)
	M.repo.Range(func(key, value interface{}) bool {
		if !tools.Contains(defaultExclusion, key.(string)) {
			f(value.(*entity.Metrics))
		}
		return true
	})
}
