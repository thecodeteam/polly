package main

import (
	"bytes"
	"os"

	"testing"

	log "github.com/Sirupsen/logrus"
	gofig "github.com/akutz/gofig"
	core "github.com/emccode/polly/core"
	"github.com/stretchr/testify/assert"
)

func TestPolly(t *testing.T) {
	cfg := gofig.New()

	yml := []byte(defaultConfig)
	if err := cfg.ReadConfig(bytes.NewReader(yml)); err != nil {
		log.WithError(err).Fatal("problem reading config")
		os.Exit(1)
	}

	pc, err := core.NewWithConfig(cfg)
	if err != nil {
		log.WithError(err).Fatal("problem creating new polly core")
		os.Exit(1)
	}

	assert.NotEqual(t, pc, nil)
	assert.NotEqual(t, pc.Store, nil)
	assert.Equal(t, pc.Store.StoreType(), "boltdb")

	vols, err := pc.LsClient.Volumes()
	if err != nil {
		log.WithError(err).Fatal("problem getting volumes from libstorage")
		os.Exit(1)
	}

	assert.Len(t, vols, 3)
	assert.Equal(t, vols[0])

}
