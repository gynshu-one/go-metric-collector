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
	CfgPath   string `mapstructure:"CONFIG"`
}

var instance *config
var once sync.Once

func GetConfig() *config {
	once.Do(func() {
		instance = &config{}
		// Order matters if we want to prioritize ENV over flags
		prior := readOs()
		flags := readAgentFlags()
		smartSet(flags, prior)
		if prior.CfgPath != "" {
			json := readConfigJSON(prior.CfgPath)
			smartSet(json, prior)
		}
		prior.Server.Address = "http://" + prior.Server.Address
		instance = prior
		log.Debug().Interface("config", instance).Msg("Agent started with configs")
	})
	return instance
}

// readOs reads config from environment variables
// This func will replace Config parameters if any presented in os environment vars
func readOs() *config {
	var cfg config
	// load config from environment variables
	v := viper.New()
	v.AutomaticEnv()
	if v.GetString("KEY") != "" {
		cfg.Key = v.GetString("KEY")
	}
	if v.GetDuration("POLL_INTERVAL") != 0 {
		cfg.Agent.PollInterval = v.GetDuration("POLL_INTERVAL")
	}
	if v.GetDuration("REPORT_INTERVAL") != 0 {
		cfg.Agent.ReportInterval = v.GetDuration("REPORT_INTERVAL")
	}
	if v.GetInt("RATE_LIMIT") != 0 {
		cfg.Agent.RateLimit = v.GetInt("RATE_LIMIT")
	}
	if v.GetString("ADDRESS") != "" {
		cfg.Server.Address = v.GetString("ADDRESS")
	}
	if v.GetString("CRYPTO_KEY") != "" {
		cfg.CryptoKey = v.GetString("CRYPTO_KEY")
	}
	if v.GetString("CONFIG") != "" {
		cfg.CfgPath = v.GetString("CONFIG")
	}
	return &cfg
}

// readAgentFlags separate function required bec of similar variable names required for agent and server
func readAgentFlags() *config {
	var cfg config
	// read flags
	appFlags := flag.NewFlagSet("go-metric-collector", flag.ContinueOnError)

	appFlags.StringVar(&cfg.Server.Address, "a", "localhost:8080", "server address")
	appFlags.StringVar(&cfg.Key, "k", "", "hash key")
	appFlags.DurationVar(&cfg.Agent.PollInterval, "p", 2*time.Second, "poll interval")
	appFlags.DurationVar(&cfg.Agent.ReportInterval, "r", 10*time.Second, "report interval")
	appFlags.IntVar(&cfg.Agent.RateLimit, "l", 2, "rate limit")
	appFlags.StringVar(&cfg.CryptoKey, "crypto-key", "", "crypto key")
	appFlags.StringVar(&cfg.CfgPath, "c", "config", "config file")

	// Parse the flags using the new flag set
	err := appFlags.Parse(os.Args[1:])
	if err != nil {
		log.Debug().Err(err).Msg("Failed to parse flags")
	}
	return &cfg
}

func readConfigJSON(path string) *config {
	var cfg config
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(path)
	v.SetConfigType("json")
	err := v.ReadInConfig()
	if err != nil {
		log.Debug().Err(err).Msg("Failed to read config file")
	}
	err = v.Unmarshal(&cfg)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to unmarshal config file")
	}
	return &cfg
}

// smartSet sets new config values only if they are not empty
// meaning old is prioritized over new
func smartSet(new, old *config) {
	if old.Key == "" {
		old.Key = new.Key
	}
	if old.Agent.PollInterval == 0 {
		old.Agent.PollInterval = new.Agent.PollInterval
	}
	if old.Agent.ReportInterval == 0 {
		old.Agent.ReportInterval = new.Agent.ReportInterval
	}
	if old.Agent.RateLimit == 0 {
		old.Agent.RateLimit = new.Agent.RateLimit
	}
	if old.Server.Address == "" {
		old.Server.Address = new.Server.Address
	}
	if old.CryptoKey == "" {
		old.CryptoKey = new.CryptoKey
	}
	if old.CfgPath == "" {
		old.CfgPath = new.CfgPath
	}
}
