package server

import (
	"github.com/akutz/gofig"

	"github.com/emccode/libstorage/api/server"
)

// New starts a server with default configuration
func New() (gofig.Config, error) {
	cfg, _, errs := server.Run("", false, "mock", "mock")
	go func() {
		err := <-errs
		panic(err)
	}()

	return cfg, nil
}

// NewWithConfig starts a server with a configuration
func NewWithConfig(config gofig.Config) (gofig.Config, error) {
	_, errs := server.RunWithConfig(config)
	go func() {
		err := <-errs
		panic(err)
	}()

	return config, nil
}
