package server

import (
	"net/http"

	ctypes "github.com/emccode/polly/core/types"

	"github.com/gorilla/mux"
)

// Router holds the router and Polly core object
type Router struct {
	r *mux.Router
	p *ctypes.Polly
}

// New creates a new router with a nested Polly Core object
func New(p *ctypes.Polly) *Router {
	r := &Router{
		r: mux.NewRouter(),
		p: p,
	}

	//volumes
	r.r.HandleFunc("/admin/volumes", r.getVolumesHandler).Methods("GET")
	r.r.HandleFunc("/admin/volumes", r.postVolumesHandler).Methods("POST")
	r.r.HandleFunc("/admin/volumes/{volid}", r.deleteVolumesHandler).Methods("DELETE")

	http.Handle("/", r.r)
	go http.ListenAndServe(":8080", nil)

	return r
}
