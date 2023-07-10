// Package storage contains all the storage implementations of the domain
// a little bit modified version of storage service
package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gynshu-one/go-metric-collector/internal/adapters"
	config "github.com/gynshu-one/go-metric-collector/internal/config/server"
	"github.com/gynshu-one/go-metric-collector/internal/domain/entity"
	"github.com/gynshu-one/go-metric-collector/internal/domain/service"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"sync"
	"time"
)

// ServerStorage is an interface for server storage
// It is used in server use case and contains all the methods of
// MemStorage interface and some additional methods such as Dump and Restore, SetFltPrc and GetFltPrc
// Which are specific for server side
type ServerStorage interface {
	service.MemStorage
	Dump(context.Context)
	Restore(context.Context)
	SetFltPrc(name, p string)
	GetFltPrc(name string) int
}
type serverUseCase struct {
	service.MemStorage
	dbAdapter adapters.DBAdapter
	// fltPrecision is for autotests iter3
	fltPrecision sync.Map
}

// NewServerUseCase creates new server storage, context is for filesDaemon
// Which is used to ether restore previous state from file
// or DB or to dump current state to file or DB (separate goroutine)
func NewServerUseCase(ctx context.Context, MemStorage service.MemStorage, dbAdapter adapters.DBAdapter) *serverUseCase {
	s := &serverUseCase{
		MemStorage:   MemStorage,
		dbAdapter:    dbAdapter,
		fltPrecision: sync.Map{},
	}
	log.Info().Msg("Server storage initialized")
	s.filesDaemon(ctx)
	return s
}

// SetFltPrc sets precision for float metrics, it is used in iter3
// to pass autotests
func (S *serverUseCase) SetFltPrc(name, p string) {
	precision := strings.Split(p, ".")
	if len(precision) < 2 {
		S.fltPrecision.Store(name, 0)
		return
	}
	S.fltPrecision.Store(name, len(precision[1]))
}

// GetFltPrc returns precision for float metrics
func (S *serverUseCase) GetFltPrc(name string) int {
	if v, ok := S.fltPrecision.Load(name); ok {
		return v.(int)
	}
	return 0
}
func (S *serverUseCase) filesDaemon(ctx context.Context) {
	if config.GetConfig().Server.Restore {
		S.Restore(ctx)
	}
	if config.GetConfig().Server.StoreInterval != 0 && config.GetConfig().Database.Address == "" {
		go func() {
			ticker := time.NewTicker(config.GetConfig().Server.StoreInterval)
			for {
				t := <-ticker.C
				S.Dump(ctx)
				fmt.Println("Saved to file at", t)
			}
		}()
	}
}

// Dump dumps current state to file or DB
func (S *serverUseCase) Dump(ctx context.Context) {
	if config.GetConfig().Database.Address != "" {
		S.toDB(ctx)
	} else {
		S.toFile()
	}

}

// Restore  dumps current state to file or DB
func (S *serverUseCase) Restore(ctx context.Context) {
	if config.GetConfig().Database.Address != "" {
		S.fromDB(ctx)
	} else {
		S.fromFile()
	}
}
func (S *serverUseCase) fromDB(ctx context.Context) {
	log.Info().Msg("Restoring previous state from DB...")
	metrics, err := S.dbAdapter.GetMetrics(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error restoring from DB")
		return
	}
	for _, m := range metrics {
		S.Set(m)
	}
	log.Info().Msg("Successfully restored from DB")
}

func (S *serverUseCase) toDB(ctx context.Context) {
	allMetrics := S.GetAll()
	err := S.dbAdapter.StoreMetrics(ctx, allMetrics)
	if err != nil {
		log.Error().Err(err).Msg("Error storing to DB")
		return
	}
	log.Info().Msg("Successfully stored to DB")
}

func (S *serverUseCase) fromFile() {
	log.Info().Msg("Restoring previous state from file...")
	file, err := os.OpenFile(config.GetConfig().Server.StoreFile, os.O_RDONLY, 0666)
	if err != nil {
		log.Warn().Err(err).Msg("Error opening file to restore")
		return
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Trace().Err(err).Msg("Error closing file")
		}
	}()
	var metrics []*entity.Metrics
	err = json.NewDecoder(file).Decode(&metrics)
	if err != nil {
		log.Error().Err(err).Msg("Error decoding json may be file is empty:")
		return
	}
	if metrics == nil {
		log.Warn().Msg("File is empty")
		return
	}
	for _, m := range metrics {
		S.Set(m)
	}
	log.Info().Msg("Successfully restored from file")
	metrics = nil
}

func (S *serverUseCase) toFile() {
	allMetrics := S.GetAll()
	jsonData, err := json.Marshal(allMetrics)
	if err != nil {
		log.Error().Err(err).Msg("Error marshaling metrics to json")
	}
	err = os.WriteFile(config.GetConfig().Server.StoreFile, jsonData, 0644)
	if err != nil {
		log.Error().Err(err).Msg("Error writing to file")
	}
}
