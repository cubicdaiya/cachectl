package cachectl

import (
	"errors"
	"github.com/BurntSushi/toml"
	"os"
)

type ConfToml struct {
	Targets []SectionTarget `toml:"targets"`
}

type SectionTarget struct {
	Path          string  `toml:"path"`
	PurgeInterval uint64  `toml:"purge_interval"`
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
			return errors.New("invalid rate")
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
