package config

import (
	"bytes"

	gofig "github.com/akutz/gofig"
	goof "github.com/akutz/goof"
)

// DefaultConfig can be used globally
var DefaultConfig = `
polly:
  host: tcp://127.0.0.1:7978
  store:
    type: boltdb
    endpoints: /tmp/boltdb
    bucket: MyBoltDb_test
  libstorage:
    host: tcp://127.0.0.1:7981
    embedded: false
    server:
      endpoints:
        localhost:
          address: tcp://:7981
      services:
        mock:
          libstorage:
            driver: mock
`

// New returns a new configuration object
func New() (gofig.Config, error) {
	cfg := gofig.New()

	if !cfg.IsSet("polly.store") && !cfg.IsSet("polly.libstorage") {
		yml := []byte(DefaultConfig)
		if err := cfg.ReadConfig(bytes.NewReader(yml)); err != nil {
			return nil, goof.WithError("problem reading config", err)
		}
	}

	return cfg, nil
}

// NewWithConfig returna a new configuration object from a yaml string
func NewWithConfig(yamlConfig string) (gofig.Config, error) {
	cfg := gofig.New()
	yml := []byte(yamlConfig)
	if err := cfg.ReadConfig(bytes.NewReader(yml)); err != nil {
		return nil, goof.WithError("problem reading config", err)
	}

	return cfg, nil
}
