package store

import (
	"bytes"
	"log"

	"testing"

	gofig "github.com/akutz/gofig"
	types "github.com/emccode/libstorage/api/types"
)

const (
	libStorageConfigBaseBolt = `
polly:
  store:
    type: boltdb
    endpoints: /tmp/boltdb
    bucket: MyBoltDb
`
	libStorageConfigBaseConsul = `
polly:
  store:
    type: consul
    endpoints: 127.0.0.1:8500
`
)

func TestSaveData(t *testing.T) {
	config := gofig.New()

	configYamlBuf := []byte(libStorageConfigBaseBolt)
	if err := config.ReadConfig(bytes.NewReader(configYamlBuf)); err != nil {
		panic(err)
	}

	ps, err := NewWithConfig(config.Scope("polly.store"))
	if err != nil {
		log.Fatal("Failed to create PollyStore")
	}

	//save
	volume1 := types.Volume{ID: "myid1"}
	volume1.Fields = make(map[string]string)
	volume1.Fields["mykey1"] = "myvalue1"
	volume1.Fields["mykey2"] = "myvalue2"

	err = ps.SaveVolumeMetadata(&volume1)
	if err != nil {
		log.Fatal("Failed to save metadata")
	}

	//get
	volume2 := types.Volume{ID: "myid1"}
	volume2.Fields = make(map[string]string)
	err = ps.GetVolumeMetadata(&volume2)
	if err != nil {
		log.Fatal("Failed to save metadata")
	}

	log.Print("After GET")
	for key, value := range volume2.Fields {
		log.Print(key, " = ", value)
	}

	//update
	volume2.Fields["mykey1"] = "myvalue3"

	err = ps.SaveVolumeMetadata(&volume2)
	if err != nil {
		log.Fatal("Failed to save metadata")
	}

	log.Print("After UPDATE")
	for key, value := range volume2.Fields {
		log.Print(key, " = ", value)
	}

	//delte thru save
	volume2.Fields = make(map[string]string)
	volume2.Fields["mykey1"] = "myvalue3"

	err = ps.SaveVolumeMetadata(&volume2)
	if err != nil {
		log.Fatal("Failed to save metadata")
	}

	log.Print("After DELETE THRU SAVE")
	for key, value := range volume2.Fields {
		log.Print(key, " = ", value)
	}

	//delete all metadata
	err = ps.DeleteVolumeMetadata(&volume2)
	if err != nil {
		log.Fatal("Failed to save metadata")
	}
}
