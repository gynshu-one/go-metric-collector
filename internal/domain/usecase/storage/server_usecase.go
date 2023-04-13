package storage

import (
	"encoding/json"
	"fmt"
	"github.com/gynshu-one/go-metric-collector/internal/adapters"
	config "github.com/gynshu-one/go-metric-collector/internal/config/server"
	"github.com/gynshu-one/go-metric-collector/internal/domain/entity"
	"github.com/gynshu-one/go-metric-collector/internal/domain/service"
	"log"
	"os"
	"time"
)

type ServerStorage interface {
	service.MemStorage
	Dump()
	Restore()
}
type serverUseCase struct {
	service.MemStorage
	dbAdapter adapters.DBAdapter
}

func NewServerUseCase(MemStorage service.MemStorage, dbAdapter adapters.DBAdapter) *serverUseCase {
	s := &serverUseCase{
		MemStorage: MemStorage,
		dbAdapter:  dbAdapter,
	}
	s.filesDaemon()
	return s
}

func (S *serverUseCase) filesDaemon() {
	if config.GetConfig().Server.Restore {
		S.Restore()
	}
	if config.GetConfig().Server.StoreInterval != 0 && config.GetConfig().Database.Address == "" {
		ticker := time.NewTicker(config.GetConfig().Server.StoreInterval)
		go func() {
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
	metrics, err := S.dbAdapter.GetMetrics()
	if err != nil {
		return
	}
	for _, m := range metrics {
		S.Set(m)
	}
}

func (S *serverUseCase) toDB() {
	allMetrics := S.GetAll()
	err := S.dbAdapter.StoreMetrics(allMetrics)
	if err != nil {
		log.Fatal(err)
	}
}

func (S *serverUseCase) fromFile() {
	file, err := os.OpenFile(config.GetConfig().Server.StoreFile, os.O_RDONLY, 0666)
	if err != nil {
		log.Printf("Nothing to resore from storage file: %v", err)
		return
	}
	defer file.Close()
	var metrics []*entity.Metrics
	err = json.NewDecoder(file).Decode(&metrics)
	if err != nil {
		log.Printf("Error decoding json may be file is empty: %v", err)
		return
	}
	for _, m := range metrics {
		S.Set(m)
	}
	metrics = nil
}

func (S *serverUseCase) toFile() {
	allMetrics := S.GetAll()
	// save to jsonData file
	jsonData, err := json.Marshal(allMetrics)
	if err != nil {
		log.Fatal(err)
	}
	// save to file
	err = os.WriteFile(config.GetConfig().Server.StoreFile, jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
