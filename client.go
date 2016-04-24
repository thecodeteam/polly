package main

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/akutz/gofig"
	"github.com/akutz/goof"

	"github.com/emccode/libstorage/client"
	"github.com/emccode/polly/pkg/store"
)

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

func newClient(config gofig.Config) (client.Client, error) {
	client, err := client.New(config)
	if err != nil {
		return nil, goof.WithFieldE(
			"host", config.Get("libstorage.host"),
			"error dialing libStorage service", err)
	}
	return client, nil
}

func startClient(config gofig.Config) error {
	var err error
	fmt.Println(fmt.Sprintf("%+v", config))
	lsClient, err = newClient(config)
	if err != nil {
		return goof.WithError("cannot connect to libstorage client", err)
	}

	services, err = lsClient.Services()
	if err != nil {
		return goof.WithError("cannot instantiate client services", err)
	}

	// todo remove this persistent store stuff when core has it
	pcfg := gofig.New()

	configYamlBuf := []byte(libStorageConfigBaseBolt)
	if err = pcfg.ReadConfig(bytes.NewReader(configYamlBuf)); err != nil {
		panic(err)
	}

	ps, err = store.NewWithConfig(pcfg.Scope("polly.store"))
	if err != nil {
		log.Fatal("Failed to create PollyStore")
	}

	return nil
}
