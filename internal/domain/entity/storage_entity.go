package entity

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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
	InvalidHash           = "invalid hash"
	HashNotProvided       = "env var KEY is set but hash is missing"

	DBConnError   = "db connection error"
	InvalidMetric = "invalid metric"
	EmptyMetric   = "empty metric"
)

type Metrics struct {
	ID    string   `json:"id" db:"id,primarykey"`
	MType string   `json:"type" db:"type"`
	Delta *int64   `json:"delta,omitempty" db:"delta,omitempty" `
	Value *float64 `json:"value,omitempty" db:"value,omitempty"`
	Hash  string   `json:"hash,omitempty" db:"hash,omitempty"`
}

type ApplyToAll func(*Metrics)

// CalculateAndWriteHash calculates HMAC hash of the message with the key and writes it to the Hash field
// May be this function is violating the single responsibility principle (?)
// In my opinion this is the best way to do it, otherwise we would have to calculate the hash
// in other package where we should import Metrics struct for simplicity, or in both handlers
func (M *Metrics) CalculateAndWriteHash(key string) []byte {
	//HMAC Sign  message hash with key
	if key == "" {
		M.Hash = ""
		return nil
	}
	h := hmac.New(sha256.New, []byte(key))
	message := ""
	switch M.MType {
	case GaugeType:
		value := *M.Value
		message = fmt.Sprintf("%s:%s:%f", M.ID, M.MType, value)
	case CounterType:
		delta := *M.Delta
		message = fmt.Sprintf("%s:%s:%d", M.ID, M.MType, delta)
	}
	h.Write([]byte(message))
	M.Hash = hex.EncodeToString(h.Sum(nil))
	return h.Sum(nil)
}
