package agent

import (
	"flag"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"sync"
	"time"
)

type config struct {
	Key   string `mapstructure:"KEY"`
	Agent struct {
		PollInterval   time.Duration `mapstructure:"POLL_INTERVAL"`
		ReportInterval time.Duration `mapstructure:"REPORT_INTERVAL"`
		RateLimit      int           `mapstructure:"RATE_LIMIT"`
	}
	Server struct {
		Address string `mapstructure:"ADDRESS"`
	}
	CryptoKey string `mapstructure:"CRYPTO_KEY"`
}

var instance *config
var once sync.Once

func GetConfig() *config {
	once.Do(func() {
		instance = &config{}
		// Order matters if we want to prioritize ENV over flags
		instance.readAgentFlags()
		instance.readOs()
		log.Debug().Interface("config", instance).Msg("Agent started with configs")
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
	if v.Get("RATE_LIMIT") != nil {
		config.Agent.RateLimit = v.GetInt("RATE_LIMIT")
	}
	if v.Get("CRYPTO_KEY") != nil {
		config.CryptoKey = v.GetString("CRYPTO_KEY")
	}
	config.Server.Address = "http://" + config.Server.Address
}

// readAgentFlags separate function required bec of similar variable names required for agent and server
func (config *config) readAgentFlags() {
	// read flags
	appFlags := flag.NewFlagSet("go-metric-collector", flag.ContinueOnError)

	appFlags.StringVar(&config.Server.Address, "a", "localhost:8080", "server address")
	appFlags.StringVar(&config.Key, "k", "", "hash key")
	appFlags.DurationVar(&config.Agent.PollInterval, "p", 2*time.Second, "poll interval")
	appFlags.DurationVar(&config.Agent.ReportInterval, "r", 3*time.Second, "report interval")
	appFlags.IntVar(&config.Agent.RateLimit, "l", 2, "rate limit")
	appFlags.StringVar(&config.CryptoKey, "crypto-key", "", "crypto key")

	// Parse the flags using the new flag set
	err := appFlags.Parse(os.Args[1:])
	if err != nil {
		log.Debug().Err(err).Msg("Failed to parse flags")
	}
}
