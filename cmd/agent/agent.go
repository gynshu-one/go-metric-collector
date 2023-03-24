package main

import (
	"flag"
	"github.com/gynshu-one/go-metric-collector/internal/handlers"
	"time"
)

func main() {
	pollInterval := flag.Duration("pollInterval", 2*time.Second, "Poll interval for the server")
	reportInterval := flag.Duration("reportInterval", 10*time.Second, "Report interval for the agent")
	serverAddr := flag.String("serverAddr", "localhost:8080", "Server address for the agent")
	// Create a new agent
	*serverAddr = "http://" + *serverAddr
	agent := handlers.NewAgent(pollInterval, reportInterval, serverAddr)

	// Start the agent
	agent.Poll()
}
