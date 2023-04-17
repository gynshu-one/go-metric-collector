package entity

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

type Metrics struct {
	ID    string   `json:"id" db:"id,primarykey"`
	MType string   `json:"type" db:"type"`
	Delta *int64   `json:"delta,omitempty" db:"delta,omitempty" `
	Value *float64 `json:"value,omitempty" db:"value,omitempty"`
	Hash  string   `json:"hash,omitempty" db:"hash,omitempty"`
}

func (M *Metrics) String() string {
	delta := int64(0)
	value := float64(0)
	if M.Delta != nil {
		delta = *M.Delta
	}
	if M.Value != nil {
		value = *M.Value
	}
	return fmt.Sprintf("\nID: %s, \nType: %s, \nDelta: %d, \nValue: %f, \nHash: %s", M.ID, M.MType, delta, value, M.Hash)
}

type ApplyToAll func(*Metrics)

// CalculateHash calculates HMAC hash of the message with the key and writes it to the Hash field
// May be this function is violating the single responsibility principle (?)
// In my opinion this is the best way to do it, otherwise we would have to calculate the hash
// in other package where we should import Metrics struct for simplicity, or in both handlers
func (M *Metrics) CalculateHash(key string) string {
	if key == "" {
		return ""
	}
	h := hmac.New(sha256.New, []byte(key))
	message := ""
	switch M.MType {
	case GaugeType:
		message = fmt.Sprintf("%s:%s:%f", M.ID, M.MType, *M.Value)
	case CounterType:
		message = fmt.Sprintf("%s:%s:%d", M.ID, M.MType, *M.Delta)
	default:
		return ""
	}
	h.Write([]byte(message))
	M.Hash = hex.EncodeToString(h.Sum(nil))
	return M.Hash
}
