package configs

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Address        string `mapstructure:"ADDRESS"`
	Port           string `mapstructure:"PORT"`
	PollInterval   int    `mapstructure:"POLL_INTERVAL"`
	ReportInterval int    `mapstructure:"REPORT_INTERVAL"`
}

var CFG = &Config{}

func (config *Config) LoadConfig(path string) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(err)
	}
}
