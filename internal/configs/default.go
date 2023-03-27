package configs

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

type Config struct {
	Address        string `mapstructure:"ADDRESS"`
	PollInterval   int    `mapstructure:"POLL_INTERVAL"`
	ReportInterval int    `mapstructure:"REPORT_INTERVAL"`
}

var CFG = &Config{}

func (config *Config) LoadConfig(file string) {
	// load config from environment variables
	v := viper.New()
	v.AutomaticEnv()
	err := mapstructure.Decode(v.AllSettings(), &config)
	if err != nil {
		panic(fmt.Errorf("error decoding config: %s", err))
	}
	if config.Address == "" || config.PollInterval == 0 || config.ReportInterval == 0 {
		// load
		v.SetConfigName("app")
		v.SetConfigType("env")
		v.AddConfigPath(file)
		err = v.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("error reading config: %s", err))
		}
		// Retrieve the value using the key
		p := v.Get("POLL_INTERVAL")
		r := v.Get("REPORT_INTERVAL")
		// Remove the "s" for autotests
		if strings.HasSuffix(p.(string), "s") {
			p = strings.TrimSuffix(p.(string), "s")
		}
		if strings.HasSuffix(r.(string), "s") {
			r = strings.TrimSuffix(r.(string), "s")
		}
		// Convert to int64
		p, err = strconv.ParseInt(p.(string), 10, 64)
		if err != nil {
			panic(fmt.Errorf("error decoding config: %s", err))
		}
		r, err = strconv.ParseInt(r.(string), 10, 64)
		if err != nil {
			panic(fmt.Errorf("error decoding config: %s", err))
		}
		// Set the value
		v.Set("POLL_INTERVAL", p.(int64))
		v.Set("REPORT_INTERVAL", r.(int64))
		err = v.Unmarshal(&config)
		if err != nil {
			panic(fmt.Errorf("error decoding config: %s", err))
		}
	}
}
