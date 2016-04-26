/*
Polly provides a a REST API based on libstorage,
The REST API provides both administrator and scheduler functions.
*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	lstypes "github.com/emccode/libstorage/api/types"
	apihttp "github.com/emccode/libstorage/api/types/http"
	"github.com/emccode/polly/types"
	"github.com/gorilla/mux"
)

// VersionPollyAdminAPI is version of admin API
var VersionPollyAdminAPI = "0.1.0"

// VersionPollySchedulerAPI is version of scheduler API
var VersionPollySchedulerAPI = "0.1.0"

// VersionPollyBuild is build version of polly
var VersionPollyBuild = "0.0.0"

// Volumes is a slice of known volumes
var Volumes = []*types.Volume{}

// getVersionHandler is gorilla mux handler for GET version on REST API
func getVersionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ver := &types.PollyVersion{
		PollyAdminAPIVersion:     VersionPollyAdminAPI,
		PollySchedulerAPIVersion: VersionPollySchedulerAPI,
		PollyBuildVersion:        VersionPollyBuild,
	}

	j, _ := json.Marshal(&ver)
	w.Write(j)
}

// errorHandler is gorilla mux handler for various REST API resources
func errorHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", 405)
	return
}

func getServicesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get provider list from libstorage (temporary using services instead)
	serviceMap, err := lsClient.Services()
	if err != nil {
		http.Error(w, "could not retrieve service map",
			http.StatusInternalServerError)
		return
	}

	var provs []*types.Provider
	for _, servinfo := range serviceMap {
		provNew := &types.Provider{
			ProviderName: servinfo.Name,
			DriverName:   servinfo.Driver.Name,
			Labels:       make(map[string]string),
		}
		provs = append(provs, provNew)
	}

	j, _ := json.Marshal(&provs)
	w.Write(j)
}

func applyVolumeFilter(v *types.Volume, r *http.Request) bool {
	for key, value := range r.URL.Query() {
		fmt.Printf("applyVolumeFilter Key:%s=%s\n", key, value[0])
		switch key {
		case "availabilityZone":
			if v.AvailabilityZone == value[0] {
				fmt.Println("availabilityZone filter match", value[0], "on volume ",
					v.Name)
				continue
			}
			fmt.Println("REJECT volume ", v.Name, " availabilityZone",
				v.AvailabilityZone, "!=", value[0])
			return false
		case "iops":
			i, err := strconv.ParseInt(value[0], 10, 64)
			if err != nil {
				return false
			}
			if v.IOPS == i {
				fmt.Println("iops filter match", value[0], "on volume ", v.Name)
				continue
			}
			fmt.Println("REJECT volume ", v.Name, " iops", v.IOPS, "!=", i)
			return false
		case "size":
			i, err := strconv.ParseInt(value[0], 10, 64)
			if err != nil {
				return false
			}
			if v.Size == i {
				fmt.Println("size filter match", value[0], "on volume ", v.Name)
				continue
			}
			fmt.Println("REJECT volume ", v.Name, " size", v.Size, "!=", i)
			return false
		case "serviceName":
			if v.ServiceName == value[0] {
				fmt.Println("serviceName filter match", value[0], "on volume ", v.Name)
				continue
			}
			return false
		default:
			// non-standard filter key
			fmt.Printf("non-standard filter Key:%s=%s\n", key, value[0])
			match := false
			for k2, v2 := range v.Fields {
				if k2 != key {
					continue
				}
				if v2 == value[0] {
					match = true
				}
				fmt.Printf("REJECT key %s %s != %s on Volume %s", key, v2, value[0],
					v.Name)
				break
			}
			if match {
				fmt.Printf("filter math key %s == %s on Volume %s", key, value[0],
					v.Name)
				continue
			}
			fmt.Printf("REJECT, key %s not defined on volume %s", key, v.Name)
			return false
		}
	}
	fmt.Println("passed all filter(s) on Volume ", v.Name)
	return true
}

func getVolumesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get volume list from libstorage
	// note that libstorage returns a limited amount of metadata for each volume
	// Polly augments the per volume metatdata with items retained in the Polly
	// persistent store
	serviceVolumeMap, err := lsClient.Volumes(false)
	if err != nil {
		http.Error(w, "could not retrieve volume map",
			http.StatusInternalServerError)
		return
	}

	var vols []*types.Volume
	for serviceName, volumeMap := range serviceVolumeMap {
		for _, vol := range volumeMap {
			// todo use service name + volume name to retrieve polly metadata
			schedulerName := "AcmeSched"   // todo fill this from persistence store
			providerName := "AcmeProvider" // todo fill this in from persistence store

			volNew := &types.Volume{
				Volume:              vol,
				ServiceName:         serviceName,
				Scheduler:           schedulerName,
				StorageProviderName: providerName,
				Labels:              make(map[string]string),
			}

			volNew.Labels["mycustomkey"] = "mycustomvalue"

			// filter makes append decision for each volume
			if applyVolumeFilter(volNew, r) {
				vols = append(vols, volNew)
			}
		}
	}

	j, _ := json.Marshal(&vols)
	w.Write(j)
}

func postVolumesHandler(w http.ResponseWriter, r *http.Request) {
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

	if _, ok := services[strings.ToLower(m.ServiceName)]; !ok {
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

	Volumes = append(Volumes, m)

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

	reply, err := lsClient.VolumeCreate(m.ServiceName, volumeCreateRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	volNew := &types.Volume{
		Volume:              reply,
		ServiceName:         m.ServiceName,
		Scheduler:           m.Scheduler,
		StorageProviderName: m.StorageProviderName,
	}

	volNew.Labels["MyAdminField"] = "42"

	// todo persist this
	err = ps.SaveVolumeMetadata(volNew)
	if err != nil {
		log.Fatal("Failed to save metadata")
	}

	j, _ := json.Marshal(volNew)

	w.Write(j)
}

func deleteVolumesHandler(w http.ResponseWriter, r *http.Request) {
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

	err := lsClient.VolumeRemove(vs[0], vs[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	return
}

// notAllowedHandler is used to set a status code for unsupported operations
func notAllowedHandler(allow ...string) func(w http.ResponseWriter, req *http.Request) {
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
func notImplementedHandler() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "Not Implemented Yet", http.StatusNotImplemented)
		return
	}
}

func main() {

	cfg, err := startServer()
	if err != nil {
		panic(err)
	}

	config = cfg

	if err = startClient(); err != nil {
		panic(err)
	}
	fmt.Println("client started")

	r := mux.NewRouter()
	//volumes
	r.HandleFunc("/admin/volumes", getVolumesHandler).Methods("GET")
	r.HandleFunc("/admin/volumes", postVolumesHandler).Methods("POST")
	r.HandleFunc("/admin/version",
		notAllowedHandler("GET", "POST")).Methods("PUT", "PATCH", "DELETE")
	r.HandleFunc("/admin/volumes/{volid}", deleteVolumesHandler).Methods("DELETE")
	r.HandleFunc("/admin/version", getVersionHandler).Methods("GET")
	r.HandleFunc("/admin/version",
		notAllowedHandler("GET")).Methods("POST", "PUT", "PATCH", "DELETE")

	//storage providers
	r.HandleFunc("/admin/storageproviders", getServicesHandler).Methods("GET")
	r.HandleFunc("/admin/storageproviders",
		notImplementedHandler()).Methods("POST")
	r.HandleFunc("/admin/storageproviders",
		notAllowedHandler("GET", "POST")).Methods("PUT", "PATCH", "DELETE")
	r.HandleFunc("/admin/storageproviders/{provid}",
		notImplementedHandler()).Methods("GET", "POST", "PUT", "PATCH", "DELETE")

	r.HandleFunc("/admin/storagepools", getServicesHandler).Methods("GET")
	r.HandleFunc("/admin/storagepools",
		notImplementedHandler()).Methods("POST")
	r.HandleFunc("/admin/storagepools",
		notAllowedHandler("GET", "POST")).Methods("PUT", "PATCH", "DELETE")
	r.HandleFunc("/admin/storagepools/{poolid}",
		notImplementedHandler()).Methods("GET", "POST", "PUT", "PATCH", "DELETE")

	http.Handle("/", r)
	go http.ListenAndServe(":8080", nil)

	// iterate over all volumes from libstorage
	// if any are missing from the Polly store, add them now, with scheduler id = "UNCLAIMED"
	serviceVolumeMap, err := lsClient.Volumes(false)
	if err != nil {
		panic(err)
	}

	ps.EraseStore()

	for serviceName, volumeMap := range serviceVolumeMap {
		for _, lsvol := range volumeMap {
			vol := &types.Volume{
				Volume:              lsvol,
				ServiceName:         serviceName,
				Scheduler:           "",
				StorageProviderName: "libstorage",
				Labels:              make(map[string]string),
			}

			var exists bool
			exists, err = ps.Exists(vol)
			if err != nil {
				panic(err)
			}
			if !exists {
				ps.SaveVolumeMetadata(vol)
				fmt.Println("Added unclaimed libstorage volume to Polly persistent store ",
					vol.Name, " as ", vol.ID)
			}
		}
	}

	// read in volume records from Polly persistent store and cache them
	ids, err := ps.GetVolumeIds()
	if err != nil {
		log.Fatal("Failed to retrieve volume IDs from persistent store", err)
	}

	for _, id := range ids {
		v := &types.Volume{
			Volume: &lstypes.Volume{
				ID: id,
			},
		}
		err = ps.SetVolumeMetadata(v)
		if err != nil {
			log.Fatal("Failed to retrieve volume metadata from persistent store")
		}
		Volumes = append(Volumes, v)
	}
	j, _ := json.Marshal(&Volumes)
	fmt.Println("volumes=", string(j))

	forever := make(chan bool)
	<-forever
}
