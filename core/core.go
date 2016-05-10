package core

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/akutz/gofig"
	"github.com/akutz/goof"
	"github.com/akutz/gotil"
	"github.com/emccode/libstorage/api/context"
	apivolroute "github.com/emccode/libstorage/api/server/router/volume"
	apitypes "github.com/emccode/libstorage/api/types"
	adminserver "github.com/emccode/polly/api/admin/server"
	catypes "github.com/emccode/polly/api/types"
	"github.com/emccode/polly/core/libstorage/client"
	"github.com/emccode/polly/core/libstorage/server"
	store "github.com/emccode/polly/core/store"
	ctypes "github.com/emccode/polly/core/types"
	util "github.com/emccode/polly/util"
	"net/http"
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

	if debug, _ := strconv.ParseBool(os.Getenv("LIBSTORAGE_DEBUG")); debug {
		log.SetLevel(log.DebugLevel)
		os.Setenv("LIBSTORAGE_SERVER_HTTP_LOGGING_ENABLED", "true")
		os.Setenv("LIBSTORAGE_SERVER_HTTP_LOGGING_LOGREQUEST", "true")
		os.Setenv("LIBSTORAGE_SERVER_HTTP_LOGGING_LOGRESPONSE", "true")
		os.Setenv("LIBSTORAGE_CLIENT_HTTP_LOGGING_ENABLED", "true")
		os.Setenv("LIBSTORAGE_CLIENT_HTTP_LOGGING_LOGREQUEST", "true")
		os.Setenv("LIBSTORAGE_CLIENT_HTTP_LOGGING_LOGRESPONSE", "true")
	}

}

func globalRegistration() *gofig.Registration {
	r := gofig.NewRegistration("Global")
	r.Yaml(`
polly:
  host: :7978
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

	ctx := context.Background()
	if lscfg.GetString("libstorage.client.requestPath") == "" {
		lscfg.Set("libstorage.client.requestPath", "admin")
	}

	lsc, err := client.NewWithConfig(ctx, lscfg)
	if err != nil {
		return err
	}
	p.LsClient = lsc

	// filterVolume sets a filter in place for requests
	filterVolume := func(
		ctx apitypes.Context,
		req *http.Request,
		store apitypes.Store,
		volume *apitypes.Volume) (bool, error) {

		// set the polly volume settings
		volumeNew, err := client.NewVolume(lsc, volume, ctx.ServiceName())
		if err != nil {
			return false, err
		}

		volume.Fields = make(map[string]string)
		volume.Fields["polly.id"] = volumeNew.VolumeID

		// apply filtering
		rp := req.Header.Get("requestpath")
		ctx.WithField("requestPath", rp).Info("volume response on request path")

		if ctx.Route() == "volumeCreate" {
			// establish new volume metadata for new libstorage inbound requests
			ctx.WithField("route", ctx.Route()).Debug("volumes create route")
			volumeNew.Schedulers = []string{ctx.ServiceName()}
			err = p.Store.SaveVolumeMetadata(volumeNew)
			if err != nil {
				return false, goof.WithError("failed to save metadata", err)
			}

			return updateVolume(ctx, p, volumeNew, volume, true, rp)
		} else if rp == "admin" {
			ctx.WithField("requestPath", rp).Debug("volumes from admin request")
			return updateVolume(ctx, p, volumeNew, volume, false, rp)
		}
		ctx.WithField("requestPath", rp).Debug("volumes from non-admin request")

		return updateVolume(ctx, p, volumeNew, volume, true, rp)
	}

	apivolroute.OnVolume = filterVolume

	_ = adminserver.Start(p)
	return nil
}

func updateVolume(ctx apitypes.Context, p *ctypes.Polly,
	volumeNew *catypes.Volume, volume *apitypes.Volume,
	mustExist bool, rp string) (bool, error) {

	if exists, err := p.Store.Exists(volumeNew); err != nil {
		return false, err
	} else if mustExist && !exists {
		return false, err
	} else if exists {
		if err := p.Store.SetVolumeAdminLabels(volumeNew); err != nil {
			return false, err
		}

		if _, err := p.Store.SetVolumeMetadata(volumeNew); err != nil {
			return false, err
		}

		for k, v := range volumeNew.Labels {
			volume.Fields[fmt.Sprintf("polly.labels.%s", k)] = v
		}
	}

	if !util.ContainsString(volumeNew.Schedulers, ctx.ServiceName()) &&
		rp != "admin" {
		return false, nil
	}

	return true, nil
}

// Run starts the Polly core services and blocks
func Run(p *ctypes.Polly) error {
	if err := Start(p); err != nil {
		goof.WithError("could not run polly core services", err)
	}

	select {}
}
