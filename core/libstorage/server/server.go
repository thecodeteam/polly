package server

import (
	"github.com/akutz/gofig"

	"github.com/emccode/libstorage/api/server"
)

// New starts a server with default configuration
func New(config gofig.Config) (gofig.Config, error) {
	return config, NewWithConfig(config.Scope("polly"))
}

// NewWithConfig starts a server with a configuration
func NewWithConfig(config gofig.Config) error {
	_, errs, err := server.Serve(nil, config)
	if err != nil {
		return err
	}
	go func() {
		err := <-errs
		panic(err)
	}()

	return nil
}
