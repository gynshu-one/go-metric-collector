package main

import (
	"context"
	"fmt"
	ag "github.com/gynshu-one/go-metric-collector/internal/controller/http/agent"
	"github.com/gynshu-one/go-metric-collector/internal/domain/service"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
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
	go agent.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	log.Info().Msg("Shutdown Agent ...")

	// run func with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	go func() {
		agent.Stop()
		cancel()
	}()
	select {
	case <-ctx.Done():
		log.Info().Msg("Agent stopped successfully, Exiting...")
	case <-time.After(15 * time.Second):
		log.Error().Msg("Agent shutdown timeout. Exiting...")
	}
}
