package storage

import "sync"

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
	Collection *sync.Map
}
