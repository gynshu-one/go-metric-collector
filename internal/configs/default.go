package configs

import (
	"fmt"
	"github.com/gynshu-one/go-metric-collector/internal/tools"
	"github.com/mitchellh/mapstructure"
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

func (config *Config) LoadConfig() {
	// load config from environment variables
	v := viper.New()
	v.AutomaticEnv()
	dir := tools.GetProjectRoot()
	err := mapstructure.Decode(v.AllSettings(), &config)
	if err != nil {
		panic(fmt.Errorf("error decoding config: %s", err))
	}
	if config.Address == "" || config.PollInterval == 0 || config.ReportInterval == 0 {
		// load
		v.SetConfigName("app")
		v.SetConfigType("env")
		v.AddConfigPath(dir)
		err = v.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("error reading config: %s", err))
		}
		// Retrieve the value using the key

		err = v.Unmarshal(&config)
		if err != nil {
			panic(fmt.Errorf("error decoding config: %s", err))
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
}
