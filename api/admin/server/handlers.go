package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/akutz/goof"
	"github.com/emccode/polly/api/types"
	"github.com/emccode/polly/core/version"
	"github.com/gorilla/mux"
)

func (rtr *Router) getVolumesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Debug("getVolumesHandler")
	vols, err := rtr.vsc.Volumes(r.URL.Query())
	if err != nil {
		http.Error(w, goof.WithError("problem getting volumes", err).Error(),
			http.StatusInternalServerError)
		return
	}

	j, _ := json.Marshal(&vols)
	w.Write(j)
}

func (rtr *Router) getVolumesAllHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Debug("getVolumesAllHandler")
	vols, err := rtr.vsc.VolumesAll(r.URL.Query())
	if err != nil {
		http.Error(w, goof.WithError("problem getting volumes", err).Error(),
			http.StatusInternalServerError)
		return
	}

	j, _ := json.Marshal(&vols)
	w.Write(j)
}

func (rtr *Router) getVolumeInspectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Debug("getVolumeInspectHandler")
	volumeID := mux.Vars(r)["volumeID"]
	vol, err := rtr.vsc.VolumeInspect(volumeID)
	if err != nil {
		http.Error(w, goof.WithError("problem getting volumes", err).Error(),
			http.StatusInternalServerError)
		return
	}

	j, _ := json.Marshal(&vol)
	w.Write(j)
}

func (rtr *Router) postVolumeOfferHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var o *types.VolumeOfferRequest
	b, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(b, &o)
	if err != nil {
		http.Error(w, "json is unparsable", http.StatusBadRequest)
		return
	}

	// Process mandatory elements on create
	if o.VolumeID == "" {
		http.Error(w, "mandatory volumeID missing or empty", 422)
		return
	}

	vol, err := rtr.vsc.VolumeOffer(o.VolumeID, o.Schedulers)
	if err != nil {
		http.Error(w, "problem performing volume offer", 422)
	}

	j, _ := json.Marshal(vol)
	w.Write(j)
}

func (rtr *Router) postVolumeOfferRevokeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var o *types.VolumeOfferRevokeRequest
	b, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(b, &o)
	if err != nil {
		http.Error(w, "json is unparsable", http.StatusBadRequest)
		return
	}

	// Process mandatory elements on create
	if o.VolumeID == "" {
		http.Error(w, "mandatory volumeID missing or empty", 422)
		return
	}

	vol, err := rtr.vsc.VolumeOfferRevoke(o.VolumeID, o.Schedulers)
	if err != nil {
		http.Error(w, "problem performing volume offer", 422)
	}

	j, _ := json.Marshal(vol)
	w.Write(j)
}

func (rtr *Router) postVolumeLabelHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var o *types.VolumeLabelRequest
	b, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(b, &o)
	if err != nil {
		http.Error(w, "json is unparsable", http.StatusBadRequest)
		return
	}

	// Process mandatory elements on create
	if o.VolumeID == "" {
		http.Error(w, "mandatory volumeID missing or empty", 422)
		return
	}

	vol, err := rtr.vsc.VolumeLabel(o.VolumeID, o.Labels)
	if err != nil {
		http.Error(w, "problem performing volume offer", 422)
	}

	j, _ := json.Marshal(vol)
	w.Write(j)
}

func (rtr *Router) postVolumeLabelsRemoveHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var o *types.VolumeLabelsRemoveRequest
	b, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(b, &o)
	if err != nil {
		http.Error(w, "json is unparsable", http.StatusBadRequest)
		return
	}

	// Process mandatory elements on create
	if o.VolumeID == "" {
		http.Error(w, "mandatory volumeID missing or empty", 422)
		return
	}

	vol, err := rtr.vsc.VolumeLabelsRemove(o.VolumeID, o.Labels)
	if err != nil {
		http.Error(w, "problem performing volume offer", 422)
	}

	j, _ := json.Marshal(vol)
	w.Write(j)
}

func (rtr *Router) postVolumesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var m *types.VolumeCreateRequest
	b, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(b, &m)
	if err != nil {
		http.Error(w, "json is unparsable", http.StatusBadRequest)
		return
	}

	// Process mandatory elements on create
	if m.ServiceName == "" {
		http.Error(w, "mandatory ServiceName missing or empty", 422)
		return
	}

	if _, ok := rtr.p.LsClient.Services[strings.ToLower(m.ServiceName)]; !ok {
		http.Error(w, "ServiceName is not defined", http.StatusNotFound)
		return
	}

	if m.Name == "" {
		http.Error(w, "mandatory Name missing or empty", 422)
		return
	}

	volNew, err := rtr.vsc.VolumeCreate(m)
	if err != nil {
		log.WithError(err).Error("volume creation failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	j, _ := json.Marshal(volNew)

	w.Write(j)
}

func (rtr *Router) deleteVolumesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ok bool
	var volid string
	vars := mux.Vars(r)

	if volid, ok = vars["volid"]; !ok {
		http.Error(w, "volid missing", http.StatusBadRequest)
		return
	}

	vs := strings.SplitN(volid, "-", 2)
	if len(vs) != 2 {
		http.Error(w, "valid volid is service-volumeid", http.StatusBadRequest)
		return
	}

	err := rtr.vsc.VolumeRemove(volid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	return
}

// getVersionHandler is gorilla mux handler for GET version on REST API
func (rtr *Router) getVersionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ver := &types.VersionResponse{
		VersionPollyAdminAPI:     version.VersionStr,
		VersionPollySchedulerAPI: version.VersionStr,
		VersionPollyBuild:        version.SemVer,
	}

	j, _ := json.Marshal(&ver)
	w.Write(j)
}

// notAllowedHandler is used to set a status code for unsupported operations
func (rtr *Router) notAllowedHandler(allow ...string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		mesg := "Method Not Allowed. "
		if len(allow) > 0 {
			mesg += "Allow " + strings.Join(allow, ", ")
		}
		http.Error(w, mesg, http.StatusMethodNotAllowed)
		return
	}
}

// notImplementedHandler is used to set a status code for future operations
func (rtr *Router) notImplementedHandler() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "Not Implemented Yet", http.StatusNotImplemented)
		return
	}
}
