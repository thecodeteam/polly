package client

import (
	"crypto/tls"
	"net"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/akutz/gofig"
	"github.com/akutz/goof"
	"github.com/akutz/gotil"

	"github.com/emccode/libstorage/api/utils"
	apiclient "github.com/emccode/polly/api/admin/client"

	// load the drivers
	_ "github.com/emccode/libstorage/drivers/os"
)

func init() {
	registerConfig()
}

const (
	clientScope          = "polly.client"
	hostKey              = "polly.host"
	logEnabledKey        = "polly.client.http.logging.enabled"
	logOutKey            = "polly.client.http.logging.out"
	logErrKey            = "polly.client.http.logging.err"
	logRequestsKey       = "polly.client.http.logging.logrequest"
	logResponsesKey      = "polly.client.http.logging.logresponse"
	disableKeepAlivesKey = "polly.client.http.disableKeepAlives"
)

type pc struct {
	apiclient.Client
}

// New returns a new Client.
func New(config gofig.Config) (Client, error) {

	logFields := log.Fields{}

	if config == nil {
		return nil, goof.New("missing configuration when configuring client")
	}

	addr := config.GetString(hostKey)

	proto, lAddr, err := gotil.ParseAddress("tcp://127.0.0.1" + addr)
	if err != nil {
		return nil, err
	}

	tlsConfig, err := utils.ParseTLSConfig(
		config.Scope(clientScope), logFields)
	if err != nil {
		return nil, err
	}

	return &pc{
		Client: apiclient.Client{
			Host:         getHost(proto, lAddr, tlsConfig),
			Headers:      http.Header{},
			LogRequests:  config.GetBool(logRequestsKey),
			LogResponses: config.GetBool(logResponsesKey),
			Client: &http.Client{
				Transport: &http.Transport{
					Dial: func(string, string) (net.Conn, error) {
						if tlsConfig == nil {
							return net.Dial(proto, lAddr)
						}
						return tls.Dial(proto, lAddr, tlsConfig)
					},
					DisableKeepAlives: config.GetBool(disableKeepAlivesKey),
				},
			},
		},
	}, nil
}

func getHost(proto, lAddr string, tlsConfig *tls.Config) string {
	if tlsConfig != nil && tlsConfig.ServerName != "" {
		return tlsConfig.ServerName
	} else if proto == "unix" {
		return "polly-server"
	} else {
		return lAddr
	}
}

func registerConfig() {
	r := gofig.NewRegistration("polly Client")
	r.Key(gofig.Bool, "", false, "", logEnabledKey)
	r.Key(gofig.String, "", "", "", logOutKey)
	r.Key(gofig.String, "", "", "", logErrKey)
	r.Key(gofig.Bool, "", false, "", logRequestsKey)
	r.Key(gofig.Bool, "", false, "", logResponsesKey)
	r.Key(gofig.Bool, "", false, "", disableKeepAlivesKey)
	gofig.Register(r)
}
