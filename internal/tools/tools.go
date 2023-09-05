package tools

import (
	"github.com/gynshu-one/go-metric-collector/internal/domain/entity"
	"github.com/gynshu-one/go-metric-collector/proto"
)

func Contains(sl []string, s string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}
	return false
}
func Int64Ptr(i int64) *int64 {
	return &i
}
func Float64Ptr(f float64) *float64 {
	return &f
}

func MarshalMetric(m *entity.Metrics) *proto.Metric {
	metric := &proto.Metric{
		ID:    m.ID,
		MType: m.MType,
	}
	if m.Value != nil {
		metric.Value = *m.Value
	}
	if m.Delta != nil {
		metric.Delta = *m.Delta
	}
	return metric
}

func UnmarshalMetric(m *proto.Metric) *entity.Metrics {
	return &entity.Metrics{
		ID:    m.ID,
		MType: m.MType,
		Value: &m.Value,
		Delta: &m.Delta,
	}
}
