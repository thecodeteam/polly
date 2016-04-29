package daemon

import (
	"bytes"
	"os"

	"testing"

	log "github.com/Sirupsen/logrus"
	gofig "github.com/akutz/gofig"
	core "github.com/emccode/polly/core"
	config "github.com/emccode/polly/core/config"
	"github.com/stretchr/testify/assert"
)

func TestPolly(t *testing.T) {
	cfg := gofig.New()

	yml := []byte(config.DefaultConfig)
	if err := cfg.ReadConfig(bytes.NewReader(yml)); err != nil {
		log.WithError(err).Fatal("problem reading config")
		os.Exit(1)
	}

	p := core.NewWithConfig(cfg)
	if err := core.Start(p); err != nil {
		log.WithError(err).Fatal("problem starting polly core")
		os.Exit(1)
	}

	assert.NotEqual(t, p, nil)
	assert.NotEqual(t, p.Store, nil)
	assert.Equal(t, p.Store.StoreType(), "boltdb")

	vols, err := p.LsClient.Volumes()
	if err != nil {
		log.WithError(err).Fatal("problem getting volumes from libstorage")
		os.Exit(1)
	}

	assert.Len(t, vols, 3)
	// assert.Equal(t, vols[0])

}
