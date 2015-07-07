package cachectl

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

type ConfToml struct {
	Targets []SectionTarget `toml:"targets"`
}

type SectionTarget struct {
	Path          string  `toml:"path"`
	PurgeInterval int     `toml:"purge_interval"`
	Filter        string  `toml:"filter"`
	Rate          float64 `toml:"rate"`
}

func ValidateConf(confToml *ConfToml) error {
	for i, target := range confToml.Targets {
		_, err := os.Stat(target.Path)
		if err != nil {
			return err
		}
		if target.Filter == "" || target.Filter == "*" {
			confToml.Targets[i].Filter = ".*"
		}
		if target.Rate < 0 || target.Rate > 1.0 {
			log.Printf("[warning] target: %s, rate is invalid: %f. zero will be assigned.", target.Path, target.Rate)
			confToml.Targets[i].Rate = 0
		}
		if target.PurgeInterval == 0 {
			log.Printf("[warning] target: %s, purge_interval is invalid: %d, or not set. 3600 will be assigned.", target.Path, target.PurgeInterval)
			confToml.Targets[i].PurgeInterval = 3600
		}
	}
	return nil
}

func LoadConf(confPath string, confToml *ConfToml) error {
	_, err := toml.DecodeFile(confPath, confToml)
	if err != nil {
		return err
	}
	return nil
}
