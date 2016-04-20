package main

import (
	"bytes"
	"log"

	"github.com/akutz/gofig"
	"github.com/akutz/goof"
	apihttp "github.com/emccode/libstorage/api/types/http"
	"github.com/emccode/libstorage/client"
	store "github.com/emccode/polly/pkg/store"
)

// todo remove this persistent store stuff when core has it
const (
	libStorageConfigBaseBolt = `
polly:
  store:
    type: boltdb
    endpoints: /tmp/boltdb
    bucket: MyBoltDb
`
	libStorageConfigBaseConsul = `
polly:
  store:
    type: consul
    endpoints: 127.0.0.1:8500
`
)

var services apihttp.ServicesMap
var lsClient client.Client
var ps *store.PollyStore

//var ps

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

	// todo remove this persistent store stuff when core has it
	config := gofig.New()

	configYamlBuf := []byte(libStorageConfigBaseBolt)
	if err = config.ReadConfig(bytes.NewReader(configYamlBuf)); err != nil {
		panic(err)
	}

	ps, err = store.NewWithConfig(config.Scope("polly.store"))
	if err != nil {
		log.Fatal("Failed to create PollyStore")
	}

	return nil
}
