package main

import "github.com/BurntSushi/toml"

type config struct {
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

func loadConfigurationFile(path string) (cfg *config, err error) {
	cfg = &config{}

	if _, err = toml.DecodeFile(path, cfg); err != nil {
		return
	}

	return
}