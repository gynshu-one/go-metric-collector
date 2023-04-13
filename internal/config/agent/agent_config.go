package agent

import (
	"flag"
	"github.com/fatih/color"
	"github.com/spf13/viper"
	"sync"
	"time"
)

type config struct {
	Key   string `mapstructure:"KEY"`
	Agent struct {
		PollInterval   time.Duration `mapstructure:"POLL_INTERVAL"`
		ReportInterval time.Duration `mapstructure:"REPORT_INTERVAL"`
	}
	Server struct {
		Address string `mapstructure:"ADDRESS"`
	}
}

var instance *config
var once sync.Once

func GetConfig() *config {
	once.Do(func() {
		instance = &config{}
		// Order matters if we want to prioritize ENV over flags
		instance.readAgentFlags()
		instance.readOs()
		color.Cyan("Configs: %+v", instance)
	})
	return instance
}

// readOs reads config from environment variables
// This func will replace Config parameters if any presented in os environment vars
func (config *config) readOs() {
	// load config from environment variables
	v := viper.New()
	v.AutomaticEnv()
	if v.Get("ADDRESS") != nil {
		config.Server.Address = v.GetString("ADDRESS")
	}
	if v.Get("POLL_INTERVAL") != nil {
		config.Agent.PollInterval = v.GetDuration("POLL_INTERVAL")
	}
	if v.Get("REPORT_INTERVAL") != nil {
		config.Agent.ReportInterval = v.GetDuration("REPORT_INTERVAL")
	}
	if v.Get("KEY") != nil {
		config.Key = v.GetString("KEY")
	}
	config.Server.Address = "http://" + config.Server.Address
}

// readAgentFlags separate function required bec of similar variable names required for agent and server
func (config *config) readAgentFlags() {
	// read flags
	flag.CommandLine.Init("go-metric-collector", flag.ContinueOnError)
	flag.StringVar(&config.Server.Address, "a", "localhost:8080", "server address")
	flag.StringVar(&config.Key, "k", "/tmp/devops-metrics-db.json", "hash key")
	flag.DurationVar(&config.Agent.PollInterval, "p", 1*time.Second, "poll interval")
	flag.DurationVar(&config.Agent.ReportInterval, "r", 2*time.Second, "report interval")
	flag.Parse()
}
