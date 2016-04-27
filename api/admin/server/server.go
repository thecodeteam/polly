package server

import (
	"net/http"

	ctypes "github.com/emccode/polly/core/types"
	"github.com/emccode/polly/core/volumes"
	"github.com/gorilla/mux"
)

// Router holds the router and Polly core object
type Router struct {
	r   *mux.Router
	p   *ctypes.Polly
	vsc *volumes.Vsc
}

// Start creates a new router with a nested Polly Core object
func Start(p *ctypes.Polly) *Router {

	r := &Router{
		r:   mux.NewRouter(),
		p:   p,
		vsc: volumes.New(p),
	}

	//volumes
	r.r.HandleFunc("/admin/volumes", r.getVolumesHandler).Methods("GET")
	r.r.HandleFunc("/admin/volumesall", r.getVolumesAllHandler).Methods("GET")
	r.r.HandleFunc("/admin/volumes/{volumeID}", r.getVolumeInspectHandler).Methods("GET")
	r.r.HandleFunc("/admin/volumes", r.postVolumesHandler).Methods("POST")
	r.r.HandleFunc("/admin/volumeoffer", r.postVolumeOfferHandler).Methods("POST")
	r.r.HandleFunc("/admin/volumeofferrevoke", r.postVolumeOfferRevokeHandler).Methods("POST")
	r.r.HandleFunc("/admin/volumelabel", r.postVolumeLabelHandler).Methods("POST")
	r.r.HandleFunc("/admin/volumelabelsremove", r.postVolumeLabelsRemoveHandler).Methods("POST")
	r.r.HandleFunc("/admin/volumes/{volid}", r.deleteVolumesHandler).Methods("DELETE")

	http.Handle("/", r.r)

	go func() {
		err := http.ListenAndServe(p.Config.GetString("polly.host"), nil)
		panic(err)
	}()

	return r
}
