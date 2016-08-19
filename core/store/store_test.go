package store

import (
	"bytes"
	"os"
	"strconv"
	"testing"

	log "github.com/Sirupsen/logrus"

	gofig "github.com/akutz/gofig"
	lstypes "github.com/emccode/libstorage/api/types"
	"github.com/emccode/polly/api/types"
	lsclient "github.com/emccode/polly/core/libstorage/client"
	"github.com/stretchr/testify/assert"
)

const (
	libStorageConfigBaseTestBolt = `
polly:
  store:
    type: boltdb
    endpoints: /tmp/boltdb-test
    bucket: MyBoltDb_test
  server:
    services:
      vfs:
        libstorage:
          storage:
            driver: vfs
`
	libStorageConfigBaseBenchBolt = `
polly:
  store:
    type: boltdb
    endpoints: /tmp/boltdb-bench
    bucket: MyBoltDb_bench
  server:
    services:
      vfs:
        libstorage:
          storage:
            driver: vfs
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
	vol, _ := lsclient.NewVolume(nil, lsvol, service)
	return vol
}

func TestMain(m *testing.M) {
	log.SetLevel(log.ErrorLevel)
	os.Setenv("POLLY_DEBUG", "true")
	config := gofig.New()

	configYamlBuf := []byte(libStorageConfigBaseTestBolt)
	if err := config.ReadConfig(bytes.NewReader(configYamlBuf)); err != nil {
		panic(err)
	}

	os.Remove("/tmp/boltdb-test")
	os.Remove("/tmp/boltdb-bench")

	var err error
	ps, err = NewWithConfig(config)
	if err != nil {
		log.WithError(err).Fatal("Failed to create PollyStore")
	}

	m.Run()
}

//Testing
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

	configYamlBuf := []byte(libStorageConfigBaseTestBolt)
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

//Benchmarks
func BenchmarkGetVolumeIDNotExist(b *testing.B) {
	myConfig := gofig.New()

	configYamlBuf := []byte(libStorageConfigBaseBenchBolt)
	if err := myConfig.ReadConfig(bytes.NewReader(configYamlBuf)); err != nil {
		log.WithError(err).Fatal("Failed to create PollyStore")
	}

	psBench, err := NewWithConfig(myConfig.Scope("polly.store"))
	if err != nil {
		log.Fatal("Failed to create PollyStore")
	}

	err = psBench.EraseStore()
	if err != nil {
		log.WithError(err).Fatal("Failed to reset PollyStore")
	}

	volume := newVolume("pollytestpkg2", "testiddoesntexist")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			psBench.Exists(volume)
		}
	})
}

func BenchmarkGetVolumeIDs(b *testing.B) {
	myConfig := gofig.New()

	configYamlBuf := []byte(libStorageConfigBaseBenchBolt)
	if err := myConfig.ReadConfig(bytes.NewReader(configYamlBuf)); err != nil {
		log.WithError(err).Fatal("Failed to create PollyStore")
	}

	psBench, err := NewWithConfig(myConfig.Scope("polly.store"))
	if err != nil {
		log.Fatal("Failed to create PollyStore")
	}

	err = psBench.EraseStore()
	if err != nil {
		log.WithError(err).Fatal("Failed to reset PollyStore")
	}

	volume := newVolume("pollytestpkg1", "testid1")

	err = psBench.SaveVolumeMetadata(volume)
	if err != nil {
		log.WithError(err).Fatal("Failed to SaveVolumeMetadata")
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			psBench.GetVolumeIds()
		}
	})
}

func BenchmarkNewVolumeMetadata(b *testing.B) {
	myConfig := gofig.New()

	configYamlBuf := []byte(libStorageConfigBaseBenchBolt)
	if err := myConfig.ReadConfig(bytes.NewReader(configYamlBuf)); err != nil {
		log.WithError(err).Fatal("Failed to create PollyStore")
	}

	psBench, err := NewWithConfig(myConfig.Scope("polly.store"))
	if err != nil {
		log.Fatal("Failed to create PollyStore")
	}

	err = psBench.EraseStore()
	if err != nil {
		log.WithError(err).Fatal("Failed to reset PollyStore")
	}

	volume := newVolume("pollytestpkg1", "testid1")
	volume.Schedulers = []string{"testScheduler"}

	err = psBench.SaveVolumeMetadata(volume)
	if err != nil {
		log.WithError(err).Fatal("Failed to SaveVolumeMetadata")
	}

	volume = newVolume("pollytestpkg1", "testid1")
	_, err = psBench.SetVolumeMetadata(volume)
	if err != nil {
		log.WithError(err).Fatal("Failed to SaveVolumeMetadata")
	}

	volume.Labels = make(map[string]string)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		volume.Labels["testkey"+strconv.Itoa(i)] = "testval"
		psBench.SaveVolumeMetadata(volume)
	}
}

func BenchmarkUpdateVolumeMetadata(b *testing.B) {
	myConfig := gofig.New()

	configYamlBuf := []byte(libStorageConfigBaseBenchBolt)
	if err := myConfig.ReadConfig(bytes.NewReader(configYamlBuf)); err != nil {
		log.WithError(err).Fatal("Failed to create PollyStore")
	}

	psBench, err := NewWithConfig(myConfig.Scope("polly.store"))
	if err != nil {
		log.Fatal("Failed to create PollyStore")
	}

	err = psBench.EraseStore()
	if err != nil {
		log.WithError(err).Fatal("Failed to reset PollyStore")
	}

	volume := newVolume("pollytestpkg1", "testid1")
	volume.Schedulers = []string{"testScheduler"}

	err = psBench.SaveVolumeMetadata(volume)
	if err != nil {
		log.WithError(err).Fatal("Failed to SaveVolumeMetadata")
	}

	volume = newVolume("pollytestpkg1", "testid1")
	_, err = psBench.SetVolumeMetadata(volume)
	if err != nil {
		log.WithError(err).Fatal("Failed to SaveVolumeMetadata")
	}

	volume.Labels = make(map[string]string)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		volume.Labels["testkey1"] = "testval" + strconv.Itoa(i)
		psBench.SaveVolumeMetadata(volume)
	}
}
