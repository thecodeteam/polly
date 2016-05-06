package store

import (
	"bytes"

	log "github.com/Sirupsen/logrus"

	"testing"

	"os"

	gofig "github.com/akutz/gofig"
	lstypes "github.com/emccode/libstorage/api/types"
	"github.com/emccode/polly/api/types"
	lsclient "github.com/emccode/polly/core/libstorage/client"
	"github.com/stretchr/testify/assert"
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

var ps *PollyStore

func newVolume(service, volumeID string) *types.Volume {
	lsvol := &lstypes.Volume{
		ID: volumeID,
	}
	return lsclient.NewVolume(lsvol, service)
}

func TestMain(m *testing.M) {
	log.SetLevel(log.DebugLevel)
	os.Setenv("POLLY_DEBUG", "true")
	config := gofig.New()

	configYamlBuf := []byte(libStorageConfigBaseBolt)
	if err := config.ReadConfig(bytes.NewReader(configYamlBuf)); err != nil {
		panic(err)
	}

	var err error
	ps, err = NewWithConfig(config.Scope("polly.store"))
	if err != nil {
		log.Fatal("Failed to create PollyStore")
	}

	vol := newVolume("pollytestpkg1", "testid1")
	_ = ps.RemoveVolumeMetadata(vol)

	if err := ps.EraseStore(); err != nil {
		log.Fatal("Could not clear polly store")
	}

	m.Run()
}

func TestGenerateRootKey(t *testing.T) {
	key, err := ps.GenerateRootKey(VolumeInternalLabelsType)
	assert.NoError(t, err)
	assert.Equal(t, "polly/volumeinternallabels/", key)
}

func TestGenerateObjectKey(t *testing.T) {
	key, err := ps.GenerateObjectKey(VolumeInternalLabelsType, "pollytestpkg1-testid1")
	assert.NoError(t, err)
	assert.Equal(t, "polly/volumeinternallabels/pollytestpkg1-testid1/", key)
}

func TestGenerateObjectKeyInvalid(t *testing.T) {
	_, err := ps.GenerateObjectKey(VolumeInternalLabelsType, "")
	assert.Error(t, err)
}

func TestVersionOfStore(t *testing.T) {
	version, err := ps.Version()
	assert.NoError(t, err)
	assert.Equal(t, version, "v0.1.0")
}

func TestNotExist(t *testing.T) {
	volume := newVolume("pollytestpkg2", "testiddoesntexist")

	exists, err := ps.Exists(volume)
	assert.NoError(t, err)
	assert.Equal(t, false, exists)
}

func TestGetVolumeIDs(t *testing.T) {
	volume := newVolume("pollytestpkg1", "testid1")

	err := ps.SaveVolumeMetadata(volume)
	assert.NoError(t, err)

	ids, err := ps.GetVolumeIds()
	assert.NoError(t, err)
	assert.Equal(t, len(ids), 1)
}

func TestRemovingSchedulers(t *testing.T) {
	volume := newVolume("pollytestpkg1", "testid1")
	volume.Schedulers = []string{"testScheduler"}

	err := ps.SaveVolumeMetadata(volume)
	assert.NoError(t, err)

	volume.Schedulers = nil
	err = ps.SaveVolumeMetadata(volume)
	assert.NoError(t, err)

	volume = newVolume("pollytestpkg1", "testid1")
	_, err = ps.SetVolumeMetadata(volume)
	assert.NoError(t, err)
	assert.Len(t, volume.Schedulers, 0)
}

func TestSaveVolumeMetadata(t *testing.T) {
	volume := newVolume("pollytestpkg1", "testid1")
	volume.Schedulers = []string{"testScheduler"}

	err := ps.SaveVolumeMetadata(volume)
	assert.NoError(t, err)

	volume = newVolume("pollytestpkg1", "testid1")
	_, err = ps.SetVolumeMetadata(volume)
	assert.NoError(t, err)
	assert.Contains(t, volume.Schedulers, "testScheduler")

	volume.Labels = make(map[string]string)
	volume.Labels["testkey1"] = "testval1"
	err = ps.SaveVolumeMetadata(volume)
	assert.NoError(t, err)

	volume = newVolume("pollytestpkg1", "testid1")
	_, err = ps.SetVolumeMetadata(volume)
	assert.NoError(t, err)
	assert.Equal(t, "testval1", volume.Labels["testkey1"])
}

func TestSaveVolumeMetadataNoChanges(t *testing.T) {
	volume := newVolume("pollytestpkg1", "testid1")
	volume.Schedulers = []string{"testScheduler"}

	err := ps.SaveVolumeMetadata(volume)
	assert.NoError(t, err)

	volume = newVolume("pollytestpkg1", "testid1")
	_, err = ps.SetVolumeMetadata(volume)
	assert.NoError(t, err)
	assert.Contains(t, volume.Schedulers, "testScheduler")

	err = ps.SaveVolumeMetadata(volume)
	assert.NoError(t, err)

	volume = newVolume("pollytestpkg1", "testid1")
	_, err = ps.SetVolumeMetadata(volume)
	assert.NoError(t, err)
	assert.Contains(t, volume.Schedulers, "testScheduler")
}

func TestUpdateVolumeMetadata(t *testing.T) {
	volume := newVolume("pollytestpkg1", "testid1")
	volume.Schedulers = []string{"testScheduler"}

	err := ps.SaveVolumeMetadata(volume)
	assert.NoError(t, err)

	volume = newVolume("pollytestpkg1", "testid1")
	_, err = ps.SetVolumeMetadata(volume)
	assert.NoError(t, err)
	assert.Contains(t, volume.Schedulers, "testScheduler")

	volume = newVolume("pollytestpkg1", "testid1")
	_, err = ps.SetVolumeMetadata(volume)
	assert.NoError(t, err)
	volume.Schedulers = []string{"testScheduler2"}

	err = ps.SaveVolumeMetadata(volume)
	assert.NoError(t, err)

	volume = newVolume("pollytestpkg1", "testid1")
	_, err = ps.SetVolumeMetadata(volume)
	assert.NoError(t, err)
	assert.Contains(t, volume.Schedulers, "testScheduler2")
}

func TestRemoveVolumeMetadata(t *testing.T) {
	volume := newVolume("pollytestpkg1", "testid1")
	volume.Schedulers = []string{"testScheduler"}

	err := ps.SaveVolumeMetadata(volume)
	assert.NoError(t, err)

	volume = newVolume("pollytestpkg1", "testid1")
	_, err = ps.SetVolumeMetadata(volume)
	assert.NoError(t, err)
	assert.Contains(t, volume.Schedulers, "testScheduler")

	volume = newVolume("pollytestpkg1", "testid1")
	_, err = ps.SetVolumeMetadata(volume)
	assert.NoError(t, err)

	volume = newVolume("pollytestpkg1", "testid1")
	err = ps.RemoveVolumeMetadata(volume)
	assert.NoError(t, err)

	volume = newVolume("pollytestpkg1", "testid1")
	_, err = ps.SetVolumeMetadata(volume)
	assert.NoError(t, err)
	assert.Len(t, volume.Schedulers, 0)
}

func TestEraseStore(t *testing.T) {
	myConfig := gofig.New()

	configYamlBuf := []byte(libStorageConfigBaseBolt)
	if err := myConfig.ReadConfig(bytes.NewReader(configYamlBuf)); err != nil {
		panic(err)
	}

	var err error
	myPs, err := NewWithConfig(myConfig.Scope("polly.store"))
	if err != nil {
		log.Fatal("Failed to create PollyStore")
	}

	volume := newVolume("pollytestpkg1", "testid1")

	err = myPs.SaveVolumeMetadata(volume)
	assert.NoError(t, err)

	err = myPs.EraseStore()
	assert.NoError(t, err)

	ids, err := myPs.GetVolumeIds()
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}

	assert.Len(t, ids, 0)
}
