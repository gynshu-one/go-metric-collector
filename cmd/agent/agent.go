package main

import (
	ag "github.com/gynshu-one/go-metric-collector/internal/controller/http/agent"
	"github.com/gynshu-one/go-metric-collector/internal/domain/service"
	"github.com/rs/zerolog/log"
	"sync"
)

var (
	agent   ag.Handler
	storage service.MemStorage
)

func main() {
	agent = ag.NewAgent(service.NewMemService(&sync.Map{}))
	log.Info().Msg("Agent started")
	agent.Start()
}
