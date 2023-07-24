// Package agent contains the implementation of the agent
// It is responsible for collecting metrics and reporting them to the server
// Basically it collects metrics from the runtime, uses reflection to get metrics by field name
// and reports them to the server
package agent

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"github.com/go-resty/resty/v2"
	config "github.com/gynshu-one/go-metric-collector/internal/config/agent"
	"github.com/gynshu-one/go-metric-collector/internal/domain/entity"
	"github.com/gynshu-one/go-metric-collector/internal/domain/service"
	"github.com/gynshu-one/go-metric-collector/internal/tools"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v3/mem"
	"io"
	mathrand "math/rand"
	"os"
	"reflect"
	"runtime"
	"sync"
	"time"
)

var client = resty.New()

type handler struct {
	mu      sync.Mutex
	memory  service.MemStorage
	workers service.WorkerPool
}

type Handler interface {
	Start()
}

func NewAgent(storage service.MemStorage) *handler {
	return &handler{
		memory:  storage,
		workers: service.NewWorkerPool(config.GetConfig().Agent.RateLimit),
	}
}

// Start polls runtime Metrics and reports them to the server by calling Report()
func (h *handler) Start() {
	pollCount := 0
	// Common metrics collection
	go func() {
		for {
			h.mu.Lock()
			pollCount++
			h.mu.Unlock()
			h.readRuntime()
			time.Sleep(config.GetConfig().Agent.PollInterval)
		}
	}()
	// Additional metrics collection
	go func() {
		for {
			h.mu.Lock()
			pollCount++
			h.mu.Unlock()
			h.readRuntime()
			// Sleep for poll interval
			time.Sleep(config.GetConfig().Agent.PollInterval)
		}
	}()
	for {
		bef, err := cpu.Get()
		if err != nil {
			log.Error().Err(err).Msg("Error getting CPU metrics")
			bef = &cpu.Stats{}
		}
		time.Sleep(config.GetConfig().Agent.ReportInterval)
		aft, err := cpu.Get()
		if err != nil {
			log.Error().Err(err).Msg("Error getting CPU metrics")
			aft = bef
		}
		if aft.Total > 0 {
			total := float64(aft.Total-bef.Total) * 100
			h.memory.Set(&entity.Metrics{
				ID:    "CPUutilization1",
				MType: entity.GaugeType,
				Value: tools.Float64Ptr(float64(aft.System-bef.System) / total),
			})
		}
		h.readAdditionalMetrics()
		go func() {
			h.mu.Lock()
			h.memory.Set(&entity.Metrics{
				ID:    "PollCount",
				MType: entity.CounterType,
				Delta: tools.Int64Ptr(int64(pollCount)),
			})
			h.memory.Set(&entity.Metrics{
				ID:    "RandomValue",
				MType: entity.GaugeType,
				Value: tools.Float64Ptr(mathrand.Float64()),
			})
			pollCount = 0
			h.mu.Unlock()
			h.report()
		}()
	}

}
func (h *handler) readAdditionalMetrics() {
	v, _ := mem.VirtualMemory()
	h.memory.Set(&entity.Metrics{
		ID:    "TotalMemory",
		MType: entity.GaugeType,
		Value: tools.Float64Ptr(float64(v.Total)),
	})
	h.memory.Set(&entity.Metrics{
		ID:    "FreeMemory",
		MType: entity.GaugeType,
		Value: tools.Float64Ptr(float64(v.Free)),
	})
}
func (h *handler) readRuntime() {
	// read runtime metrics
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	input := reflect.ValueOf(memStats).Elem()
	for i := 0; i < input.NumField(); i++ {
		switch input.Field(i).Kind() {
		case reflect.Uint64, reflect.Uint32:
			value := float64(input.Field(i).Uint())
			m := entity.Metrics{
				ID:    input.Type().Field(i).Name,
				MType: entity.GaugeType,
				Value: &value,
			}
			h.memory.Set(&m)
		case reflect.Float64:
			value := input.Field(i).Float()
			m := entity.Metrics{
				ID:    input.Type().Field(i).Name,
				MType: entity.GaugeType,
				Value: &value,
			}
			h.memory.Set(&m)
		}
	}
	log.Info().Msg("Runtime metrics read successfully")
}
func (h *handler) report() {
	if config.GetConfig().Key != "" {
		h.memory.ApplyToAll(func(m *entity.Metrics) {
			m.Hash = m.CalculateHash(config.GetConfig().Key)
		})
	}
	log.Debug().Msg("Trying to report metrics by bulk")
	h.workers.Push(&service.Task{
		ID: "bulkReport",
		Task: func() {
			err := h.bulkReport()
			if !errors.Is(err, entity.ErrBulkReport) {
				return
			}
		},
	})
	log.Debug().Msg("Bulk report unsuccessful, reporting metrics one by one")
	h.workers.Push(&service.Task{
		ID: "makeReport",
		Task: func() {
			h.makeReport()
		},
	})

}

// MakeReport makes a report to the server
// Notice that serverAddr must include the protocol
func (h *handler) makeReport() {
	for _, m := range h.memory.GetAll() {
		var err error
		jsonData, err := json.Marshal(m)
		if err != nil {
			log.Fatal().Err(err).Msg("Error marshalling metrics")
		}
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(encryptWithPublicKey(jsonData)).
			Post(config.GetConfig().Server.Address + "/update/")
		if err != nil {
			log.Error().Err(err).Msg("Error reporting metrics one by one")
			return
		}
		err = resp.RawBody().Close()
		if err != nil {
			log.Error().Err(err).Msg("Error closing Resty response body")
			return
		}
	}
}

func (h *handler) bulkReport() error {
	m := h.memory.GetAll()
	if len(m) == 0 {
		return nil
	}
	var err error
	jsonData, err := json.Marshal(&m)
	if err != nil {
		log.Fatal().Err(err).Msg("Error marshalling metrics")
	}
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(encryptWithPublicKey(jsonData)).
		Post(config.GetConfig().Server.Address + "/updates/")
	if err != nil {
		if resp.StatusCode() == 404 {
			log.Debug().Msgf("Path is unavailable: %v", resp)
			return entity.ErrBulkReport
		}
		log.Error().Err(err).Msg("Error reporting metrics by bulk")
		return err
	}

	return nil
}

var publicKey *rsa.PublicKey

func encryptWithPublicKey(body []byte) []byte {
	if config.GetConfig().CryptoKey == "" {
		return body
	}
	once.Do(func() {
		loadPublicKey(config.GetConfig().CryptoKey)
	})
	if publicKey == nil {
		return body
	}
	// Generate a new AES key
	aesKey := make([]byte, 32) // 256 bits
	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		log.Fatal().Err(err).Msg("Error generating AES key")
		return body
	}

	// Encrypt the AES key with the RSA public key
	encryptedAESKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, aesKey, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Error encrypting AES key")
		return body
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating AES cipher block")
		return body
	}

	// Create a new GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating GCM")
		return body
	}

	// Create a new nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal().Err(err).Msg("Error creating nonce")
		return body
	}

	// Encrypt the data using AES-GCM
	ciphertext := gcm.Seal(nonce, nonce, body, nil)

	// Return the RSA-encrypted AES key and the AES-encrypted data, both base64-encoded and separated by a colon
	return []byte(base64.StdEncoding.EncodeToString(encryptedAESKey) + ":" + base64.StdEncoding.EncodeToString(ciphertext))

}

var once = sync.Once{}

func loadPublicKey(path string) {
	publicKeyData, err := os.ReadFile(path)
	if err != nil {
		log.Error().Err(err).Msg("Error reading public key")
		return
	}

	// Parse the public key
	block, _ := pem.Decode(publicKeyData)
	if block == nil {
		log.Error().Err(err).Msg("Error decoding public key")
		return
	}

	rsaPublicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing public key")
		return
	}

	var ok bool
	publicKey, ok = rsaPublicKey.(*rsa.PublicKey)
	if !ok {
		log.Error().Err(err).Msg("Error casting public key")
		return
	}
}
