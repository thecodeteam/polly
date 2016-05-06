package client

import (
	// "bytes"
	"os"

	"testing"

	log "github.com/Sirupsen/logrus"
	gofig "github.com/akutz/gofig"

	"bytes"
	"github.com/akutz/goof"
	apitypes "github.com/emccode/libstorage/api/types"
	apiutils "github.com/emccode/libstorage/api/utils"
	"github.com/emccode/polly/core"
	lsclient "github.com/emccode/polly/core/libstorage/client"
	"github.com/emccode/polly/core/types"
	"github.com/stretchr/testify/assert"
	"strings"
)

// This testing package is against the libStorage client

var p *types.Polly

var defaultConfig = `
polly:
  store:
    type: boltdb
    endpoints: /tmp/boltdb
    bucket: MyBoltDb_test
  libstorage:
    host: tcp://localhost:7979
    client:
      requestPath: client
    profiles:
      enabled: true
      groups:
      - local=127.0.0.1
    server:
      endpoints:
        localhost:
          address: tcp://localhost:7979
      services:
        mock:
          libstorage:
            driver: mock
        mock2:
          libstorage:
            driver: mock
`

func TestMain(m *testing.M) {
	log.SetLevel(log.DebugLevel)

	config := gofig.New()
	configYamlBuf := []byte(defaultConfig)
	if err := config.ReadConfig(bytes.NewReader(configYamlBuf)); err != nil {
		panic(err)
	}

	p = core.NewWithConfig(config)
	if err := core.Start(p); err != nil {
		log.Error(goof.WithError("problem starting polly", err))
		os.Exit(1)
	}

	if err := p.Store.EraseStore(); err != nil {
		log.Error(goof.WithError("problem erasing store", err))
	}

	m.Run()
}

func TestNewVolume(t *testing.T) {
	avol := &apitypes.Volume{
		Name: "mock1",
		ID:   "vol-001",
	}
	vol := lsclient.NewVolume(avol, "mock")

	assert.Equal(t, "mock1", vol.Name)
	assert.Equal(t, "mock-vol-001", vol.VolumeID)
}

func TestVolumesNone(t *testing.T) {
	if err := p.Store.EraseStore(); err != nil {
		log.Error(goof.WithError("problem erasing store", err))
	}

	vols, err := p.LsClient.Volumes()
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}

	assert.Len(t, vols, 0)
}

func TestVolumeCreate(t *testing.T) {
	az := "az1"
	vtype := "type1"
	size := int64(1)
	IOPS := int64(1)

	uuid, _ := apiutils.NewUUID()
	vn := strings.Split(uuid.String(), "-")

	request := &apitypes.VolumeCreateRequest{
		Name:             vn[0],
		AvailabilityZone: &az,
		Type:             &vtype,
		Size:             &size,
		IOPS:             &IOPS,
	}
	vol, err := p.LsClient.VolumeCreate("mock", request)
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}

	assert.NotEqual(t, vol, nil)
	assert.Equal(t, vn[0], vol.Name)
}

func TestVolumeCreateRemove(t *testing.T) {
	az := "az1"
	vtype := "type1"
	size := int64(1)
	IOPS := int64(1)

	uuid, _ := apiutils.NewUUID()
	vn := strings.Split(uuid.String(), "-")

	request := &apitypes.VolumeCreateRequest{
		Name:             vn[0],
		AvailabilityZone: &az,
		Type:             &vtype,
		Size:             &size,
		IOPS:             &IOPS,
	}
	vol, err := p.LsClient.VolumeCreate("mock", request)
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}

	err = p.LsClient.VolumeRemove("mock", vol.ID)
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}

	vol, err = p.LsClient.VolumeInspect("mock", vol.ID, false)
	assert.Error(t, err)
	if err == nil {
		t.FailNow()
	}
}

func TestVolumeInspect(t *testing.T) {
	az := "az1"
	vtype := "type1"
	size := int64(1)
	IOPS := int64(1)

	uuid, _ := apiutils.NewUUID()
	vn := strings.Split(uuid.String(), "-")

	request := &apitypes.VolumeCreateRequest{
		Name:             vn[0],
		AvailabilityZone: &az,
		Type:             &vtype,
		Size:             &size,
		IOPS:             &IOPS,
	}
	vol, err := p.LsClient.VolumeCreate("mock", request)
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}

	assert.NotEqual(t, vol, nil)
	assert.Equal(t, vn[0], vol.Name)

	vol, err = p.LsClient.VolumeInspect("mock", vol.ID, false)
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}
	assert.Equal(t, request.Name, vol.Name)
}

func TestVolumes(t *testing.T) {
	if err := p.Store.EraseStore(); err != nil {
		log.Error(goof.WithError("problem erasing store", err))
	}

	az := "az1"
	vtype := "type1"
	size := int64(1)
	IOPS := int64(1)

	uuid, _ := apiutils.NewUUID()
	vn1 := strings.Split(uuid.String(), "-")

	request := &apitypes.VolumeCreateRequest{
		Name:             vn1[0],
		AvailabilityZone: &az,
		Type:             &vtype,
		Size:             &size,
		IOPS:             &IOPS,
	}
	_, err := p.LsClient.VolumeCreate("mock", request)
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}

	uuid, _ = apiutils.NewUUID()
	vn2 := strings.Split(uuid.String(), "-")

	request = &apitypes.VolumeCreateRequest{
		Name:             vn2[0],
		AvailabilityZone: &az,
		Type:             &vtype,
		Size:             &size,
		IOPS:             &IOPS,
	}
	_, err = p.LsClient.VolumeCreate("mock", request)
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}

	vols, err := p.LsClient.Volumes()
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}

	var names []string
	for _, vol := range vols {
		names = append(names, vol.Name)
	}

	assert.Contains(t, names, vn1[0])
	assert.Contains(t, names, vn2[0])
}
