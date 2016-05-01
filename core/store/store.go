package store

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"os"
	"strings"
	"time"

	gofig "github.com/akutz/gofig"
	"github.com/akutz/goof"
	"github.com/docker/libkv"
	store "github.com/docker/libkv/store"
	"github.com/docker/libkv/store/boltdb"
	"github.com/docker/libkv/store/consul"
	"github.com/docker/libkv/store/etcd"
	"github.com/docker/libkv/store/zookeeper"
)

const (
	//VolumeType is used to identify metadata for the libstorage layer
	VolumeType = 1
	//SnapshotType self explanatory
	SnapshotType = 2
	//VolumeInternalLabelsType is used to identify metadata for the Polly admin layer
	VolumeInternalLabelsType = 3
	//VolumeAdminLabelsType is used to identify labels for the Polly admin layer
	VolumeAdminLabelsType = 4
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

func configRegistration() *gofig.Registration {
	r := gofig.NewRegistration("Store")
	r.Key(gofig.String, "", "polly", "", "root", "root")
	return r
}

//NewWithConfig This initializes new instance of this library
func NewWithConfig(config gofig.Config) (pollystore *PollyStore, err error) {
	cfg := store.Config{
		ConnectionTimeout: 10 * time.Second,
	}

	ps := &PollyStore{}
	ps.config = config
	ps.root = "polly/"

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

	initlist := []int{VolumeType, VolumeInternalLabelsType, VolumeAdminLabelsType}
	for _, category := range initlist {
		key, err := ps.GenerateRootKey(category)
		if err != nil {
			log.WithError(err).Fatal("failed key gen on store init")
			return nil, err
		}
		ps.Put(key, []byte(""))
	}

	return ps, nil
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
		parts = append(parts, "volumelibstorage")
	case SnapshotType:
		parts = append(parts, "snapshot")
	case VolumeInternalLabelsType:
		parts = append(parts, "volumeinternal")
	case VolumeAdminLabelsType:
		parts = append(parts, "volumeadmin")
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
	if err := ps.store.DeleteTree("/volumeinternal"); err != nil {
		return err
	}
	if err := ps.store.DeleteTree("/volumesnapshot"); err != nil {
		return err
	}
	if err := ps.store.DeleteTree("/volumeadmin"); err != nil {
		return err
	}
	if err := ps.store.DeleteTree("/volume"); err != nil {
		return err
	}
	return nil
}

//StoreType this generates the type of backing store to use
func (ps *PollyStore) StoreType() string {
	return ps.config.GetString("type")
}

//Root this get the type of backing store to use
func (ps *PollyStore) Root() string {
	return ps.config.GetString("root")
}

//EndPoints this gets the endpoints of the store
func (ps *PollyStore) EndPoints() string {
	return ps.config.GetString("endpoints")
}

//Bucket this get the type of backing store to use
func (ps *PollyStore) Bucket() string {
	return ps.config.GetString("bucket")
}
