package polly

import (
	"github.com/akutz/gofig"
	"github.com/akutz/gotil"

	_ "github.com/emccode/libstorage"
	_ "github.com/emccode/libstorage/imports/local"
	_ "github.com/emccode/libstorage/imports/remote"

	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"strconv"

	core "github.com/emccode/polly/core"
	ctypes "github.com/emccode/polly/core/types"
	util "github.com/emccode/polly/util"
)

// NewWithConfigFile creates a new Polly instance and configures it with a
// custom configuration file.
func NewWithConfigFile(path string) (*ctypes.Polly, error) {
	return core.NewWithConfigFile(path)
}

// NewWithConfig creates a new Polly instance and configures it with a
// custom configuration stream.
func NewWithConfig(config gofig.Config) *ctypes.Polly {
	return core.NewWithConfig(config)
}

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
  logLevel: warn
`)
	r.Key(gofig.String, "l", "warn",
		"The log level (error, warn, info, debug)", "polly.logLevel",
		"logLevel")
	r.Key(gofig.String, "", "tcp://127.0.0.1:7978", "", "polly.host")
	return r
}
