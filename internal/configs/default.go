package configs

import (
	"fmt"
	"github.com/gynshu-one/go-metric-collector/internal/tools"
	"github.com/spf13/viper"
	"os"
	"path"
	"strings"
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

func (config *Config) LoadConfig() {
	// load config from environment variables
	v := viper.New()
	dir := tools.GetProjectRoot()
	v.SetConfigName("app")
	v.SetConfigType("env")
	v.AddConfigPath(dir)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error reading config: %s", err))
	}
	// Retrieve the value using the key
	err = v.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("error decoding config: %s", err))
	}

	// I know this is bad I could use v.AutomaticEnv(). What if os env does not have all vars
	// but some of them
	if os.Getenv("ADDRESS") != "" {
		config.Address = os.Getenv("ADDRESS")
	}
	if os.Getenv("POLL_INTERVAL") != "" {
		config.PollInterval = tools.ParseDuration(os.Getenv("POLL_INTERVAL"))
	}
	if os.Getenv("REPORT_INTERVAL") != "" {
		config.ReportInterval = tools.ParseDuration(os.Getenv("REPORT_INTERVAL"))
	}
	if os.Getenv("STORE_INTERVAL") != "" {
		config.StoreInterval = tools.ParseDuration(os.Getenv("STORE_INTERVAL"))
	}
	if os.Getenv("STORE_FILE") != "" {
		config.StoreFile = os.Getenv("STORE_FILE")
	}
	if os.Getenv("RESTORE") != "" {
		config.Restore = strings.Contains(os.Getenv("RESTORE"), "true")
	}
	config.StoreFile = dir + config.StoreFile
	// get dir of the file
	dr := path.Dir(config.StoreFile)
	// check if dir exists
	if _, err = os.Stat(dr); os.IsNotExist(err) {
		// create dir
		err = os.MkdirAll(dr, os.ModePerm)
		if err != nil {
			panic(fmt.Errorf("error creating dir: %s", err))
		}
	}
}
