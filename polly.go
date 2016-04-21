/*
Polly provides a a REST API based on libstorage,
The REST API provides both administrator and scheduler functions.
*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/emccode/libstorage/api/types"
	apihttp "github.com/emccode/libstorage/api/types/http"
	"github.com/gorilla/mux"
)

// Volumes is array of known volumes
var Volumes = []Volume{}

// Volume is mandatory fields for a volume
type Volume struct {
	types.Volume

	// the Storage provider identifier
	ServiceName string `json:"serviceName,omitempty"`
}

func applyVolumeFilter(v *Volume, r *http.Request) bool {
	for key, value := range r.URL.Query() {
		fmt.Printf("applyVolumeFilter Key:%s=%s\n", key, value[0])
		switch key {
		case "availabilityZone":
			if v.AvailabilityZone == value[0] {
				fmt.Println("availabilityZone filter match", value[0], "on volume ", v.Name)
				continue
			}
			fmt.Println("REJECT volume ", v.Name, " availabilityZone", v.AvailabilityZone, "!=", value[0])
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
			return false
		case "scheduler":
			if v.ServiceName == value[0] {
				fmt.Println("scheduler filter match", value[0], "on volume ", v.Name)
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
				fmt.Printf("REJECT key %s %s != %s on Volume %s", key, v2, value[0], v.Name)
				break
			}
			if match {
				fmt.Printf("filter math key %s == %s on Volume %s", key, value[0], v.Name)
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

	serviceVolumeMap, err := lsClient.Volumes(false)
	if err != nil {
		http.Error(w, "could not retrieve volume map", http.StatusInternalServerError)
		return
	}

	// log entire query clause
	fmt.Println(r.URL.RawQuery)
	for key, value := range r.URL.Query() {
		fmt.Println("Key:", key, "=Value:", value)
	}

	var vols []*Volume
	for serviceName, volumeMap := range serviceVolumeMap {
		for _, vol := range volumeMap {
			volNew := &Volume{*vol, serviceName}
			// read persistence layer here
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

	var m Volume
	b, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(b, &m)
	if err != nil {
		http.Error(w, "json is unparsable", http.StatusBadRequest)
		return
	}

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
	Volumes = append(Volumes, m)

	az := "myzone"
	voltype := "mytype"
	sz := int64(1000)
	iops := int64(500)
	opts := map[string]interface{}{}

	volumeCreateRequest := &apihttp.VolumeCreateRequest{
		Name:             m.Name,
		AvailabilityZone: &az,
		Type:             &voltype,
		Size:             &sz,
		IOPS:             &iops,
		Opts:             opts,
	}

	reply, err := lsClient.VolumeCreate(m.ServiceName, volumeCreateRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	volNew := &Volume{*reply, m.ServiceName}

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

func main() {
	r := mux.NewRouter()
	//volumes
	r.HandleFunc("/admin/volumes", getVolumesHandler).Methods("GET")
	r.HandleFunc("/admin/volumes", postVolumesHandler).Methods("POST")
	r.HandleFunc("/admin/volumes/{volid}", deleteVolumesHandler).Methods("DELETE")

	http.Handle("/", r)
	go http.ListenAndServe(":8080", nil)

	server, errs := startServer()

	if server == nil {
		fmt.Print("server == nil\n")
	}

	if err := startClient(); err != nil {
		server.Close()
		panic(err)
	}

	go func(errs <-chan error) {
		err := <-errs
		if err != nil {
			server.Close()
			fmt.Print("server error: ", err)
			// todo should trigger a syscall to shout down
		}
	}(errs)

	forever := make(chan bool)
	<-forever
}
