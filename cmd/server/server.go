package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gynshu-one/go-metric-collector/internal/configs"
	"github.com/gynshu-one/go-metric-collector/internal/handlers"
	"github.com/gynshu-one/go-metric-collector/internal/middlewares"
	"github.com/gynshu-one/go-metric-collector/internal/routers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	router  *gin.Engine
	handler *handlers.ServerHandler
	server  *http.Server
)

func init() {
	// Order matters if we want to prioritize ENV over flags
	configs.CFG.ReadServerFlags()
	configs.CFG.ReadOs()
	// Then init files
	configs.CFG.InitFiles()
	//gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
	handler = handlers.NewServerHandler()
	// I don't know if MiscDecompress() middleware even required for this increment
	router.Use(cors.Default(), middlewares.MiscDecompress(), gzip.Gzip(gzip.DefaultCompression))
	routers.MetricsRoute(router, handler)
	// These two lines written to pass autotests (wrong code, redirect)
	// -------------------------------
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = true
	// -------------------------------
	server = &http.Server{
		Addr:    configs.CFG.Address,
		Handler: router,
	}
}

// Server that receives runtime metrics from the agent. with a configurable pollInterval.
func main() {
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	handler.Memory.Dump()
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
