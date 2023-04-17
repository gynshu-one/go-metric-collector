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
}

var instance *config
var once sync.Once

func GetConfig() *config {
	once.Do(func() {
		instance = &config{}
		// Order matters if we want to prioritize ENV over flags
		instance.readServerFlags()
		instance.readOs()
		// Then init files
		instance.initFiles()
		log.Debug().Interface("config", instance).Msg("Server started with configs")
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
	if v.Get("STORE_INTERVAL") != nil {
		config.Server.StoreInterval = v.GetDuration("STORE_INTERVAL")
	}
	if v.Get("STORE_FILE") != nil {
		config.Server.StoreFile = v.GetString("STORE_FILE")
	}
	if v.Get("RESTORE") != nil {
		config.Server.Restore = v.GetBool("RESTORE")
	}
	if v.Get("KEY") != nil {
		config.Key = v.GetString("KEY")
	}
	if v.Get("DATABASE_DSN") != nil {
		config.Database.Address = v.GetString("DATABASE_DSN")
	}
	//config.Server.Address = "http://" + config.Server.Address
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
func (config *config) readServerFlags() {
	// read flags
	appFlags := flag.NewFlagSet("go-metric-collector", flag.ContinueOnError)

	appFlags.StringVar(&config.Server.Address, "a", "localhost:8080", "server address")
	appFlags.DurationVar(&config.Server.StoreInterval, "i", 300*time.Second, "store interval")
	appFlags.StringVar(&config.Server.StoreFile, "f", "/tmp/devops-metrics-db.json", "store file")
	appFlags.StringVar(&config.Key, "k", "", "hash key")
	appFlags.BoolVar(&config.Server.Restore, "r", true, "restore")
	appFlags.StringVar(&config.Database.Address, "d", "", "DB address")
	err := appFlags.Parse(os.Args[1:])
	if err != nil {
		log.Debug().Err(err).Msg("Failed to parse flags")
	}
}
