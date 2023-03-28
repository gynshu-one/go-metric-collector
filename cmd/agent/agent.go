package main

import (
	"github.com/fatih/color"
	"github.com/gynshu-one/go-metric-collector/internal/configs"
	"github.com/gynshu-one/go-metric-collector/internal/handlers"
)

func init() {
	// Order matters if we want to prioritize ENV over flags
	configs.CFG.ReadAgentFlags()
	configs.CFG.ReadOs()
	// Then init files
	configs.CFG.InitFiles()
	configs.CFG.Address = "http://" + configs.CFG.Address
	color.Cyan("Configs: %+v", configs.CFG)

}
func main() {
	agent := handlers.NewAgent(configs.CFG.PollInterval,
		configs.CFG.ReportInterval,
		configs.CFG.Address)

	// Start the agent
	agent.Poll()
}
