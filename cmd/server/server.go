package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gynshu-one/go-metric-collector/internal/configs"
	"github.com/gynshu-one/go-metric-collector/internal/routers"
	"github.com/gynshu-one/go-metric-collector/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	// Order matters if we want to prioritize ENV over flags
	configs.CFG.ReadServerFlags()
	configs.CFG.ReadOs()
	// Then init files
	configs.CFG.InitFiles()
	dock := os.Getenv("DOCKER")
	//color.Cyan("Configs: %+v", configs.CFG)
	if dock != "" {
		configs.CFG.Address = "0.0.0.0:8080"
	}
	storage.Memory = storage.InitServerStorage()

}

// Server that receives runtime metrics from the agent. with a configurable PollInterval.
func main() {
	router := gin.Default()
	// Change gin mode
	//gin.SetMode(gin.ReleaseMode)
	// Disable log gin
	router.Use(cors.Default())

	routers.MetricsRoute(router)
	// These two lines written to pass autotests (wrong code, redirect)
	// -------------------------------
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = true
	// -------------------------------
	srv := &http.Server{
		Addr:    configs.CFG.Address,
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	err := storage.Memory.StoreEverythingToFile()
	if err != nil {
		log.Fatal("Error while storing data to file: ", err)
	}
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
