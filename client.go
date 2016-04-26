package main

import (
	"github.com/akutz/gofig"
	"github.com/akutz/goof"

	"github.com/emccode/libstorage/client"
)

func getClient(config gofig.Config) (client.Client, error) {
	client, err := client.New(config)
	if err != nil {
		return nil, goof.WithFieldE(
			"host", config.Get("libstorage.host"),
			"error dialing libStorage service", err)
	}
	return client, nil
}
