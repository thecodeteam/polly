package store

import (
	"bytes"
	log "github.com/Sirupsen/logrus"

	"testing"

	gofig "github.com/akutz/gofig"
	lstypes "github.com/emccode/libstorage/api/types"
	"github.com/emccode/polly/api/types"
	lsclient "github.com/emccode/polly/core/libstorage/client"
	"github.com/stretchr/testify/assert"
	"os"
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

	m.Run()
}

func TestGenerateRootKey(t *testing.T) {
	key, err := ps.GenerateRootKey(VolumeInternalLabelsType)
	assert.NoError(t, err)
	assert.Equal(t, "/volumeinternal/", key)
}

func TestGenerateObjectKey(t *testing.T) {
	key, err := ps.GenerateObjectKey(VolumeInternalLabelsType, "pollytestpkg1-testid1")
	assert.NoError(t, err)
	assert.Equal(t, "/volumeinternal/pollytestpkg1-testid1/", key)
}

func TestGenerateObjectKeyInvalid(t *testing.T) {
	_, err := ps.GenerateObjectKey(VolumeInternalLabelsType, "")
	assert.Error(t, err)
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
