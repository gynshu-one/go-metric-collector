package main

import (
	"fmt"
	ag "github.com/gynshu-one/go-metric-collector/internal/controller/http/agent"
	"github.com/gynshu-one/go-metric-collector/internal/domain/service"
	"github.com/rs/zerolog/log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

var (
	agent        ag.Handler
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	agent = ag.NewAgent(service.NewMemService())
	log.Info().Msg("Agent started")
	time.Sleep(1 * time.Second)
	f, err := os.Create("server_mem.prof")
	if err != nil {
		log.Fatal().Err(err).Msg("could not create memory profile")
	}
	runtime.GC()
	if err = pprof.WriteHeapProfile(f); err != nil {
		log.Fatal().Err(err).Msg("could not write memory profile")
	}
	err = f.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("could not close memory profile")
	}
	agent.Start()
}
