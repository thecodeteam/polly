package server

import (
	"github.com/akutz/gofig"

	"github.com/emccode/libstorage/api/server"
)

// New starts a server with default configuration
func New(config gofig.Config) (gofig.Config, error) {
	if config != nil {
		return config, NewWithConfig(config.Scope("polly"))
	}

	cfg, _, err, errs := server.Start("", false, "mock", "mock")
	if err != nil {
		return nil, err
	}
	go func() {
		err := <-errs
		panic(err)
	}()
	return cfg, nil

}

// NewWithConfig starts a server with a configuration
func NewWithConfig(config gofig.Config) error {
	_, err, errs := server.StartWithConfig(config)
	if err != nil {
		return err
	}
	go func() {
		err := <-errs
		panic(err)
	}()

	return nil
}
