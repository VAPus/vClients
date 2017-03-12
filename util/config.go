package util

import "github.com/BurntSushi/toml"

type Config struct {
	Mode string
	Port int
	Host string
	SSL configSSL
}

type configSSL struct {
	Enabled bool
	Key string
	Cert string
}

// LoadConfigurationFile loads a TOML file into the given struct
func LoadConfigurationFile(path string) (cfg *Config, err error) {
	cfg = &Config{}

	if _, err = toml.DecodeFile(path, cfg); err != nil {
		return
	}

	return
}