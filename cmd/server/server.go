package main

import (
	"flag"
	"github.com/gynshu-one/go-metric-collector/internal/handlers"
)

// Serve that receives runtime metrics from the agent. with a configurable PollInterval.
func main() {
	// Create a flag for the poll interval
	addr := flag.String("addr", "localhost", "Address for the server")
	port := flag.String("port", "8080", "Port for the server")

	// Create a new server
	server := handlers.NewServer(addr, port)
	// Start the server
	server.Start()
}
