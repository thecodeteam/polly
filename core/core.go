package core

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/akutz/gofig"
	"github.com/akutz/goof"
	"github.com/akutz/gotil"
	adminserver "github.com/emccode/polly/api/admin/server"
	"github.com/emccode/polly/core/libstorage/client"
	"github.com/emccode/polly/core/libstorage/server"
	store "github.com/emccode/polly/core/store"
	ctypes "github.com/emccode/polly/core/types"
	util "github.com/emccode/polly/util"
	"os"
	"strconv"
)

func init() {
	gofig.SetGlobalConfigPath(util.EtcDirPath())
	gofig.SetUserConfigPath(fmt.Sprintf("%s/.polly", gotil.HomeDir()))
	gofig.Register(globalRegistration())

	if debug, _ := strconv.ParseBool(os.Getenv("POLLY_DEBUG")); debug {
		log.SetLevel(log.DebugLevel)
		os.Setenv("POLLY_SERVER_HTTP_LOGGING_ENABLED", "true")
		os.Setenv("POLLY_SERVER_HTTP_LOGGING_LOGREQUEST", "true")
		os.Setenv("POLLY_SERVER_HTTP_LOGGING_LOGRESPONSE", "true")
		os.Setenv("POLLY_CLIENT_HTTP_LOGGING_ENABLED", "true")
		os.Setenv("POLLY_CLIENT_HTTP_LOGGING_LOGREQUEST", "true")
		os.Setenv("POLLY_CLIENT_HTTP_LOGGING_LOGRESPONSE", "true")
	}

}

func globalRegistration() *gofig.Registration {
	r := gofig.NewRegistration("Global")
	r.Yaml(`
polly:
  host: :7980
  logLevel: warn
`)
	r.Key(gofig.String, "l", "warn",
		"The log level (error, warn, info, debug)", "polly.logLevel",
		"logLevel")
	return r
}

//NewWithConfigFile init the lib
func NewWithConfigFile(path string) (*ctypes.Polly, error) {
	config := gofig.New()
	if err := config.ReadConfigFile(path); err != nil {
		return nil, goof.WithError("problem reading config", err)
	}

	p := NewWithConfig(config)

	return p, nil
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

// NewWithConfig initializes new polly object
func NewWithConfig(config gofig.Config) *ctypes.Polly {
	return &ctypes.Polly{
		Config: config,
	}
}

// Start starts the Polly core services and returns
func Start(p *ctypes.Polly) error {
	scfg, _ := p.Config.Copy()
	ps, err := store.NewWithConfig(scfg.Scope("polly.store"))
	if err != nil {
		return err
	}
	p.Store = ps

	lcfg, _ := p.Config.Copy()
	lscfg, err := server.New(lcfg.Scope("polly"))
	if err != nil {
		return err
	}
	p.LsConfig = lscfg

	lsc, err := client.NewWithConfig(lscfg)
	if err != nil {
		return err
	}
	p.LsClient = lsc

	services, err := lsc.Services()
	if err != nil {
		return goof.WithError("cannot instantiate client services", err)
	}
	p.Services = services

	_ = adminserver.Start(p)
	return nil
}

// Run starts the Polly core services and blocks
func Run(p *ctypes.Polly) error {
	if err := Start(p); err != nil {
		goof.WithError("could not run polly core services", err)
	}

	select {}
}
