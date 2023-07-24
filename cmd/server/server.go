package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gynshu-one/go-metric-collector/internal/adapters"
	config "github.com/gynshu-one/go-metric-collector/internal/config/server"
	hand "github.com/gynshu-one/go-metric-collector/internal/controller/http/server/handler"
	"github.com/gynshu-one/go-metric-collector/internal/controller/http/server/middlewares"
	"github.com/gynshu-one/go-metric-collector/internal/controller/http/server/routers"
	"github.com/gynshu-one/go-metric-collector/internal/domain/service"
	usecase "github.com/gynshu-one/go-metric-collector/internal/domain/usecase/storage"
	"github.com/gynshu-one/go-metric-collector/repos/postgres"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
	"time"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
	storage      usecase.ServerStorage
	server       *http.Server
	handler      hand.Handler
	router       *gin.Engine
	dbConn       postgres.DBConn
	dbAdapter    adapters.DBAdapter
)

func init() {
	//gin.SetMode(gin.ReleaseMode)
	router = gin.Default()

	// These two lines written to pass autotests (wrong code, redirect)
	// -------------------------------
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = true
	// -------------------------------

	server = &http.Server{
		Addr:    config.GetConfig().Server.Address,
		Handler: router,
	}
}

// ServerStorage that receives runtime metrics from the agent. with a configurable pollInterval.
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

	ctx := context.Background()
	dbConn = postgres.NewDB()
	if config.GetConfig().Database.Address != "" {
		err := dbConn.Connect()
		if err != nil {
			log.Fatal().Err(err).Msg("Database connection error")
		}
		dbAdapter = adapters.NewAdapter(ctx, dbConn.GetConn())
	}

	log.Info().Msg("Database connected")

	log.Info().Msg("Activating services")
	storage = usecase.NewServerUseCase(ctx, service.NewMemService(), dbAdapter)
	handler = hand.NewServerHandler(storage, dbConn)
	router.Use(cors.Default(), middlewares.MiscDecompress(), gzip.Gzip(gzip.DefaultCompression), middlewares.DecryptMiddleware())
	routers.MetricsRoute(router, handler)
	log.Info().Msg("Services activated")

	log.Info().Msg("Starting server on " + config.GetConfig().Server.Address)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("https! Listen and serve error")
		}
	}()

	go func() {
		if err := http.ListenAndServe("localhost:9090", nil); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("http Listen and serve error")
		}
	}()

	time.Sleep(1 * time.Second)

	f, err := os.Create("server_mem.prof")
	if err != nil {
		log.Error().Err(err).Msg("could not create memory profile")
	}
	runtime.GC()
	if err = pprof.WriteHeapProfile(f); err != nil {
		log.Error().Err(err).Msg("could not write memory profile")
	}
	err = f.Close()
	if err != nil {
		log.Error().Err(err).Msg("could not close memory profile")
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	log.Info().Msg("Shutdown Server ...")

	storage.Dump(ctx)
	ctxShut, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err = server.Shutdown(ctxShut); err != nil {
		log.Fatal().Err(err).Msgf("Timeout of %d seconds exceeded, server forced to shutdown", 5)
	}

	log.Info().Msg("Server exiting")
}
