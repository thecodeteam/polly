package core

import (
	"log"

	"github.com/akutz/gofig"
	"github.com/akutz/goof"
	adminserver "github.com/emccode/polly/api/admin/server"
	"github.com/emccode/polly/core/libstorage/client"
	"github.com/emccode/polly/core/libstorage/server"
	store "github.com/emccode/polly/core/store"
	ctypes "github.com/emccode/polly/core/types"
	"github.com/emccode/polly/core/volumes"
)

//NewWithConfigFile init the lib
func NewWithConfigFile(path string) (p *ctypes.Polly, err error) {
	config := gofig.New()

	myErr := config.ReadConfigFile(path)
	if myErr != nil {
		log.Fatal("Fatal: ", myErr)
		return nil, myErr
	}

	pollyCore, myErr := NewWithConfig(config)
	if myErr != nil {
		log.Fatal("Fatal: ", myErr)
		return nil, myErr
	}

	return pollyCore, nil
}

const (
	storeConfigBolt = `
polly:
  store:
    type: boltdb
    endpoints: /tmp/boltdb
    bucket: MyBoltDb
`
	storeConfigConsul = `
polly:
  store:
    type: consul
    endpoints: 127.0.0.1:8500
`
)

//NewWithConfig This initializes new instance of this library
func NewWithConfig(config gofig.Config) (*ctypes.Polly, error) {
	ps, err := store.NewWithConfig(config.Scope("polly.store"))
	if err != nil {
		return nil, err
	}

	lscfg, err := server.New()
	if err != nil {
		return nil, err
	}

	lsc, err := client.NewWithConfig(lscfg)
	if err != nil {
		return nil, err
	}

	services, err := lsc.Services()
	if err != nil {
		return nil, goof.WithError("cannot instantiate client services", err)
	}

	pc := &ctypes.Polly{
		Store:    ps,
		LsConfig: lscfg,
		LsClient: lsc,
		Config:   config,
		Services: services,
	}

	_ = adminserver.New(pc)
	if err := volumes.Init(pc); err != nil {
		return nil, err
	}

	return pc, nil
}

// if err := pcfg.ReadConfigFile("/etc/polly/config.yml"); err != nil {
// 	return nil, err
// }
