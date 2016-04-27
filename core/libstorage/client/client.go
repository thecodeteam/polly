package client

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/akutz/gofig"
	"github.com/akutz/goof"
	lstypes "github.com/emccode/libstorage/api/types"
	"github.com/emccode/libstorage/client"
	"github.com/emccode/polly/api/types"
)

// Client is the polly version of libstorage Client
type Client struct {
	client.Client
}

// NewWithConfig creates a new client with specified configuration object
func NewWithConfig(config gofig.Config) (*Client, error) {
	c, err := client.New(config)
	if err != nil {
		return nil, goof.WithFieldE(
			"host", config.Get("libstorage.host"),
			"error dialing libStorage service", err)
	}
	return &Client{c}, nil
}

func startClient(config gofig.Config) (*Client, error) {
	var err error
	c, err := NewWithConfig(config)
	if err != nil {
		return nil, goof.WithError("cannot connect to libstorage client", err)
	}

	// todo remove this persistent store stuff when core has it

	return c, nil
}

// NewVolume creates a Polly volume from a libStorage volume
func NewVolume(vol *lstypes.Volume, service string) *types.Volume {
	newVol := &types.Volume{
		Volume:      vol,
		ServiceName: service,
		VolumeID:    fmt.Sprintf("%s-%s", service, vol.ID),
		Labels:      make(map[string]string),
	}
	log.WithFields(log.Fields{
		"newVolume":        newVol,
		"newVolume.Volume": newVol.Volume,
	}).Debug("converted volume from libstorage to polly")
	return newVol
}

// VolumesByService returns a list of Polly volumes from libstorage
func (c Client) VolumesByService(serviceName string) ([]*types.Volume, error) {
	volumeMap, err := c.Client.VolumesByService(serviceName, false)
	if err != nil {
		return nil, err
	}

	var vols []*types.Volume
	for _, vol := range volumeMap {
		vols = append(vols, NewVolume(vol, serviceName))
	}
	return vols, nil
}

// Volumes returns a list of Polly volumes from libstorage
func (c Client) Volumes() ([]*types.Volume, error) {
	serviceVolumeMap, err := c.Client.Volumes(false)
	if err != nil {
		return nil, err
	}

	var vols []*types.Volume
	for serviceName, volumeMap := range serviceVolumeMap {
		for _, vol := range volumeMap {
			vols = append(vols, NewVolume(vol, serviceName))
		}
	}
	return vols, nil
}

// VolumeInspect returns a Polly volume
func (c Client) VolumeInspect(serviceName, volumeID string, attachments bool) (*types.Volume, error) {
	vol, err := c.Client.VolumeInspect(serviceName, volumeID, attachments)
	if err != nil {
		return nil, err
	}

	return NewVolume(vol, serviceName), nil
}
