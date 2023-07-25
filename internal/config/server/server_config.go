package server

import (
	"flag"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"path"
	"sync"
	"time"
)

type config struct {
	Key    string `mapstructure:"KEY"`
	Server struct {
		Address       string        `mapstructure:"ADDRESS"`
		StoreInterval time.Duration `mapstructure:"STORE_INTERVAL"`
		StoreFile     string        `mapstructure:"STORE_FILE"`
		Restore       bool          `mapstructure:"RESTORE"`
	}
	Database struct {
		Address string `mapstructure:"DATABASE_DSN"`
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
		flags := readServerFlags()
		smartSet(flags, prior)
		if prior.CfgPath != "" {
			json := readConfigJSON(prior.CfgPath)
			smartSet(json, prior)
		}
		instance = prior
		// Then init files
		instance.initFiles()
		log.Debug().Interface("config", instance).Msg("Server started with configs")
	})
	return instance
}

// readOs reads config from environment variables
// This func will replace Config parameters if any presented in os environment vars
func readOs() *config {
	var cfg config
	v := viper.New()
	v.AutomaticEnv()
	if v.Get("ADDRESS") != nil {
		cfg.Server.Address = v.GetString("ADDRESS")
	}
	if v.Get("STORE_INTERVAL") != nil {
		cfg.Server.StoreInterval = v.GetDuration("STORE_INTERVAL")
	}
	if v.Get("STORE_FILE") != nil {
		cfg.Server.StoreFile = v.GetString("STORE_FILE")
	}
	if v.Get("RESTORE") != nil {
		cfg.Server.Restore = v.GetBool("RESTORE")
	}
	if v.Get("KEY") != nil {
		cfg.Key = v.GetString("KEY")
	}
	if v.Get("DATABASE_DSN") != nil {
		cfg.Database.Address = v.GetString("DATABASE_DSN")
	}
	if v.Get("CRYPTO_KEY") != nil {
		cfg.CryptoKey = v.GetString("CRYPTO_KEY")
	}
	if v.Get("CONFIG") != nil {
		cfg.CfgPath = v.GetString("CONFIG")
	}
	return &cfg
}

// initFiles creates all necessary files and folders for server storage
func (config *config) initFiles() {
	// get dir of the file
	dr := path.Dir(config.Server.StoreFile)
	// check if dir exists
	if _, err := os.Stat(dr); os.IsNotExist(err) {
		// create dir
		err = os.MkdirAll(dr, os.ModePerm)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create dir for server storage")
		}
	}
}

// readServerFlags reads config from flags Run this first
func readServerFlags() *config {
	var cfg config
	// read flags
	appFlags := flag.NewFlagSet("go-metric-collector", flag.ContinueOnError)

	appFlags.StringVar(&cfg.Server.Address, "a", "localhost:8080", "server address")
	appFlags.DurationVar(&cfg.Server.StoreInterval, "i", 1*time.Second, "store interval")
	appFlags.StringVar(&cfg.Server.StoreFile, "f", "/tmp/devops-metrics-db.json", "store file")
	appFlags.StringVar(&cfg.Key, "k", "", "hash key")
	appFlags.BoolVar(&cfg.Server.Restore, "r", true, "restore")
	appFlags.StringVar(&cfg.Database.Address, "d", "", "DB address")
	appFlags.StringVar(&cfg.CryptoKey, "crypto-key", "", "crypto key")
	appFlags.StringVar(&cfg.CfgPath, "c", "config", "config file")

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
	if new == nil || old == nil {
		return
	}

	if old.Key == "" {
		old.Key = new.Key
	}
	if old.Server.Address == "" {
		old.Server.Address = new.Server.Address
	}
	if old.Server.StoreInterval == 0 {
		old.Server.StoreInterval = new.Server.StoreInterval
	}
	if old.Server.StoreFile == "" {
		old.Server.StoreFile = new.Server.StoreFile
	}
	if old.Database.Address == "" {
		old.Database.Address = new.Database.Address
	}
	if old.CryptoKey == "" {
		old.CryptoKey = new.CryptoKey
	}
	if old.CfgPath == "" {
		old.CfgPath = new.CfgPath
	}
}
