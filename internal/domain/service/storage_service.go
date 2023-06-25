package service

import (
	"github.com/gynshu-one/go-metric-collector/internal/domain/entity"
	"github.com/gynshu-one/go-metric-collector/internal/tools"
	"github.com/rs/zerolog/log"
	_ "net/http/pprof"
	"sync"
)

// MemStorage is the interface for the storage service
// It provides methods to get and set metrics as well as apply a function to all metrics
type MemStorage interface {
	Get(id string) *entity.Metrics
	Set(m *entity.Metrics) *entity.Metrics
	ApplyToAll(f entity.ApplyToAll, exclude ...string)
	GetAll() []*entity.Metrics
}
type memService struct {
	repo map[string]*entity.Metrics
	mu   sync.Mutex
}

func NewMemService() *memService {
	return &memService{repo: make(map[string]*entity.Metrics)}
}

// Get retrieves a metric from the storage
func (M *memService) Get(id string) *entity.Metrics {
	M.mu.Lock()
	defer M.mu.Unlock()
	return M.repo[id]
}

// Set stores a metric in the storage
// If the metric already exists, it will be updated
func (M *memService) Set(m *entity.Metrics) *entity.Metrics {
	if m == nil {
		return nil
	}
	M.mu.Lock()
	defer M.mu.Unlock()
	found, ok := M.repo[m.ID]
	if ok {
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

// ApplyToAll applies a function to all metrics in the storage
// It is used to update the metrics when a new interval starts
// You can exclude some metrics from the update by passing their name as a parameter
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

// GetAll returns all metrics in the storage
func (M *memService) GetAll() []*entity.Metrics {
	M.mu.Lock()
	defer M.mu.Unlock()
	metrics := make([]*entity.Metrics, 0, len(M.repo))
	for _, v := range M.repo {
		metrics = append(metrics, v)
	}
	return metrics
}
