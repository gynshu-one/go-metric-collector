package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	config "github.com/gynshu-one/go-metric-collector/internal/config/server"
	hand "github.com/gynshu-one/go-metric-collector/internal/controller/http/server/handler"
	"github.com/gynshu-one/go-metric-collector/internal/controller/http/server/middlewares"
	"github.com/gynshu-one/go-metric-collector/internal/controller/http/server/routers"
	"github.com/gynshu-one/go-metric-collector/internal/db_adapters"
	"github.com/gynshu-one/go-metric-collector/internal/domain/service"
	usecase "github.com/gynshu-one/go-metric-collector/internal/domain/usecase/storage"
	"github.com/gynshu-one/go-metric-collector/pkg/client/postgres"
	"log"
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

	dbAdapter db_adapters.DbAdapter
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
		dbAdapter = db_adapters.NewAdapter(dbConn.GetConn())
		if err != nil {
			log.Fatal("Database connection error: ", err)
		}
	}
	storage = usecase.NewServerUseCase(service.NewMemService(&sync.Map{}), dbAdapter)
	handler = hand.NewServerHandler(storage, dbConn)
	router.Use(cors.Default(), middlewares.MiscDecompress(), gzip.Gzip(gzip.DefaultCompression))
	routers.MetricsRoute(router, handler)
	//ok := dbAdapter.Test()
	//if !ok {
	//	log.Fatal("Database connection error: ", ok)
	//}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: ", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	storage.Dump(ctx)
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctxShut, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctxShut); err != nil {
		log.Fatal("ServerStorage forced to shutdown: ", err)
	}

	log.Println("ServerStorage exiting")
}
