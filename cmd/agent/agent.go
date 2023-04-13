package main

import (
	ag "github.com/gynshu-one/go-metric-collector/internal/controller/http/agent"
	"github.com/gynshu-one/go-metric-collector/internal/domain/entity"
	"github.com/gynshu-one/go-metric-collector/internal/domain/service"
)

var (
	agent   ag.Handler
	storage service.MemStorage
)

func main() {
	agent = ag.NewAgent(service.NewMemService(make(map[string]*entity.Metrics)))
	agent.Start()
}
