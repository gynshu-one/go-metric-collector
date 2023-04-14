package main

import (
	ag "github.com/gynshu-one/go-metric-collector/internal/controller/http/agent"
	"github.com/gynshu-one/go-metric-collector/internal/domain/service"
	"sync"
)

var (
	agent   ag.Handler
	storage service.MemStorage
)

func main() {
	agent = ag.NewAgent(service.NewMemService(&sync.Map{}))
	agent.Start()

}
