package polly

import (
	"bytes"
	"os"

	"testing"

	gofig "github.com/akutz/gofig"
	core "github.com/emccode/polly/core"
	"github.com/stretchr/testify/assert"
)

const (
	DefaultConfigFilename = "/etc/polly/config.yml"

	pollyConfigBaseBolt = `
polly:
  store:
    type: boltdb
    endpoints: /tmp/boltdb
    bucket: MyBoltDb
`
)

func TestPolly(t *testing.T) {
	config := gofig.New()

	loaded := false
	if _, err := os.Stat(DefaultConfigFilename); os.IsNotExist(err) {
		if err := config.ReadConfigFile(DefaultConfigFilename); err == nil {
			loaded = true
		}
	}

	if !loaded {
		configYamlBuf := []byte(pollyConfigBaseBolt)
		if err := config.ReadConfig(bytes.NewReader(configYamlBuf)); err != nil {
			panic(err)
		}
	}

	ps, err := core.NewWithConfig(config.Scope("polly.store"))
	if err != nil {
		panic(err)
	}

	assert.NotEqual(t, ps, nil)
	assert.NotEqual(t, ps.PollyStore, nil)
	assert.Equal(t, ps.PollyStore.StoreType(), "boltdb")

}
