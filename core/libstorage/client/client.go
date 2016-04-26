package client

import (
	// "bytes"
	// "fmt"
	// log "github.com/Sirupsen/logrus"
	"github.com/akutz/gofig"
	"github.com/akutz/goof"

	"github.com/emccode/libstorage/client"
	"github.com/emccode/polly/api/types"
	// "github.com/emccode/polly/core/store"
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

// VolumesByService returns a list of Polly volumes from libstorage
func (c Client) VolumesByService(name string) ([]*types.Volume, error) {
	volumeMap, err := c.Client.VolumesByService(name, false)
	if err != nil {
		return nil, err
	}

	var vols []*types.Volume
	for _, vol := range volumeMap {
		vols = append(vols, &types.Volume{
			Volume:      vol,
			ServiceName: name,
			Labels:      make(map[string]string),
		})
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
			vols = append(vols, &types.Volume{
				Volume:      vol,
				ServiceName: serviceName,
				Labels:      make(map[string]string),
			})
		}
	}
	return vols, nil
}
