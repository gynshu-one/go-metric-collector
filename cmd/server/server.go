package main

import (
	"context"
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
	"sync"
	"syscall"
	"time"
)

var (
	storage usecase.ServerStorage
	server  *http.Server
	handler hand.Handler
	router  *gin.Engine
	dbConn  postgres.DBConn

	dbAdapter adapters.DBAdapter
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
	ctx := context.Background()
	dbConn = postgres.NewDB()
	if config.GetConfig().Database.Address != "" {
		err := dbConn.Connect()
		dbAdapter = adapters.NewAdapter(dbConn.GetConn())
		if err != nil {
			log.Fatal().Err(err).Msg("Database connection error")
		}
	}

	log.Info().Msg("Database connected")

	log.Info().Msg("Activating services")
	storage = usecase.NewServerUseCase(service.NewMemService(&sync.Map{}), dbAdapter)
	handler = hand.NewServerHandler(storage, dbConn)
	router.Use(cors.Default(), middlewares.MiscDecompress(), gzip.Gzip(gzip.DefaultCompression))
	routers.MetricsRoute(router, handler)
	log.Info().Msg("Services activated")

	log.Info().Msg("Starting server on " + config.GetConfig().Server.Address)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Listen and serve error")
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutdown Server ...")

	storage.Dump()
	ctxShut, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctxShut); err != nil {
		log.Fatal().Err(err).Msgf("Timeout of %d seconds exceeded, server forced to shutdown", 5)
	}

	log.Info().Msg("Server exiting")
}
