package helpers

import (
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Output string
	Target string
	Ignore Ignore
}
type Ignore struct {
	PackagePrefix string
	Files         []string
}

func LoadConfig(configFile string) *Config {
	file, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	newConfig := new(Config)

	err = toml.Unmarshal(file, newConfig)
	if err != nil {
		log.Fatal(err)
	}

	return newConfig
}
