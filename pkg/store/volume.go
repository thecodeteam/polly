package store

import (
	"log"
	"strings"

	types "github.com/emccode/libstorage/api/types"
)

//SaveVolumeMetadata This function will save all metadata associated with a volume
func (ps *PollyStore) SaveVolumeMetadata(volume *types.Volume) error {
	key, err := ps.GenerateKey(VolumeType, volume.ID)
	if err != nil {
		log.Fatal(err)
		return err
	}

	//current keys in libkv
	kvpairs, _ := ps.store.List(key)

	if kvpairs != nil {
		//delete keys that are gone...
		for _, pair := range kvpairs {
			for _, vkey := range volume.Fields {
				found := false
				if strings.Compare(pair.Key, key+vkey) == 0 {
					found = true
				}
				if !found {
					ps.store.Delete(pair.Key)
				}
			}
		}
	}

	//new and existing keys will be added here
	for volKey, volValue := range volume.Fields {
		err := ps.store.Put(key+volKey, []byte(volValue), nil)
		if err != nil {
			log.Fatal("Fatal: ", err)
			return err
		}
	}

	return nil
}

//GetVolumeMetadata This function will get all metadata associated with a volume
func (ps *PollyStore) GetVolumeMetadata(volume *types.Volume) error {
	key, err := ps.GenerateKey(VolumeType, volume.ID)
	if err != nil {
		log.Fatal(err)
		return err
	}

	//current keys in libkv
	kvpairs, err := ps.store.List(key)
	if err != nil {
		log.Fatal(err)
		return err
	}

	volume.Fields = make(map[string]string)
	for _, pair := range kvpairs {
		key, err := ps.GetKeyFromFQKN(pair.Key)
		if err != nil {
			log.Fatal("Invalid key format.")
			continue
		}
		volume.Fields[key] = string(pair.Value)
	}

	return nil
}

//DeleteVolumeMetadata This function will save all metadata associated with a volume
func (ps *PollyStore) DeleteVolumeMetadata(volume *types.Volume) error {
	key, err := ps.GenerateKey(VolumeType, volume.ID)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = ps.store.DeleteTree(key)
	if err != nil {
		log.Fatal("Fatal: ", err)
		return err
	}

	return nil
}
