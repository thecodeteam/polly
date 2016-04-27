package client

import (
	"fmt"

	"github.com/emccode/polly/api/types"
)

// Volumes returns a list of all registered Volumes for all Services.
func (c *Client) Volumes() (reply []*types.Volume, err error) {
	url := fmt.Sprintf("/admin/volumes")
	if _, err = c.httpGet(url, &reply); err != nil {
		return nil, err
	}
	return reply, nil
}

// VolumesAll returns a list of all Volumes available for all Services.
func (c *Client) VolumesAll() (reply []*types.Volume, err error) {
	url := fmt.Sprintf("/admin/volumesall")
	if _, err = c.httpGet(url, &reply); err != nil {
		return nil, err
	}
	return reply, nil
}

// VolumeInspect will inspect a specific volume
func (c *Client) VolumeInspect(instanceID string) (reply *types.Volume, err error) {
	url := fmt.Sprintf("/admin/volumes/%s", instanceID)
	if _, err = c.httpGet(url, &reply); err != nil {
		return nil, err
	}
	return reply, nil
}

// VolumeOffer will advertise a volume to schedulers
func (c *Client) VolumeOffer(offer *types.VolumeOfferRequest) (reply *types.Volume, err error) {
	url := "/admin/volumeoffer"
	if _, err = c.httpPost(url, offer, &reply); err != nil {
		return nil, err
	}
	return reply, nil
}

// VolumeOfferRevoke will revoke an offer from schedulers
func (c *Client) VolumeOfferRevoke(offer *types.VolumeOfferRevokeRequest) (reply *types.Volume, err error) {
	url := "/admin/volumeofferrevoke"
	if _, err = c.httpPost(url, offer, &reply); err != nil {
		return nil, err
	}
	return reply, nil
}

// VolumeLabel creates labels on a volume
func (c *Client) VolumeLabel(lr *types.VolumeLabelRequest) (reply *types.Volume, err error) {
	url := "/admin/volumelabel"
	if _, err = c.httpPost(url, lr, &reply); err != nil {
		return nil, err
	}
	return reply, nil
}

// VolumeLabelsRemove removes labels from a volume
func (c *Client) VolumeLabelsRemove(lr *types.VolumeLabelsRemoveRequest) (reply *types.Volume, err error) {
	url := "/admin/volumelabelsremove"
	if _, err = c.httpPost(url, lr, &reply); err != nil {
		return nil, err
	}
	return reply, nil
}

// VolumeCreate create a volume
func (c *Client) VolumeCreate(lr *types.VolumeCreateRequest) (reply *types.Volume, err error) {
	url := "/admin/volumes"
	if _, err = c.httpPost(url, lr, &reply); err != nil {
		return nil, err
	}
	return reply, nil
}

// VolumeRemove removes a volume
func (c *Client) VolumeRemove(volumeID string) (err error) {
	url := fmt.Sprintf("/admin/volumes/%s", volumeID)
	if _, err = c.httpDelete(url, nil); err != nil {
		return err
	}
	return nil
}
