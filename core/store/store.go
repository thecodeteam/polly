package store

import (
	"errors"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	store "github.com/docker/libkv/store"
	"github.com/docker/libkv/store/boltdb"
	"github.com/docker/libkv/store/consul"
	"github.com/docker/libkv/store/etcd"
	"github.com/docker/libkv/store/zookeeper"

	"github.com/akutz/gofig"
	"github.com/akutz/goof"
	"github.com/docker/libkv"
	version "github.com/emccode/polly/core/version"
)

const (
	//VolumeType is used to identify metadata for the libstorage layer
	VolumeType = 1
	//VolumeInternalLabelsType is used to identify metadata for the Polly admin layer
	VolumeInternalLabelsType = 2
	//VolumeAdminLabelsType is used to identify labels for the Polly admin layer
	VolumeAdminLabelsType = 3
)

const (
	storeVolumeLibStorage         = "volumelibstorage"
	storeVolumeInternalLabelsType = "volumeinternallabels"
	storeVolumeAdminLabelsType    = "volumeadminlabels"
	rootKey                       = "polly"
)

var (
	//ErrObjectInvalid is an error for not being able to determine the key for an object
	ErrObjectInvalid = errors.New("Unable to determine root object path")
	//ErrInvalidStore is an error for not being able to determine the key for an object
	ErrInvalidStore = errors.New("Invalid store type")
	//ErrIndexOutOfBounds parsing a string has resulted in a -1
	ErrIndexOutOfBounds = errors.New("Index out of bounds")
)

//PollyStore representation of PollyStore
type PollyStore struct {
	config gofig.Config
	store  store.Store
	root   string
}

func init() {
	gofig.Register(configRegistration())
}

//NewWithConfig This initializes new instance of this library
func NewWithConfig(config gofig.Config) (pollystore *PollyStore, err error) {
	cfg := store.Config{
		ConnectionTimeout: 10 * time.Second,
	}

	ps := &PollyStore{}
	ps.config = config
	ps.root = rootKey + "/"

	backend := store.Backend(ps.StoreType())
	endpoints := []string{ps.EndPoints()}
	root := ps.Root()
	if len(root) > 0 {
		ps.root = root
	}

	switch backend {
	case store.CONSUL:
		consul.Register()
	case store.ZK:
		zookeeper.Register()
	case store.ETCD:
		etcd.Register()
	case store.BOLTDB:
		boltdb.Register()
		cfg.Bucket = ps.Bucket()
	default:
		return nil, ErrInvalidStore
	}

	myStore, err := libkv.NewStore(backend, endpoints, &cfg)
	if err != nil {
		log.Fatal("Fatal: ", err)
		return nil, err
	}
	ps.store = myStore

	pair, err := ps.store.Get(ps.versionKey())
	if pair != nil && err != nil {
		log.WithFields(log.Fields{
			"store": string(pair.Value),
		}).Debug("store version")
		log.WithFields(log.Fields{
			"version": version.VersionStr,
		}).Debug("current version")
		//TODO: if you need to some form of metadata update, do it here
	}

	//record the current version for the metadata
	ps.Put(ps.root, []byte(""))
	err = ps.store.Put(ps.versionKey(), []byte(version.VersionStr), nil)
	if err != nil {
		log.WithError(err).Fatal("failed to set version on store")
		return nil, err
	}

	if err := ps.initKeys([]int{VolumeType,
		VolumeInternalLabelsType, VolumeAdminLabelsType}); err != nil {
		return nil, err
	}

	return ps, nil
}

func (ps *PollyStore) initKeys(keys []int) error {
	for _, cat := range keys {
		key, err := ps.GenerateRootKey(cat)
		if err != nil {
			log.WithError(err).Fatal("failed key gen on store init")
			return err
		}
		ps.Put(key, []byte(""))
	}
	return nil
}

func (ps *PollyStore) versionKey() string {
	return ps.root + "version"
}

// Put saves key value pairs
func (ps *PollyStore) Put(key string, bytes []byte) error {
	log.WithFields(log.Fields{
		"key":         key,
		"bytes":       bytes,
		"bytesString": string(bytes),
	}).Debug("putting key value")
	return ps.store.Put(key, bytes, nil)
}

// Delete removes a key value pair
func (ps *PollyStore) Delete(key string) error {
	log.WithFields(log.Fields{
		"key": key,
	}).Debug("deleting key value")
	return ps.store.Delete(key)
}

//GenerateObjectKey generates the internal path (=key) for an object
func (ps *PollyStore) GenerateObjectKey(mytype int, guid string) (path string, err error) {
	if guid == "" {
		return "", goof.New("no guid provided")
	}
	return ps.generateKey(mytype, guid)
}

//GenerateRootKey generates the internal path for root keys
func (ps *PollyStore) GenerateRootKey(mytype int) (path string, err error) {
	return ps.generateKey(mytype, "")
}

//GenerateKey generates the internal path (=key) for persisting a value
func (ps *PollyStore) generateKey(mytype int, guid string) (path string, err error) {
	var parts = make([]string, 0, 3)

	root := ps.Root()
	if len(root) > 0 {
		parts = append(parts, root)
	}

	switch mytype {
	case VolumeType:
		parts = append(parts, storeVolumeLibStorage)
	case VolumeInternalLabelsType:
		parts = append(parts, storeVolumeInternalLabelsType)
	case VolumeAdminLabelsType:
		parts = append(parts, storeVolumeAdminLabelsType)
	default:
		return "", ErrObjectInvalid
	}

	if len(guid) > 0 {
		guid = strings.TrimSpace(guid)
		guid = strings.Trim(guid, "/") // don't allow slash in keys
		parts = append(parts, guid)
	}

	path = strings.Join(parts, "/")
	path += "/"
	return
}

//GetKeyFromFQKN get the key name from the FQKN
func (ps *PollyStore) GetKeyFromFQKN(fqkn string) (mykey string, err error) {
	pos := strings.LastIndex(fqkn, "/")
	if pos == -1 {
		return "", ErrIndexOutOfBounds
	}

	return fqkn[pos+1:], nil
}

// List lists the key values pairs for a key
func (ps *PollyStore) List(key string) ([]*store.KVPair, error) {
	list, err := ps.store.List(key)
	if err != nil {
		return nil, goof.WithError("problem listing key values", err)
	}
	log.WithFields(log.Fields{
		"key":  key,
		"list": list,
	}).Debug("got key values")
	if os.Getenv("POLLY_DEBUG") == "true" {
		for k, v := range list {
			log.WithFields(log.Fields{
				"key":   k,
				"value": v,
			})
		}
	}
	return list, nil
}

// EraseStore erases the store
func (ps *PollyStore) EraseStore() error {
	log.WithField("store", ps.store).Warning("erasing polly store trees")
	for _, t := range []int{
		VolumeInternalLabelsType, VolumeType, VolumeInternalLabelsType} {
		if err := ps.EraseType(t); err != nil {
			return err
		}
	}

	return nil
}

// EraseType erases a root of the store
func (ps *PollyStore) EraseType(mytype int) error {
	var rkey string
	var err error
	if rkey, err = ps.GenerateRootKey(mytype); err != nil {
		return err
	}
	if err := ps.store.DeleteTree(rkey); err != nil {
		return err
	}
	ps.Put(rkey, []byte(""))
	return nil
}

//StoreType this generates the type of backing store to use
func (ps *PollyStore) StoreType() string {
	return ps.config.GetString("polly.store.type")
}

//Root this get the type of backing store to use
func (ps *PollyStore) Root() string {
	return ps.config.GetString("polly.store.root")
}

//EndPoints this gets the endpoints of the store
func (ps *PollyStore) EndPoints() string {
	return ps.config.GetString("polly.store.endpoints")
}

//Bucket this get the type of backing store to use
func (ps *PollyStore) Bucket() string {
	return ps.config.GetString("polly.store.bucket")
}

// Version of the metadata in the store
func (ps *PollyStore) Version() (string, error) {
	versionKey := ps.root + "version"
	pair, err := ps.store.Get(versionKey)
	if err != nil {
		return "", err
	}
	return string(pair.Value), nil
}

func configRegistration() *gofig.Registration {
	r := gofig.NewRegistration("Store")
	r.Key(gofig.String, "", "polly", "", "polly.store.root")
	r.Key(gofig.String, "", "", "", "polly.store.endpoints")
	r.Key(gofig.String, "", "", "", "polly.store.bucket")
	r.Key(gofig.String, "", "", "", "polly.store.type")
	return r
}
