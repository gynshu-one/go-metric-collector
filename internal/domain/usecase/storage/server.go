package storage

import (
	"encoding/json"
	config "github.com/gynshu-one/go-metric-collector/internal/config/server"
	"github.com/gynshu-one/go-metric-collector/internal/domain/entity"
	"github.com/gynshu-one/go-metric-collector/internal/domain/service"
	"log"
	"os"
)

type ServerStorage interface {
	service.MemStorage
	Dump()
	Restore()
}
type serverUseCase struct {
	service.MemStorage
}

func NewServerUseCase(MemStorage service.MemStorage) *serverUseCase {
	return &serverUseCase{
		MemStorage: MemStorage,
	}
}
func (S serverUseCase) Dump() {
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
	//path := configs.CFG.StoreFile
}
func (S serverUseCase) Restore() {
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
