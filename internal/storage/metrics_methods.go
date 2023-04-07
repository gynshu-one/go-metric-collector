package storage

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gynshu-one/go-metric-collector/internal/configs"
)

// CalculateAndWriteHash calculates HMAC hash of the message with the key and writes it to the Hash field
// May be this function is violating the single responsibility principle (?)
// In my opinion this is the best way to do it, otherwise we would have to calculate the hash
// in other package where we should import Metrics struct for simplicity, or in both handlers
func (M *Metrics) CalculateAndWriteHash() []byte {
	key := configs.CFG.Key
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
