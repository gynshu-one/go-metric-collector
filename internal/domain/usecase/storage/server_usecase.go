package storage

import (
	"context"
	"encoding/json"
	"fmt"
	config "github.com/gynshu-one/go-metric-collector/internal/config/server"
	"github.com/gynshu-one/go-metric-collector/internal/db_adapters"
	"github.com/gynshu-one/go-metric-collector/internal/domain/entity"
	"github.com/gynshu-one/go-metric-collector/internal/domain/service"
	"log"
	"os"
	"time"
)

type ServerStorage interface {
	service.MemStorage
	Dump(context.Context)
	Restore(context.Context)
}
type serverUseCase struct {
	service.MemStorage
	dbAdapter db_adapters.DbAdapter
}

func NewServerUseCase(MemStorage service.MemStorage, dbAdapter db_adapters.DbAdapter) *serverUseCase {
	s := &serverUseCase{
		MemStorage: MemStorage,
		dbAdapter:  dbAdapter,
	}
	s.filesDaemon()
	return s
}

func (S *serverUseCase) filesDaemon() {
	if config.GetConfig().Server.Restore {
		S.Restore(context.Background())
	}
	if config.GetConfig().Server.StoreInterval != 0 {
		ticker := time.NewTicker(config.GetConfig().Server.StoreInterval)
		go func() {
			for {
				t := <-ticker.C
				S.Dump(context.Background())
				fmt.Println("Saved to file at", t)
			}
		}()
	}
}
func (S *serverUseCase) Dump(ctx context.Context) {
	if config.GetConfig().Database.Address != "" {
		S.toDB(ctx)
	} else {
		S.toFile()
	}

}
func (S *serverUseCase) Restore(ctx context.Context) {
	if config.GetConfig().Database.Address != "" {
		S.fromDB(ctx)
	} else {
		S.fromFile()
	}
}
func (S *serverUseCase) fromDB(ctx context.Context) {
	metrics, err := S.dbAdapter.GetMetrics(ctx)
	if err != nil {
		return
	}
	for _, m := range metrics {
		S.Set(m)
	}
}

func (S *serverUseCase) toDB(ctx context.Context) {
	allMetrics := make([]*entity.Metrics, 0)
	S.ApplyToAll(func(metrics *entity.Metrics) {
		allMetrics = append(allMetrics, metrics)
	})
	err := S.dbAdapter.StoreMetrics(ctx, allMetrics)
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
	allMetrics := make([]*entity.Metrics, 0)
	S.ApplyToAll(func(metrics *entity.Metrics) {
		allMetrics = append(allMetrics, metrics)
	})
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
