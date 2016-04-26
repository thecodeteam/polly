package main

import (
	"bytes"
	"os"

	log "github.com/Sirupsen/logrus"
	gofig "github.com/akutz/gofig"
	core "github.com/emccode/polly/core"
)

var defaultConfig = `
polly:
  store:
    type: boltdb
    endpoints: /tmp/boltdb
    bucket: MyBoltDb
  libstorage:
    host: tcp://localhost:7979
    profiles:
      enabled: true
      groups:
      - local=127.0.0.1
    server:
      endpoints:
        localhost:
          address: tcp://localhost:7979
      services:
        mock:
          libstorage:
            driver: mock
`

func main() {
	cfg := gofig.New()

	yml := []byte(defaultConfig)
	if err := cfg.ReadConfig(bytes.NewReader(yml)); err != nil {
		log.WithError(err).Fatal("problem reading config")
		os.Exit(1)
	}

	_, err := core.NewWithConfig(cfg)
	if err != nil {
		log.WithError(err).Fatal("problem creating new polly core")
		os.Exit(1)
	}

	select {}
}
