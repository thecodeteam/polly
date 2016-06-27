package client

import (
	// "bytes"
	"os"

	"testing"

	log "github.com/Sirupsen/logrus"
	// gofig "github.com/akutz/gofig"
	"fmt"

	"github.com/akutz/goof"
	config "github.com/emccode/polly/core/config"
	"github.com/emccode/polly/core/store"
	"github.com/emccode/polly/daemon"
	"github.com/stretchr/testify/assert"
)

var tpc Client

var defaultConfig = `
polly:
  host: tcp://127.0.0.1:7978
  store:
    type: boltdb
    endpoints: /tmp/boltdb
    bucket: MyBoltDb_test
libstorage:
  host: tcp://localhost:7981
  server:
    endpoints:
      localhost:
        address: tcp://localhost:7981
    services:
      mockService:
        libstorage:
          driver: mock
`

func TestMain(m *testing.M) {
	log.SetLevel(log.DebugLevel)
	cfg, err := config.NewWithConfig(defaultConfig)
	if err != nil {
		log.Error(goof.WithError("problem getting config", err))
		os.Exit(1)
	}

	scfg, _ := cfg.Copy()
	ps, err := store.NewWithConfig(scfg.Scope("polly.store"))
	if err != nil {
		log.Error(goof.WithError("problem initialization store", err))
		os.Exit(1)
	}

	if err = ps.EraseStore(); err != nil {
		log.Error(goof.WithError("problem erasing store", err))
	}

	init := make(chan error)
	stop := make(chan os.Signal)

	go func() {
		daemon.Start(cfg, init, stop)
	}()

	i := <-init
	if i != nil {
		log.Fatal(goof.WithError("got init error", err))
		os.Exit(1)
	}

	tpc, err = New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	m.Run()
}

func TestPolly(t *testing.T) {
	assert.NotEqual(t, tpc, nil)
}

func TestVolumesAll(t *testing.T) {
	vols, err := tpc.VolumesAll()

	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}
	assert.Len(t, vols, 3)
}

func TestVolumes(t *testing.T) {
	vols, err := tpc.Volumes()
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}
	assert.Len(t, vols, 0)
}

func TestVolumeInspect(t *testing.T) {
	vol, err := tpc.VolumeInspect("mock-vol-000")
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}
	assert.Equal(t, "mock-vol-000", vol.VolumeID)
	assert.Equal(t, "mockservice", vol.ServiceName)
}

func TestVolumeOffer(t *testing.T) {
	vol, err := tpc.VolumeOffer("mock-vol-001", []string{"mesos"})
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}

	assert.Equal(t, "mock-vol-001", vol.VolumeID)
	assert.Equal(t, "mockservice", vol.ServiceName)
	assert.Contains(t, vol.Schedulers, "mesos")

	vol, err = tpc.VolumeInspect("mock-vol-001")
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}

	assert.Equal(t, "mock-vol-001", vol.VolumeID)
	assert.Contains(t, vol.Schedulers, "mesos")
}

func TestVolumeOfferMultiple(t *testing.T) {
	vol, err := tpc.VolumeOffer("mock-vol-001", []string{"mesos", "docker"})
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}

	assert.Equal(t, "mock-vol-001", vol.VolumeID)
	assert.Equal(t, "mockservice", vol.ServiceName)
	assert.Contains(t, vol.Schedulers, "mesos")
	assert.Contains(t, vol.Schedulers, "docker")

	vol, err = tpc.VolumeInspect("mock-vol-001")
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}
	assert.Equal(t, "mock-vol-001", vol.VolumeID)
	assert.Contains(t, vol.Schedulers, "mesos")
	assert.Contains(t, vol.Schedulers, "docker")
}

func TestVolumeOfferRevoke(t *testing.T) {
	vol, err := tpc.VolumeOffer("mock-vol-001", []string{"mesos", "docker"})
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}

	assert.Equal(t, "mock-vol-001", vol.VolumeID)
	assert.Equal(t, "mockservice", vol.ServiceName)
	assert.Contains(t, vol.Schedulers, "mesos")
	assert.Contains(t, vol.Schedulers, "docker")

	vol, err = tpc.VolumeOfferRevoke("mock-vol-001", []string{"mesos", "docker"})
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}

	assert.Equal(t, "mock-vol-001", vol.VolumeID)
	assert.Equal(t, "mockservice", vol.ServiceName)
	assert.Len(t, vol.Schedulers, 0)

	vol, err = tpc.VolumeInspect("mock-vol-001")
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}
	assert.Len(t, vol.Schedulers, 0)
}

func TestVolumeLabel(t *testing.T) {
	vol, err := tpc.VolumeInspect("mock-vol-000")
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}
	assert.Equal(t, "mock-vol-000", vol.VolumeID)
	assert.Equal(t, "mockservice", vol.ServiceName)

	key1 := "key1"
	value1 := "value1"
	key2 := "key2"
	value2 := "value2"
	labels := []string{fmt.Sprintf("%s=%s", key1, value1), fmt.Sprintf("%s=%s", key2, value2)}
	vol, err = tpc.VolumeLabel("mock-vol-000", labels)
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}

	assert.Equal(t, "mock-vol-000", vol.VolumeID)
	assert.Equal(t, "mockservice", vol.ServiceName)
	assert.Len(t, vol.Labels, 2)

	vol, err = tpc.VolumeInspect("mock-vol-000")
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}

	_, ok := vol.Labels[key1]
	assert.Equal(t, true, ok)
	assert.Equal(t, value1, vol.Labels[key1])
	_, ok = vol.Labels[key2]
	assert.Equal(t, true, ok)
	assert.Equal(t, value2, vol.Labels[key2])
	assert.Len(t, vol.Labels, 2)

}

func TestVolumeLabelsRemove(t *testing.T) {
	TestVolumeLabel(t)

	key1 := "key1"
	value1 := "value1"
	key2 := "key2"
	value2 := "value2"
	vol, err := tpc.VolumeInspect("mock-vol-000")
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}
	_, ok := vol.Labels[key1]
	assert.Equal(t, true, ok)
	assert.Equal(t, value1, vol.Labels[key1])

	vol, err = tpc.VolumeLabelsRemove("mock-vol-000", []string{key1})
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}
	_, ok = vol.Labels[key1]
	assert.Equal(t, false, ok)

	vol, err = tpc.VolumeInspect("mock-vol-000")
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}
	_, ok = vol.Labels[key1]
	assert.Equal(t, false, ok)
	_, ok = vol.Labels[key2]
	assert.Equal(t, true, ok)
	assert.Equal(t, value2, vol.Labels[key2])
	assert.Len(t, vol.Labels, 1)

}

func TestVolumeCreate(t *testing.T) {

	availabilityZone := "az1"
	IOPS := int64(1000)
	name := "NewVolume"
	size := int64(1000)
	vtype := "silver"
	key1 := "key1"
	value1 := "value1"
	key2 := "key2"
	value2 := "value2"
	schedulers := []string{"scheduler1", "scheduler2"}
	labels := []string{fmt.Sprintf("%s=%s", key1, value1), fmt.Sprintf("%s=%s", key2, value2)}
	fields := []string{fmt.Sprintf("%s=%s", key1, value1), fmt.Sprintf("%s=%s", key2, value2)}
	service := "mockservice"
	driver := "mock"
	vol, err := tpc.VolumeCreate(service, name, vtype, size, IOPS, availabilityZone, schedulers, labels, fields)
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}
	_, ok := vol.Labels[key1]
	assert.Equal(t, true, ok)
	assert.Equal(t, value1, vol.Labels[key1])
	_, ok = vol.Labels[key2]
	assert.Equal(t, true, ok)
	assert.Equal(t, value2, vol.Labels[key2])
	assert.Contains(t, vol.Schedulers, "scheduler1")
	assert.Contains(t, vol.Schedulers, "scheduler2")
	assert.Equal(t, availabilityZone, vol.Volume.AvailabilityZone)
	assert.Equal(t, IOPS, vol.Volume.IOPS)
	assert.Equal(t, name, vol.Volume.Name)
	assert.Equal(t, size, vol.Volume.Size)
	assert.Equal(t, vtype, vol.Volume.Type)
	assert.Equal(t, fmt.Sprintf("%s-%s", driver, "vol-004"), vol.VolumeID)

	vol, err = tpc.VolumeInspect(fmt.Sprintf("%s-%s", driver, "vol-004"))
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}
	_, ok = vol.Labels[key1]
	assert.Equal(t, true, ok)
	assert.Equal(t, value1, vol.Labels[key1])
	_, ok = vol.Labels[key2]
	assert.Equal(t, true, ok)
	assert.Equal(t, value2, vol.Labels[key2])
	assert.Contains(t, vol.Schedulers, "scheduler1")
	assert.Contains(t, vol.Schedulers, "scheduler2")
	assert.Equal(t, availabilityZone, vol.Volume.AvailabilityZone)
	assert.Equal(t, IOPS, vol.Volume.IOPS)
	assert.Equal(t, name, vol.Volume.Name)
	assert.Equal(t, size, vol.Volume.Size)
	assert.Equal(t, vtype, vol.Volume.Type)
	assert.Equal(t, fmt.Sprintf("%s-%s", driver, "vol-004"), vol.VolumeID)

}

func TestVolumeRemove(t *testing.T) {
	vol, err := tpc.VolumeInspect(fmt.Sprintf("%s-%s", "mock", "vol-001"))
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}

	err = tpc.VolumeRemove(vol.VolumeID)
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}

	_, err = tpc.VolumeInspect(vol.VolumeID)
	assert.Error(t, err)
	if err == nil {
		t.FailNow()
	}
}
