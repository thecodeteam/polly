package store

import (
	"errors"
	"log"
	"strings"
	"time"

	gofig "github.com/akutz/gofig"
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
		key, err := ps.GenerateKey(category, "")
		if err != nil {
			log.Fatal("failed key gen on store init", err)
			return nil, err
		}
		ps.store.Put(key, []byte(""), nil)
	}

	return ps, nil
}

//GenerateKey generates the internal path (=key) for persisting a value
func (ps *PollyStore) GenerateKey(mytype int, guid string) (path string, err error) {
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
	path = store.Normalize(path)
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
