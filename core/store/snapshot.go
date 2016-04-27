package store

import (
	// log "github.com/Sirupsen/logrus"
	"strings"

	"github.com/emccode/polly/api/types"
)

//SaveSnapshotMetadata This function will save all metadata associated with a volume
func (ps *PollyStore) SaveSnapshotMetadata(snapshot *types.Snapshot) error {
	key, err := ps.GenerateObjectKey(SnapshotType, snapshot.SnapshotID)
	if err != nil {
		return err
	}

	//current keys in libkv
	kvpairs, _ := ps.List(key)

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
			return err
		}
	}

	return nil
}

//GetSnapshotMetadata This function will get all metadata associated with a volume
func (ps *PollyStore) GetSnapshotMetadata(snapshot *types.Snapshot) error {
	key, err := ps.GenerateObjectKey(SnapshotType, snapshot.SnapshotID)
	if err != nil {
		return err
	}

	//current keys in libkv
	kvpairs, err := ps.List(key)
	if err != nil {
		return err
	}

	snapshot.Fields = make(map[string]string)
	for _, pair := range kvpairs {
		key, err := ps.GetKeyFromFQKN(pair.Key)
		if err != nil {
			continue
		}
		snapshot.Fields[key] = string(pair.Value)
	}

	return nil
}

//DeleteSnapshotMetadata This function will save all metadata associated with a volume
func (ps *PollyStore) DeleteSnapshotMetadata(snapshot *types.Snapshot) error {
	key, err := ps.GenerateObjectKey(SnapshotType, snapshot.SnapshotID)
	if err != nil {
		return err
	}

	err = ps.store.DeleteTree(key)
	if err != nil {
		return err
	}

	return nil
}
