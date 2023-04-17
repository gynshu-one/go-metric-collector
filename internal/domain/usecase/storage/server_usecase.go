package storage

import (
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

type ServerStorage interface {
	service.MemStorage
	Dump()
	Restore()
	SetFltPrc(name, p string)
	GetFltPrc(name string) int
}
type serverUseCase struct {
	service.MemStorage
	dbAdapter adapters.DBAdapter
	// fltPrecision is for autotests iter3
	fltPrecision sync.Map
}

func NewServerUseCase(MemStorage service.MemStorage, dbAdapter adapters.DBAdapter) *serverUseCase {
	s := &serverUseCase{
		MemStorage:   MemStorage,
		dbAdapter:    dbAdapter,
		fltPrecision: sync.Map{},
	}
	log.Info().Msg("Server storage initialized")
	s.filesDaemon()
	return s
}
func (S *serverUseCase) SetFltPrc(name, p string) {
	precision := strings.Split(p, ".")
	if len(precision) < 2 {
		S.fltPrecision.Store(name, 0)
		return
	}
	S.fltPrecision.Store(name, len(precision[1]))
}
func (S *serverUseCase) GetFltPrc(name string) int {
	if v, ok := S.fltPrecision.Load(name); ok {
		return v.(int)
	}
	return 0
}
func (S *serverUseCase) filesDaemon() {
	if config.GetConfig().Server.Restore {
		S.Restore()
	}
	if config.GetConfig().Server.StoreInterval != 0 && config.GetConfig().Database.Address == "" {
		go func() {
			ticker := time.NewTicker(config.GetConfig().Server.StoreInterval)
			for {
				t := <-ticker.C
				S.Dump()
				fmt.Println("Saved to file at", t)
			}
		}()
	}
}
func (S *serverUseCase) Dump() {
	if config.GetConfig().Database.Address != "" {
		S.toDB()
	} else {
		S.toFile()
	}

}
func (S *serverUseCase) Restore() {
	if config.GetConfig().Database.Address != "" {
		S.fromDB()
	} else {
		S.fromFile()
	}
}
func (S *serverUseCase) fromDB() {
	log.Info().Msg("Restoring previous state from DB...")
	metrics, err := S.dbAdapter.GetMetrics()
	if err != nil {
		log.Error().Err(err).Msg("Error restoring from DB")
		return
	}
	for _, m := range metrics {
		S.Set(m)
	}
	log.Info().Msg("Successfully restored from DB")
}

func (S *serverUseCase) toDB() {
	allMetrics := S.GetAll()
	err := S.dbAdapter.StoreMetrics(allMetrics)
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
		log.Error().Err(err).Msg("Error opening file")
		return
	}
	defer file.Close()
	var metrics []*entity.Metrics
	err = json.NewDecoder(file).Decode(&metrics)
	if err != nil {
		log.Error().Err(err).Msg("Error decoding json may be file is empty:")
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
