package store

import (
	"errors"
	"log"
	"strings"
	"time"

	gofig "github.com/akutz/gofig"
	libkv "github.com/docker/libkv"
	store "github.com/docker/libkv/store"
	"github.com/docker/libkv/store/boltdb"
	"github.com/docker/libkv/store/consul"
	"github.com/docker/libkv/store/etcd"
	"github.com/docker/libkv/store/zookeeper"
	"github.com/emccode/libstorage/api/types"
)

const (
	//VolumeType self explanatory
	VolumeType = 1
	//SnapshotType self explanatory
	SnapshotType = 2
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

//IPollyStore a representation of the libkv Store
type IPollyStore interface {
	StoreType() string
	Root() string
	EndPoints() string
	Bucket() string

	GenerateKey(mytype int, guid string) (path string, err error)
	GetKeyFromFQKN(fqkn string) (mykey string, err error)

	SaveVolumeMetadata(volume *types.Volume) error
	GetVolumeMetadata(volume *types.Volume) error
	DeleteVolumeMetadata(volume *types.Volume) error

	SaveSnapshotMetadata(snapshot *types.Snapshot) error
	GetSnapshotMetadata(snapshot *types.Snapshot) error
	DeleteSnapshotMetadata(snapshot *types.Snapshot) error
}

//NewWithConfig This initializes new instance of this library
func NewWithConfig(config gofig.Config) (pollystore IPollyStore, err error) {
	cfg := store.Config{
		ConnectionTimeout: 10 * time.Second,
	}

	ps := new(PollyStore)
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

	return ps, nil
}

//GenerateKey generates the internal path (=key) for persisting a value
func (ps *PollyStore) GenerateKey(mytype int, guid string) (path string, err error) {
	switch mytype {
	case VolumeType:
		return ps.Root() + "volume/" + guid + "/", nil
	case SnapshotType:
		return ps.Root() + "snapshot/" + guid + "/", nil
	}

	return "", ErrObjectInvalid
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
