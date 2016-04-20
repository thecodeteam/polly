package main

import (
	"github.com/akutz/goof"
	apihttp "github.com/emccode/libstorage/api/types/http"
	"github.com/emccode/libstorage/client"
)

var services apihttp.ServicesMap
var lsClient client.Client

func startClient() error {
	var err error
	lsClient, err = getClient(config)
	if err != nil {
		return goof.WithError("cannot connect to libstorage client", err)
	}

	services, err = lsClient.Services()
	if err != nil {
		return goof.WithError("cannot instantiate client services", err)
	}
	return nil
}
