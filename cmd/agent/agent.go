package main

import (
	"github.com/gynshu-one/go-metric-collector/internal/configs"
	"github.com/gynshu-one/go-metric-collector/internal/handlers"
	"time"
)

func main() {
	configs.CFG.LoadConfig(".")
	configs.CFG.Address = "http://" + configs.CFG.Address
	agent := handlers.NewAgent(time.Duration(configs.CFG.PollInterval)*time.Second,
		time.Duration(configs.CFG.ReportInterval)*time.Second,
		configs.CFG.Address)

	// Start the agent
	agent.Poll()
}
