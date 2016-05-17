package server

import (
	"net/http"

	"github.com/akutz/gotil"
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
	r.r.HandleFunc("/admin/version", r.getVersionHandler).Methods("GET")
	r.r.HandleFunc("/admin/version",
		r.notAllowedHandler("GET")).Methods("POST", "PUT", "PATCH", "DELETE")
	r.r.HandleFunc("/admin/volumes", r.getVolumesHandler).Methods("GET")
	r.r.HandleFunc("/admin/volumes", r.postVolumesHandler).Methods("POST")
	r.r.HandleFunc("/admin/volumes",
		r.notAllowedHandler("GET", "POST")).Methods("PUT", "PATCH", "DELETE")
	r.r.HandleFunc("/admin/volumesall", r.getVolumesAllHandler).Methods("GET")
	r.r.HandleFunc("/admin/volumesall",
		r.notAllowedHandler("GET")).Methods("POST", "PUT", "PATCH", "DELETE")
	r.r.HandleFunc("/admin/volumes/{volumeID}", r.getVolumeInspectHandler).Methods("GET")
	r.r.HandleFunc("/admin/volumes/{volid}", r.deleteVolumesHandler).Methods("DELETE")
	r.r.HandleFunc("/admin/volumes/{volumeID}",
		r.notAllowedHandler("GET", "DELETE")).Methods("PUT", "PATCH", "POST")
	r.r.HandleFunc("/admin/volumeoffer", r.postVolumeOfferHandler).Methods("POST")
	r.r.HandleFunc("/admin/volumeoffer",
		r.notAllowedHandler("POST")).Methods("GET", "PUT", "PATCH", "DELETE")
	r.r.HandleFunc("/admin/volumeofferrevoke", r.postVolumeOfferRevokeHandler).Methods("POST")
	r.r.HandleFunc("/admin/volumeofferrevoke",
		r.notAllowedHandler("POST")).Methods("GET", "PUT", "PATCH", "DELETE")
	r.r.HandleFunc("/admin/volumelabel", r.postVolumeLabelHandler).Methods("POST")
	r.r.HandleFunc("/admin/volumelabel",
		r.notAllowedHandler("POST")).Methods("GET", "PUT", "PATCH", "DELETE")
	r.r.HandleFunc("/admin/volumelabelsremove", r.postVolumeLabelsRemoveHandler).Methods("POST")
	r.r.HandleFunc("/admin/volumelabelsremove",
		r.notAllowedHandler("POST")).Methods("GET", "PUT", "PATCH", "DELETE")

	http.Handle("/", r.r)

	_, lAddr, err := gotil.ParseAddress(p.Config.GetString("polly.host"))
	if err != nil {
		panic(err)
	}

	go func() {
		err := http.ListenAndServe(lAddr, nil)
		panic(err)
	}()

	return r
}
