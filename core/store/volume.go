package store

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/akutz/goof"
	store "github.com/docker/libkv/store"
	"github.com/emccode/polly/api/types"
)

// Exists returns true if a key for the specified Volume exists in the store
func (ps *PollyStore) Exists(volume *types.Volume) (bool, error) {
	key, err := ps.GenerateObjectKey(VolumeInternalLabelsType, volume.VolumeID)
	if err != nil {
		return false, err
	}

	log.WithField("key", key).Debug("checking for the existence of key")
	exists, err := ps.store.Exists(key)
	if err != nil {
		return exists, goof.WithError("problem checking key", err)
	}
	log.WithFields(log.Fields{
		"exists": exists,
		"key":    key,
	}).Debug("key check result")
	return exists, nil
}

//GetVolumeIds return an array of IDs for all volumes in the store
func (ps *PollyStore) GetVolumeIds() (ids []string, err error) {
	key, err := ps.GenerateRootKey(VolumeInternalLabelsType)
	if err != nil {
		return nil, err
	}

	kvpairs, err := ps.List(key)
	if err != nil {
		return nil, err
	}

	for _, pair := range kvpairs {
		path := strings.Split(pair.Key, "/")
		if len(path) == 4 && path[3] == "ID" {
			ids = append(ids, path[2])
		}
	}
	return
}

//SaveVolumeFields saves libstorage fields associated with a volume
func (ps *PollyStore) SaveVolumeFields(volume *types.Volume) error {
	// key, err := ps.GenerateObjectKey(VolumeType, volume.VolumeID)
	// if err != nil {
	// 	return err
	// }

	// ps.Put(key, []byte("")) // init directory at key
	//
	// //current keys in libkv
	// kvpairs, _ := ps.List(key)
	//
	// if kvpairs != nil {
	// 	//delete keys that are gone...
	// 	for _, pair := range kvpairs {
	// 		for _, vkey := range volume.Fields {
	// 			found := false
	// 			if strings.Compare(pair.Key, key+vkey) == 0 {
	// 				found = true
	// 			}
	// 			if !found {
	// 				ps.store.Delete(pair.Key)
	// 			}
	// 		}
	// 	}
	// }

	// //new and existing keys will be added here
	// for volKey, volValue := range volume.Fields {
	// 	err := ps.Put(key+volKey, []byte(volValue))
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

//SaveVolumeAdminLabels saves admin controlled labels associated with a volume
func (ps *PollyStore) SaveVolumeAdminLabels(volume *types.Volume) error {
	key, err := ps.GenerateObjectKey(VolumeAdminLabelsType, volume.VolumeID)
	if err != nil {
		return err
	}

	if err := ps.Put(key, []byte("")); err != nil {
		return err
	}

	kvpairs, _ := ps.List(key)

	if kvpairs != nil {
		for _, pair := range kvpairs {
			found := false
			for vkey := range volume.Labels {
				log.WithFields(log.Fields{
					"pair.key": pair.Key,
					"key+vkey": fmt.Sprintf("%s%s", key, vkey),
				}).Debug("comparing keys")
				if strings.Compare(pair.Key, key+vkey) == 0 {
					found = true
					break
				}
			}
			if !found {
				log.WithField("key", pair.Key).Debug("delete key")
				ps.store.Delete(pair.Key)
			}
		}
	}

	//new and existing keys will be added here
	for volKey, volValue := range volume.Labels {
		err := ps.Put(key+volKey, []byte(volValue))
		if err != nil {
			return err
		}
	}

	return nil
}

//SaveVolumeMetadata saves all metadata associated with a volume
func (ps *PollyStore) SaveVolumeMetadata(volume *types.Volume) error {
	key, err := ps.GenerateObjectKey(VolumeInternalLabelsType, volume.VolumeID)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"vol": volume,
		"key": key}).Info("saving volume metadata")

	if err = ps.Put(key, []byte("")); err != nil {
		return err
	}

	err = ps.Put(key+"ID", []byte(volume.VolumeID))
	if err != nil {
		return err
	}

	if volume.Schedulers == nil || len(volume.Schedulers) == 0 {
		err = ps.Delete(key + "Schedulers")
		if err != nil {
			return err
		}
	} else {
		js, err := json.Marshal(&volume.Schedulers)
		if err != nil {
			return err
		}

		err = ps.Put(key+"Schedulers", []byte(js))
		if err != nil {
			return err
		}
	}

	err = ps.Put(key+"ServiceName", []byte(volume.ServiceName))
	if err != nil {
		return err
	}

	err = ps.SaveVolumeFields(volume)
	if err != nil {
		return goof.New("failed to save fields")
	}

	err = ps.SaveVolumeAdminLabels(volume)
	if err != nil {
		return goof.New("failed to save fields")
	}

	return nil
}

//SetVolumeAdminLabels sets volume admin labels from persistent store
func (ps *PollyStore) SetVolumeAdminLabels(volume *types.Volume) error {
	key, err := ps.GenerateObjectKey(VolumeAdminLabelsType, volume.VolumeID)
	if err != nil {
		return err
	}

	if err = ps.Put(key, []byte("")); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"vol": volume,
		"key": key}).Info("saving volume admin labels")

	kvpairs, err := ps.List(key)
	if err == store.ErrKeyNotFound {
		ps.Put(key, []byte("")) // make empty list in store
		return nil
	} else if err != nil {
		return err
	}

	volume.Labels = make(map[string]string)
	for _, pair := range kvpairs {
		key, err := ps.GetKeyFromFQKN(pair.Key)
		if err != nil {
			continue
		}
		if key != "" {
			volume.Labels[key] = string(pair.Value)
		}
	}

	return nil
}

//SetVolumeMetadata This function will get all metadata associated with a volume
func (ps *PollyStore) SetVolumeMetadata(volume *types.Volume) (bool, error) {
	var exists bool
	var err error
	exists, err = ps.Exists(volume)
	if err != nil {
		return exists, err
	}

	if !exists {
		log.Debug("volume does not exist yet in store")
		return exists, nil
	}

	key, err := ps.GenerateObjectKey(VolumeInternalLabelsType, volume.VolumeID)
	if err != nil {
		return exists, err
	}

	//current keys in libkv
	kvpairs, err := ps.List(key)
	if err != nil {
		return exists, err
	}

	for _, pair := range kvpairs {
		key, err := ps.GetKeyFromFQKN(pair.Key)
		if err != nil {
			continue
		}
		switch key {
		case "Schedulers":
			err := json.Unmarshal(pair.Value, &volume.Schedulers)
			if err != nil {
				return exists, err
			}
		case "ServiceName":
			volume.ServiceName = string(pair.Value)
		}
	}

	err = ps.SetVolumeAdminLabels(volume)
	if err != nil {
		return exists, goof.New("Failed to retrieve volume labels from persistent store")
	}
	return exists, nil
}

//RemoveVolumeMetadata This function will save all metadata associated with a volume
func (ps *PollyStore) RemoveVolumeMetadata(volume *types.Volume) error {
	deletelist := []int{VolumeType, VolumeInternalLabelsType, VolumeAdminLabelsType}
	for _, deleteme := range deletelist {
		key, err := ps.GenerateObjectKey(deleteme, volume.VolumeID)
		if err != nil {
			return err
		}

		err = ps.store.DeleteTree(key)
		if err != nil {
			return err
		}
	}

	return nil
}
