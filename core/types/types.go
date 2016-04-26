package types

import (
	"github.com/akutz/gofig"
	"github.com/emccode/libstorage/api/types/http"
	lsclient "github.com/emccode/polly/core/libstorage/client"
	store "github.com/emccode/polly/core/store"
)

//Polly this represents the "core" functionality for Polly
type Polly struct {
	Store    *store.PollyStore
	LsClient *lsclient.Client
	Config   gofig.Config
	LsConfig gofig.Config
	Services http.ServicesMap
	// server.Server
}
