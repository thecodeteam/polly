package client

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/akutz/gofig"
	"github.com/akutz/goof"
	"github.com/emccode/libstorage"
	apitypes "github.com/emccode/libstorage/api/types"
	pcontext "github.com/emccode/polly/api/context"
	"github.com/emccode/polly/api/types"
	"strings"
)

// Client is the polly version of libstorage Client
type Client struct {
	apitypes.Client
	ctx            apitypes.Context
	config         gofig.Config
	Services       apitypes.ServicesMap
	ServiceDrivers map[string]string
	DriverService  map[string]string
}

// NewWithConfig creates a new client with specified configuration object
func NewWithConfig(ctx apitypes.Context, config gofig.Config) (*Client, error) {
	config = config.Scope("polly")
	_, err, errs := libstorage.Serve(config)
	if err != nil {
		return nil, goof.WithError(
			"error starting libstorage server and client", err)
	}
	go func() {
		err := <-errs
		if err != nil {
			log.Error(err)
		}
	}()

	c, err := libstorage.Dial(config)
	if err != nil {
		return nil, goof.WithFieldE(
			"host", config.Get("libstorage.host"),
			"error dialing libStorage service", err)
	}

	services, err := c.API().Services(ctx)
	if err != nil {
		return nil, goof.WithError("cannot instantiate client services", err)
	}
	for _, service := range services {
		if strings.Contains(service.Name, "-") {
			return nil, goof.New("illegal character in serviceName '-'")
		}
	}

	serviceDrivers := make(map[string]string)
	for _, service := range services {
		serviceDrivers[service.Name] = service.Driver.Name
	}

	driverService := make(map[string]string)
	for _, s := range services {
		driverService[s.Driver.Name] = s.Name
	}

	return &Client{c, ctx, config, services, serviceDrivers, driverService}, nil
}

func getDriver(c *Client, s string) (string, error) {
	if service, ok := c.Services[s]; ok {
		return service.Driver.Name, nil
	}

	return "", goof.New("no service found by name")
}

// NewVolume creates a Polly volume from a libStorage volume
func NewVolume(c *Client, vol *apitypes.Volume, service string) (*types.Volume, error) {
	var d string
	var err error
	if c != nil {
		d, err = getDriver(c, service)
		if err != nil {
			return nil, err
		}
	} else {
		d = service
	}

	newVol := &types.Volume{
		Volume:      vol,
		ServiceName: service,
		VolumeID:    fmt.Sprintf("%s-%s", d, vol.ID),
		Labels:      make(map[string]string),
	}
	log.WithFields(log.Fields{
		"newVolume":        newVol,
		"newVolume.Volume": newVol.Volume,
	}).Debug("converted volume from libstorage to polly")
	return newVol, nil
}

// VolumesByService returns a list of Polly volumes from libstorage
func (c *Client) VolumesByService(serviceName string) ([]*types.Volume, error) {
	if c.ctx.Value(pcontext.RequestPathHeaderKey) == nil {
		c.ctx = c.ctx.WithValue(pcontext.RequestPathHeaderKey, "admin")
	}
	volumeMap, err := c.Client.API().VolumesByService(
		c.ctx, serviceName, false)
	if err != nil {
		return nil, err
	}

	var vols []*types.Volume
	for _, vol := range volumeMap {
		nv, err := NewVolume(c, vol, serviceName)
		if err != nil {
			return nil, err
		}
		vols = append(vols, nv)
	}
	return vols, nil
}

// Volumes returns a list of Polly volumes from libstorage
func (c *Client) Volumes() ([]*types.Volume, error) {
	if c.ctx.Value(pcontext.RequestPathHeaderKey) == nil {
		c.ctx = c.ctx.WithValue(pcontext.RequestPathHeaderKey, "admin")
	}
	serviceVolumeMap, err := c.Client.API().Volumes(
		c.ctx, false)
	if err != nil {
		return nil, err
	}

	var vols []*types.Volume
	for serviceName, volumeMap := range serviceVolumeMap {
		for _, vol := range volumeMap {
			nv, err := NewVolume(c, vol, serviceName)
			if err != nil {
				return nil, err
			}
			vols = append(vols, nv)
		}
	}
	return vols, nil
}

// VolumeInspect returns a Polly volume
func (c *Client) VolumeInspect(serviceName, volumeID string, attachments bool) (*types.Volume, error) {
	if c.ctx.Value(pcontext.RequestPathHeaderKey) == nil {
		c.ctx = c.ctx.WithValue(pcontext.RequestPathHeaderKey, "admin")
	}

	vol, err := c.Client.API().VolumeInspect(c.ctx, serviceName, volumeID, attachments)
	if err != nil {
		return nil, err
	}

	nv, err := NewVolume(c, vol, serviceName)
	if err != nil {
		return nil, err
	}

	return nv, nil
}

// VolumeCreate creates a Polly Volume
func (c *Client) VolumeCreate(serviceName string, request *apitypes.VolumeCreateRequest) (*types.Volume, error) {
	if c.ctx.Value(pcontext.RequestPathHeaderKey) == nil {
		c.ctx = c.ctx.WithValue(pcontext.RequestPathHeaderKey, "admin")
	}
	vol, err := c.Client.API().VolumeCreate(c.ctx, serviceName, request)
	if err != nil {
		return nil, err
	}

	nv, err := NewVolume(c, vol, serviceName)
	if err != nil {
		return nil, err
	}
	return nv, nil
}

// VolumeRemove removes a Polly Volume
func (c *Client) VolumeRemove(serviceName string, volumeID string) error {
	if c.ctx.Value(pcontext.RequestPathHeaderKey) == nil {
		c.ctx = c.ctx.WithValue(pcontext.RequestPathHeaderKey, "admin")
	}
	err := c.Client.API().VolumeRemove(c.ctx, serviceName, volumeID)
	if err != nil {
		return err
	}

	return nil
}

func (c Client) requestPath() string {
	return c.config.GetString("libstorage.client.requestPath")
}
