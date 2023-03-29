package configs

import (
	"flag"
	"fmt"
	"github.com/gynshu-one/go-metric-collector/internal/tools"
	"github.com/spf13/viper"
	"os"
	"path"
	"time"
)

type Config struct {
	Address        string        `mapstructure:"ADDRESS"`
	PollInterval   time.Duration `mapstructure:"POLL_INTERVAL"`
	ReportInterval time.Duration `mapstructure:"REPORT_INTERVAL"`
	StoreInterval  time.Duration `mapstructure:"STORE_INTERVAL"`
	StoreFile      string        `mapstructure:"STORE_FILE"`
	Restore        bool          `mapstructure:"RESTORE"`
}

var CFG = &Config{}

// ReadOs reads config from environment variables
// This func will replace Config parameters if any presented in os environment vars
func (config *Config) ReadOs() {
	// load config from environment variables
	v := viper.New()
	v.AutomaticEnv()
	if v.Get("ADDRESS") != nil {
		config.Address = v.GetString("ADDRESS")
	}
	if v.Get("POLL_INTERVAL") != nil {
		config.PollInterval = v.GetDuration("POLL_INTERVAL")
	}
	if v.Get("REPORT_INTERVAL") != nil {
		config.ReportInterval = v.GetDuration("REPORT_INTERVAL")
	}
	if v.Get("STORE_INTERVAL") != nil {
		config.StoreInterval = v.GetDuration("STORE_INTERVAL")
	}
	if v.Get("STORE_FILE") != nil {
		config.StoreFile = v.GetString("STORE_FILE")
	}
	if v.Get("RESTORE") != nil {
		config.Restore = v.GetBool("RESTORE")
	}
}

// InitFiles creates all necessary files and folders for server storage
func (config *Config) InitFiles() {
	dir := tools.GetProjectRoot()
	// Make temp files dir absolute
	config.StoreFile = dir + config.StoreFile
	// get dir of the file
	dr := path.Dir(config.StoreFile)
	// check if dir exists
	if _, err := os.Stat(dr); os.IsNotExist(err) {
		// create dir
		err = os.MkdirAll(dr, os.ModePerm)
		if err != nil {
			panic(fmt.Errorf("error creating dir: %s", err))
		}
	}
}

// ReadServerFlags reads config from flags Run this first
func (config *Config) ReadServerFlags() {
	// read flags
	flag.StringVar(&config.Address, "a", "localhost:8080", "server address")
	flag.DurationVar(&config.StoreInterval, "i", 300*time.Second, "store interval")
	flag.StringVar(&config.StoreFile, "f", "/tmp/devops-metrics-db.json", "store file")
	flag.BoolVar(&config.Restore, "r", true, "restore")
	flag.Parse()
}

// ReadAgentFlags separate function required bec of similar variable names required for agent and server
func (config *Config) ReadAgentFlags() {
	// read flags
	flag.StringVar(&config.Address, "a", "localhost:8080", "server address")
	flag.DurationVar(&config.PollInterval, "p", 2*time.Second, "poll interval")
	flag.DurationVar(&config.ReportInterval, "r", 10*time.Second, "report interval")
	flag.Parse()
}
