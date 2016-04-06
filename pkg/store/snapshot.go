package store

import (
	"log"
	"strings"

	types "github.com/emccode/libstorage/api/types"
)

//SaveSnapshotMetadata This function will save all metadata associated with a volume
func (ps *PollyStore) SaveSnapshotMetadata(snapshot *types.Snapshot) error {
	key, err := ps.GenerateKey(SnapshotType, snapshot.ID)
	if err != nil {
		log.Fatal(err)
		return err
	}

	//current keys in libkv
	kvpairs, _ := ps.store.List(key)

	if kvpairs != nil {
		//delete keys that are gone...
		for _, pair := range kvpairs {
			for _, vkey := range snapshot.Fields {
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

	for snapKey, snapValue := range snapshot.Fields {
		err := ps.store.Put(key+snapKey, []byte(snapValue), nil)
		if err != nil {
			log.Fatal("Fatal: ", err)
			return err
		}
	}

	return nil
}

//GetSnapshotMetadata This function will get all metadata associated with a volume
func (ps *PollyStore) GetSnapshotMetadata(snapshot *types.Snapshot) error {
	key, err := ps.GenerateKey(SnapshotType, snapshot.ID)
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

	snapshot.Fields = make(map[string]string)
	for _, pair := range kvpairs {
		key, err := ps.GetKeyFromFQKN(pair.Key)
		if err != nil {
			log.Fatal("Invalid key format.")
			continue
		}
		snapshot.Fields[key] = string(pair.Value)
	}

	return nil
}

//DeleteSnapshotMetadata This function will save all metadata associated with a volume
func (ps *PollyStore) DeleteSnapshotMetadata(snapshot *types.Snapshot) error {
	key, err := ps.GenerateKey(SnapshotType, snapshot.ID)
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
