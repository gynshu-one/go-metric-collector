package main

import (
	"github.com/gynshu-one/go-metric-collector/internal/configs"
	"github.com/gynshu-one/go-metric-collector/internal/handlers"
)

func main() {
	configs.CFG.LoadConfig()
	configs.CFG.Address = "http://" + configs.CFG.Address
	agent := handlers.NewAgent(configs.CFG.PollInterval,
		configs.CFG.ReportInterval,
		configs.CFG.Address)

	// Start the agent
	agent.Poll()
}
