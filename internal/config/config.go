package config

import (
	"fmt"
	"flag"
)


type Config struct {
	FilePath string
}


func FromFlagArgs() (*Config, error) {
	conf := &Config{}
	flag.StringVar(&conf.FilePath, "input", "", "input file path")
	flag.StringVar(&conf.FilePath, "i", "", "input file path")
	flag.Parse()

	if conf.FilePath == "" {
		return nil, fmt.Errorf("input flag is required use -input \"input file path\"")
	}
	return conf, nil
}
