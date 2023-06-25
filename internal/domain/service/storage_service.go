package service

import (
	"github.com/gynshu-one/go-metric-collector/internal/domain/entity"
	"github.com/gynshu-one/go-metric-collector/internal/tools"
	"github.com/rs/zerolog/log"
	_ "net/http/pprof"
	"sync"
)

type MemStorage interface {
	Get(id string) *entity.Metrics
	Set(m *entity.Metrics) *entity.Metrics
	ApplyToAll(f entity.ApplyToAll, exclude ...string)
	GetAll() []*entity.Metrics
}

type memService struct {
	repo map[string]*entity.Metrics
	mu   *sync.Mutex
}

func NewMemService() *memService {
	return &memService{repo: make(map[string]*entity.Metrics), mu: new(sync.Mutex)}
}

func (M memService) Get(id string) *entity.Metrics {
	return M.repo[id]
}

func (M memService) Set(m *entity.Metrics) *entity.Metrics {
	M.mu.Lock()
	defer M.mu.Unlock()
	found := M.repo[m.ID]
	if found != nil {
		if m.MType == entity.CounterType {
			m.Delta = tools.Int64Ptr(*found.Delta + *m.Delta)
		}
		M.repo[m.ID] = m
	} else {
		M.repo[m.ID] = m
	}
	found = M.repo[m.ID]
	if found == nil {
		log.Error().Msgf("Failed to store metric: %s", m.String())
		return nil
	}
	return found
}

// ApplyToAll applies the function f to all metrics in the MemStorage
// Default exclusion is "PauseNs", "PauseEnd", "EnableGC", "DebugGC", "BySize"
// Additional exclusion can be passed as a variadic argument
// This function does not change the value of any metric
func (M memService) ApplyToAll(f entity.ApplyToAll, exclude ...string) {
	var defaultExclusion = []string{"PauseNs", "PauseEnd", "EnableGC", "DebugGC", "BySize"}
	defaultExclusion = append(defaultExclusion, exclude...)
	for k, v := range M.repo {
		M.mu.Lock()
		if !tools.Contains(defaultExclusion, k) {
			f(v)
		}
		M.mu.Unlock()
	}
}

func (M memService) GetAll() []*entity.Metrics {
	M.mu.Lock()
	defer M.mu.Unlock()
	metrics := make([]*entity.Metrics, len(M.repo))
	for _, v := range M.repo {
		metrics = append(metrics, v)
	}
	return metrics
}
