package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	apihttp "github.com/emccode/libstorage/api/types/http"

	log "github.com/Sirupsen/logrus"
	"github.com/emccode/polly/api/types"
	"github.com/gorilla/mux"
)

func applyVolumeFilter(v *types.Volume, r *http.Request) bool {
	for key, value := range r.URL.Query() {
		log.WithFields(log.Fields{
			"vol":   v,
			"lsvol": v.Volume,
			"key":   key,
			"value": value[0]}).Info("applyVolumeFilter")

		switch key {
		case "availabilityZone":
			if v.AvailabilityZone == value[0] {
				log.WithFields(log.Fields{
					"vol":   v,
					"lsvol": v.Volume,
					"key":   key,
					"value": value[0]}).Info("availabilityZone filter match")
				continue
			}
			log.WithFields(log.Fields{
				"vol":   v,
				"lsvol": v.Volume,
				"key":   key,
				"value": value[0]}).Info("rejected volume by AZ")
			return false
		case "iops":
			i, err := strconv.ParseInt(value[0], 10, 64)
			if err != nil {
				return false
			}
			if v.IOPS == i {
				log.WithFields(log.Fields{
					"vol":   v,
					"lsvol": v.Volume,
					"key":   key,
					"value": value[0]}).Info("IOPS filter match")
				continue
			}
			log.WithFields(log.Fields{
				"vol":   v,
				"lsvol": v.Volume,
				"key":   key,
				"value": value[0]}).Info("rejected volume by IOPS")
			return false
		case "size":
			i, err := strconv.ParseInt(value[0], 10, 64)
			if err != nil {
				return false
			}
			if v.Size == i {
				log.WithFields(log.Fields{
					"vol":   v,
					"lsvol": v.Volume,
					"key":   key,
					"value": value[0]}).Info("size filter match")
				continue
			}
			log.WithFields(log.Fields{
				"vol":   v,
				"lsvol": v.Volume,
				"key":   key,
				"value": value[0]}).Info("rejected volume by size")
			return false
		case "serviceName":
			if v.ServiceName == value[0] {
				log.WithFields(log.Fields{
					"vol":   v,
					"lsvol": v.Volume,
					"key":   key,
					"value": value[0]}).Info("service filter match")
				continue
			}
			return false
		default:
			log.WithFields(log.Fields{
				"vol":   v,
				"lsvol": v.Volume,
				"key":   key,
				"value": value[0]}).Info("non-standard filter key")

			match := false
			for k2, v2 := range v.Fields {
				if k2 != key {
					continue
				}
				if v2 == value[0] {
					match = true
				}
				log.WithFields(log.Fields{
					"vol":   v,
					"lsvol": v.Volume,
					"key":   key,
					"v2":    v2,
					"value": value[0]}).Info("reject key")

				break
			}
			if match {
				log.WithFields(log.Fields{
					"vol":   v,
					"lsvol": v.Volume,
					"key":   key,
					"value": value[0]}).Info("filter match")
				continue
			}
			log.WithFields(log.Fields{
				"vol":   v,
				"lsvol": v.Volume,
				"key":   key,
				"value": value[0]}).Info("rejected volume")
			return false
		}
	}
	fmt.Println("passed all filter(s) on Volume ", v.Name)
	return true
}

func (rtr *Router) getVolumesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vols, err := rtr.p.LsClient.Volumes()
	if err != nil {
		http.Error(w, "could not retrieve volumes",
			http.StatusInternalServerError)
		return
	}

	for _, vol := range vols {
		// todo use service name + volume name to retrieve polly metadata
		schedulerName := "AcmeSched"   // todo fill this from persistence store
		providerName := "AcmeProvider" // todo fill this in from persistence store
		vol.Scheduler = schedulerName
		vol.StorageProviderName = providerName
		vol.Labels["mycustomkey"] = "mycustomvalue"

		// // filter makes append decision for each volume
		// if applyVolumeFilter(volNew, r) {
		// 	vols = append(vols, volNew)
		// }
	}

	j, _ := json.Marshal(&vols)
	w.Write(j)
}

func (rtr *Router) postVolumesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var m *types.Volume
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

	if _, ok := rtr.p.Services[strings.ToLower(m.ServiceName)]; !ok {
		http.Error(w, "ServiceName is not defined", http.StatusNotFound)
		return
	}

	if m.Name == "" {
		http.Error(w, "mandatory Name missing or empty", 422)
		return
	}
	m.Volume.Name = m.Name

	if m.Scheduler == "" {
		http.Error(w, "mandatory Scheduler missing or empty", 422)
		return
	}

	// Volumes = append(Volumes, m)

	// These elements are optional with defaults
	m.AvailabilityZone = ""
	m.Type = "ext4"
	m.Size = int64(1000)
	m.IOPS = int64(100)
	opts := map[string]interface{}{}

	volumeCreateRequest := &apihttp.VolumeCreateRequest{
		Name:             m.Name,
		AvailabilityZone: &m.AvailabilityZone,
		Type:             &m.Type,
		Size:             &m.Size,
		IOPS:             &m.IOPS,
		Opts:             opts,
	}

	reply, err := rtr.p.LsClient.VolumeCreate(m.ServiceName, volumeCreateRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	volNew := &types.Volume{
		Volume:              reply,
		ServiceName:         m.ServiceName,
		Scheduler:           m.Scheduler,
		StorageProviderName: m.StorageProviderName,
		Labels:              make(map[string]string),
	}

	volNew.Labels["MyAdminField"] = "42"

	// todo persist this
	err = rtr.p.Store.SaveVolumeMetadata(volNew)
	if err != nil {
		log.Fatal("Failed to save metadata")
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
	}

	vs := strings.SplitN(volid, "-", 2)
	if len(vs) < 2 {
		http.Error(w, "valid volid is service-volume", http.StatusBadRequest)
		return
	}

	err := rtr.p.LsClient.VolumeRemove(vs[0], vs[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	return
}
