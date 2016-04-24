package store

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	store "github.com/docker/libkv/store"
	"github.com/emccode/polly/types"
)

// Exists returns true if a key for the specified Volume exists in the store
func (ps *PollyStore) Exists(volume *types.Volume) (bool, error) {
	key, err := ps.GenerateKey(VolumeInternalLabelsType, volume.ID)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	fmt.Println("checking", key)
	exists, err := ps.store.Exists(key)
	fmt.Println(key, "exists=", exists, "err=", err)
	return exists, err
}

//GetVolumeIds return an array of IDs for all volumes in the store
func (ps *PollyStore) GetVolumeIds() (ids []string, err error) {

	key, err := ps.GenerateKey(VolumeInternalLabelsType, "")
	if err != nil {
		log.Fatal("Generate key failed", err)
		return ids, err
	}

	kvpairs, err := ps.store.List(key)
	if err != nil {
		log.Fatal("List key failed ", key, err)
		return
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
	key, err := ps.GenerateKey(VolumeType, volume.ID)
	if err != nil {
		log.Fatal(err)
		return err
	}

	ps.store.Put(key, []byte(""), nil) // init directory at key

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
		fmt.Println("write ", key+volKey, "=", volValue)
		err := ps.store.Put(key+volKey, []byte(volValue), nil)
		if err != nil {
			log.Fatal("Fatal: ", err)
			return err
		}
	}

	return nil
}

//SaveVolumeAdminLabels saves admin controlled labels associated with a volume
func (ps *PollyStore) SaveVolumeAdminLabels(volume *types.Volume) error {
	key, err := ps.GenerateKey(VolumeAdminLabelsType, volume.ID)
	if err != nil {
		log.Fatal(err)
		return err
	}

	//current keys in libkv
	kvpairs, _ := ps.store.List(key)

	if kvpairs != nil {
		//delete keys that are gone...
		for _, pair := range kvpairs {
			found := false
			for _, vkey := range volume.Labels {
				if strings.Compare(pair.Key, key+vkey) == 0 {
					found = true
					break
				}
			}
			if !found {
				ps.store.Delete(pair.Key)
			}
		}
	}

	//new and existing keys will be added here
	for volKey, volValue := range volume.Labels {
		err := ps.store.Put(key+volKey, []byte(volValue), nil)
		if err != nil {
			log.Fatal("Fatal: ", err)
			return err
		}
	}

	return nil
}

//SaveVolumeMetadata saves all metadata associated with a volume
func (ps *PollyStore) SaveVolumeMetadata(volume *types.Volume) error {
	key, err := ps.GenerateKey(VolumeInternalLabelsType, volume.ID)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("write ", key+"AvailabilityZone", "=", volume.AvailabilityZone)

	err = ps.store.Put(key+"AvailabilityZone", []byte(volume.AvailabilityZone),
		nil)
	if err != nil {
		log.Fatal("Fatal: ", err)
		return err
	}
	err = ps.store.Put(key+"ID", []byte(volume.ID), nil)
	if err != nil {
		log.Fatal("Fatal: ", err)
		return err
	}
	err = ps.store.Put(key+"IOPS", []byte(strconv.FormatInt(volume.IOPS, 10)),
		nil)
	if err != nil {
		log.Fatal("Fatal: ", err)
		return err
	}
	err = ps.store.Put(key+"NetworkName", []byte(volume.NetworkName), nil)
	if err != nil {
		log.Fatal("Fatal: ", err)
		return err
	}
	err = ps.store.Put(key+"Scheduler", []byte(volume.Scheduler), nil)
	if err != nil {
		log.Fatal("Fatal: ", err)
		return err
	}
	err = ps.store.Put(key+"ServiceName", []byte(volume.ServiceName), nil)
	if err != nil {
		log.Fatal("Fatal: ", err)
		return err
	}
	err = ps.store.Put(key+"Size", []byte(strconv.FormatInt(volume.Size, 10)),
		nil)
	if err != nil {
		log.Fatal("Fatal: ", err)
		return err
	}
	err = ps.store.Put(key+"Status", []byte(volume.Status), nil)
	if err != nil {
		log.Fatal("Fatal: ", err)
		return err
	}
	err = ps.store.Put(key+"StorageProviderName",
		[]byte(volume.StorageProviderName), nil)
	if err != nil {
		log.Fatal("Fatal: ", err)
		return err
	}
	err = ps.store.Put(key+"Type", []byte(volume.Type), nil)
	if err != nil {
		log.Fatal("Fatal: ", err)
		return err
	}

	err = ps.SaveVolumeFields(volume)
	if err != nil {
		log.Fatal("Failed to save fields")
	}
	err = ps.SaveVolumeAdminLabels(volume)
	if err != nil {
		log.Fatal("Failed to save fields")
	}

	return nil
}

//SetVolumeFields sets volume fields associated with libstorage
func (ps *PollyStore) SetVolumeFields(volume *types.Volume) error {
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

//SetVolumeAdminLabels sets volume admin labels from persistent store
func (ps *PollyStore) SetVolumeAdminLabels(volume *types.Volume) error {
	key, err := ps.GenerateKey(VolumeAdminLabelsType, volume.ID)
	if err != nil {
		log.Fatal(err)
		return err
	}

	kvpairs, err := ps.store.List(key)
	if err == store.ErrKeyNotFound {
		ps.store.Put(key, []byte(""), nil) // make empty list in store
		return nil
	} else if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("got List of", len(kvpairs))
	volume.Labels = make(map[string]string)
	for _, pair := range kvpairs {
		key, err := ps.GetKeyFromFQKN(pair.Key)
		if err != nil {
			log.Fatal("Invalid key format.")
			continue
		}
		volume.Labels[key] = string(pair.Value)
	}

	return nil
}

//SetVolumeMetadata This function will get all metadata associated with a volume
func (ps *PollyStore) SetVolumeMetadata(volume *types.Volume) error {
	if err := ps.SetVolumeFields(volume); err != nil {
		log.Fatal("Failed to retrieve volume fields from persistent store")
	}

	key, err := ps.GenerateKey(VolumeInternalLabelsType, volume.ID)
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

	for _, pair := range kvpairs {
		key, err := ps.GetKeyFromFQKN(pair.Key)
		if err != nil {
			log.Fatal("Invalid key format.")
			continue
		}
		switch key {
		case "AvailabilityZone":
			volume.AvailabilityZone = string(pair.Value)
		case "ID":
			volume.ID = string(pair.Value)
		case "IOPS":
			volume.IOPS, err = strconv.ParseInt(string(pair.Value), 10, 64)
			if err != nil {
				volume.IOPS = 0
			}
		case "NetworkName":
			volume.NetworkName = string(pair.Value)
		case "Scheduler":
			volume.Scheduler = string(pair.Value)
		case "ServiceName":
			volume.ServiceName = string(pair.Value)
		case "Size":
			volume.Size, err = strconv.ParseInt(string(pair.Value), 10, 64)
			if err != nil {
				volume.Size = 0
			}
		case "Status":
			volume.Status = string(pair.Value)
		case "StorageProviderName":
			volume.StorageProviderName = string(pair.Value)
		case "Type":
			volume.Type = string(pair.Value)
		}
	}

	err = ps.SetVolumeAdminLabels(volume)
	if err != nil {
		log.Fatal("Failed to retrieve volume labels from persistent store")
	}
	return nil
}

//DeleteVolumeMetadata This function will save all metadata associated with a volume
func (ps *PollyStore) DeleteVolumeMetadata(volume *types.Volume) error {
	deletelist := []int{VolumeType, VolumeInternalLabelsType, VolumeAdminLabelsType}
	for _, deleteme := range deletelist {
		key, err := ps.GenerateKey(deleteme, volume.ID)
		if err != nil {
			log.Fatal(err)
			return err
		}

		err = ps.store.DeleteTree(key)
		if err != nil {
			log.Fatal("Fatal: ", err)
			return err
		}
	}

	return nil
}

// EraseStore erases the store
func (ps *PollyStore) EraseStore() {
	ps.store.DeleteTree("/volumeinternal")
	ps.store.DeleteTree("/volumeadmin")
	ps.store.DeleteTree("/volume")
}
