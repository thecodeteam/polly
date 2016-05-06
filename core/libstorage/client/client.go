package client

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/akutz/gofig"
	"github.com/akutz/goof"
	apitypes "github.com/emccode/libstorage/api/types"
	"github.com/emccode/libstorage/client"
	"github.com/emccode/polly/api/types"
)

// Client is the polly version of libstorage Client
type Client struct {
	apitypes.Client
	ctx    apitypes.Context
	config gofig.Config
}

// NewWithConfig creates a new client with specified configuration object
func NewWithConfig(ctx apitypes.Context, config gofig.Config) (*Client, error) {
	c, err := client.New(config)
	if err != nil {
		return nil, goof.WithFieldE(
			"host", config.Get("libstorage.host"),
			"error dialing libStorage service", err)
	}

	return &Client{c, ctx, config}, nil
}

// NewVolume creates a Polly volume from a libStorage volume
func NewVolume(vol *apitypes.Volume, service string) *types.Volume {
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
	c.Client.API().AddHeader("requestPath", c.requestPath())
	volumeMap, err := c.Client.API().VolumesByService(
		c.ctx, serviceName, false)
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
	c.Client.API().AddHeader("requestPath", c.requestPath())
	serviceVolumeMap, err := c.Client.API().Volumes(
		c.ctx, false)
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
	c.Client.API().AddHeader("requestPath", c.requestPath())
	vol, err := c.Client.API().VolumeInspect(c.ctx, serviceName, volumeID, attachments)
	if err != nil {
		return nil, err
	}

	return NewVolume(vol, serviceName), nil
}

// VolumeCreate creates a Polly Volume
func (c *Client) VolumeCreate(serviceName string, request *apitypes.VolumeCreateRequest) (*types.Volume, error) {
	c.Client.API().AddHeader("requestPath", "admin")
	vol, err := c.Client.API().VolumeCreate(c.ctx, serviceName, request)
	if err != nil {
		return nil, err
	}

	return NewVolume(vol, serviceName), nil
}

// VolumeRemove removes a Polly Volume
func (c *Client) VolumeRemove(serviceName string, volumeID string) error {
	c.Client.API().AddHeader("requestPath", c.requestPath())
	err := c.Client.API().VolumeRemove(c.ctx, serviceName, volumeID)
	if err != nil {
		return err
	}

	return nil
}

func (c Client) requestPath() string {
	return c.config.GetString("libstorage.client.requestPath")
}
