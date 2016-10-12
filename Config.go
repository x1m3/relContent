package main

import (
	"github.com/BurntSushi/toml"
	"fmt"
)

type Config struct {
	RunConsolidatorEverySeconds int `toml:"run_consolidator_every"`
	SessionClosedAfterSeconds int `toml:"session_closed_after"`
	MinSessionSize int `toml:"min_session_size"`
	DatabaseFilename string `toml:"database_filename"`
	NonceSecret string `toml:"nonce_secret"`
}

func (conf *Config) load(configFileName string) bool {

	conf.defaultValues()

	if _, err := toml.DecodeFile(configFileName, conf); err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println(conf)

	return true
}

func (conf *Config) defaultValues() {

}