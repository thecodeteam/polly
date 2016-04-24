package polly

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	gofig "github.com/akutz/gofig"
	lstypes "github.com/emccode/libstorage/api/types"
	apihttp "github.com/emccode/libstorage/api/types/http"
	core "github.com/emccode/polly/core"
	"github.com/emccode/polly/pkg/store"

	lsclient "github.com/emccode/libstorage/client"
	"github.com/emccode/polly/types"
	"github.com/gorilla/mux"
)

//New init the lib
func New() (p *core.PollyCore, err error) {
	config := gofig.New()

	myErr := config.ReadConfigFile("/etc/polly/config.yml")
	if myErr != nil {
		log.Fatal("Fatal: ", myErr)
		return nil, myErr
	}

	pollyCore, myErr := core.NewWithConfig(config)
	if myErr != nil {
		log.Fatal("Fatal: ", myErr)
		return nil, myErr
	}
	log.Print("PollyStore Type: ", pollyCore.PollyStore.StoreType())

	return pollyCore, nil
}

//NewWithConfigFile init the lib
func NewWithConfigFile(path string) (p *core.PollyCore, err error) {
	config := gofig.New()

	myErr := config.ReadConfigFile(path)
	if myErr != nil {
		log.Fatal("Fatal: ", myErr)
		return nil, myErr
	}

	pollyCore, myErr := core.NewWithConfig(config)
	if myErr != nil {
		log.Fatal("Fatal: ", myErr)
		return nil, myErr
	}
	log.Print("PollyStore Type: ", pollyCore.PollyStore.StoreType())

	return pollyCore, nil
}

// Volumes is a slice of known volumes
var Volumes = []*types.Volume{}

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

var services apihttp.ServicesMap
var lsClient lsclient.Client
var ps *store.PollyStore

func main() {
	cfg, err := startServer()
	if err != nil {
		panic(err)
	}

	if err := startClient(cfg); err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	//volumes
	r.HandleFunc("/admin/volumes", getVolumesHandler).Methods("GET")
	r.HandleFunc("/admin/volumes", postVolumesHandler).Methods("POST")
	r.HandleFunc("/admin/volumes/{volid}", deleteVolumesHandler).Methods("DELETE")

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
			exists, err := ps.Exists(vol)
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
