package configs

import (
	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog/log"
	"os"
)

// ReadFile reads the specified configuration, file if it exists, overwriting what is already
// in the config. The file format is TOML
func (c *AppConf) ReadFile(f string) error {

	// Does the config file exist?
	_, err := os.Stat(f)
	if err != nil {
		return err
	}

	// Can we parse it?
	if _, err := toml.DecodeFile(f, c); err != nil {
		log.Error().Err(err).Msg("Error decoding config file")
		return err
	}

	return nil
}

