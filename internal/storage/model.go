package storage

import (
	"sync"
)

const (
	GaugeType             = "gauge"
	CounterType           = "counter"
	InvalidType           = "invalid type"
	TypeValueMismatch     = "type and value mismatch"
	NameTypeMismatch      = "name and type you have sent mismatch with the one in the storage"
	MetricTypeNotProvided = "metric type not provided"
	MetricNameNotProvided = "metric name not provided"
	MetricNotFound        = "metric not found"
)

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

// MemStorage is a struct that stores all metrics
// It Should be initialized using InitAgentStorage() before using
// because it has a predefined set of metrics
type MemStorage struct {
	repo sync.Map
}

type MemActions interface {
	Get(m *Metrics) *Metrics
	Set(m *Metrics) *Metrics
	ApplyToAll(f ApplyToAll, exclude ...string)
	Dump()
	Restore()
}
type ApplyToAll func(*Metrics)
