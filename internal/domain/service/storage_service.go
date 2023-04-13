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
	GetAll() []*entity.Metrics
}

type memService struct {
	mu   sync.RWMutex
	repo map[string]*entity.Metrics
}

func NewMemService(repo map[string]*entity.Metrics) *memService {
	return &memService{repo: repo}
}

func (M *memService) Get(m *entity.Metrics) *entity.Metrics {
	M.mu.Lock()
	defer M.mu.Unlock()
	return M.repo[m.ID]
}

func (M *memService) Set(m *entity.Metrics) *entity.Metrics {
	M.mu.Lock()
	defer M.mu.Unlock()
	_, ok := M.repo[m.ID]
	if ok {
		if M.repo[m.ID].MType != m.MType {
			log.Printf("name and type you have sent mismatch with the one in the storage: %s", m.ID)
			return nil
		}
		switch m.MType {
		case entity.GaugeType:
			M.repo[m.ID].Value = m.Value
		case entity.CounterType:
			*M.repo[m.ID].Delta = *M.repo[m.ID].Delta + *m.Delta
		}
	} else {
		M.repo[m.ID] = m
	}
	return M.repo[m.ID]
}

// ApplyToAll applies the function f to all metrics in the MemStorage
// Default exclusion is "PauseNs", "PauseEnd", "EnableGC", "DebugGC", "BySize"
// Additional exclusion can be passed as a variadic argument
// This function does not change the value of any metric
func (M *memService) ApplyToAll(f entity.ApplyToAll, exclude ...string) {
	var defaultExclusion = []string{"PauseNs", "PauseEnd", "EnableGC", "DebugGC", "BySize"}
	defaultExclusion = append(defaultExclusion, exclude...)
	M.mu.Lock()
	defer M.mu.Unlock()
	for k, v := range M.repo {
		if !tools.Contains(defaultExclusion, k) {
			f(v)
		}
	}
}

func (M *memService) GetAll() []*entity.Metrics {
	var metrics []*entity.Metrics
	M.mu.Lock()
	defer M.mu.Unlock()
	for _, v := range M.repo {
		metrics = append(metrics, v)
	}
	return metrics
}
