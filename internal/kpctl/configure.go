package kpctl

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kube-peering/internal/pkg/connectors"
	"github.com/pelletier/go-toml"
)

type Config struct {
	VPNConfig connectors.VPNConfig `toml:"vpn"`
}

func IsConfigExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (c *Config) WriteFile(path string) error {
	b, err := toml.Marshal(c)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	os.MkdirAll(dir, 0755)

	if err := ioutil.WriteFile(path, b, 0644); err != nil {
		return err
	}

	return nil
}

func ReadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return nil, nil
}
