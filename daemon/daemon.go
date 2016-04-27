package daemon

import (
	log "github.com/Sirupsen/logrus"
	gofig "github.com/akutz/gofig"
	"github.com/akutz/goof"
	core "github.com/emccode/polly/core"
	"os"
	"strconv"
)

// Run the Polly daemon
func Run(cfg gofig.Config) error {
	p := core.NewWithConfig(cfg)

	if err := core.Run(p); err != nil {
		return goof.New("problem starting polly core")
	}
	return nil
}

// Start the Polly daemon
func Start(cfg gofig.Config, init chan error, stop <-chan os.Signal) error {
	p := core.NewWithConfig(cfg)

	if debug, _ := strconv.ParseBool(os.Getenv("POLLY_DEBUG")); debug {
		log.SetLevel(log.DebugLevel)
		os.Setenv("POLLY_SERVER_HTTP_LOGGING_ENABLED", "true")
		os.Setenv("POLLY_SERVER_HTTP_LOGGING_LOGREQUEST", "true")
		os.Setenv("POLLY_SERVER_HTTP_LOGGING_LOGRESPONSE", "true")
		os.Setenv("POLLY_CLIENT_HTTP_LOGGING_ENABLED", "true")
		os.Setenv("POLLY_CLIENT_HTTP_LOGGING_LOGREQUEST", "true")
		os.Setenv("POLLY_CLIENT_HTTP_LOGGING_LOGRESPONSE", "true")
	}

	var err error
	if err = core.Start(p); err != nil {
		init <- goof.WithError("problem starting polly core", err)
	}

	if init != nil {
		close(init)
	}

	if err != nil {
		log.Error(err)
		log.WithError(err).Error("service initialization failed")
	}

	log.Info("service successfully initialized, waiting on stop signal")

	if stop != nil {
		<-stop
		log.Info("Service received stop signal")
	}
	return nil
}
