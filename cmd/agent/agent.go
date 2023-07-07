package main

import (
	ag "github.com/gynshu-one/go-metric-collector/internal/controller/http/agent"
	"github.com/gynshu-one/go-metric-collector/internal/domain/service"
	"github.com/rs/zerolog/log"
	"os"
	"runtime"
	"runtime/pprof"
)

var (
	agent ag.Handler
)

func main() {
	f, err := os.Create("server_mem.prof")
	if err != nil {
		log.Fatal().Err(err).Msg("could not create memory profile")
	}
	runtime.GC()
	if err = pprof.WriteHeapProfile(f); err != nil {
		log.Fatal().Err(err).Msg("could not write memory profile")
	}
	_ = f.Close()

	agent = ag.NewAgent(service.NewMemService())
	log.Info().Msg("Agent started")
	agent.Start()
}
