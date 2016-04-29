package client

import (
	"github.com/emccode/polly/api/types"
	"strings"
)

func (c *pc) Volumes() ([]*types.Volume, error) {
	return c.Client.Volumes()
}

func (c *pc) VolumesAll() ([]*types.Volume, error) {
	return c.Client.VolumesAll()
}

func (c *pc) VolumeInspect(volumeID string) (*types.Volume, error) {
	return c.Client.VolumeInspect(volumeID)
}

func (c *pc) VolumeOffer(volumeID string,
	schedulers []string) (*types.Volume, error) {
	offer := &types.VolumeOfferRequest{
		VolumeID:   volumeID,
		Schedulers: schedulers,
	}
	return c.Client.VolumeOffer(offer)
}

func (c *pc) VolumeOfferRevoke(volumeID string,
	schedulers []string) (*types.Volume, error) {
	offer := &types.VolumeOfferRevokeRequest{
		VolumeID:   volumeID,
		Schedulers: schedulers,
	}
	return c.Client.VolumeOfferRevoke(offer)
}

func splitLabel(label string) (key, value string) {
	kv := strings.SplitN(label, "=", 2)
	if len(kv) != 2 {
		return "", ""
	}
	return kv[0], kv[1]
}

func labelMap(labels []string) map[string]string {
	lm := make(map[string]string)
	for _, l := range labels {
		if k, v := splitLabel(l); k != "" {
			lm[k] = v
		}
	}
	return lm
}

func (c *pc) VolumeLabel(volumeID string,
	labels []string) (*types.Volume, error) {
	lc := &types.VolumeLabelRequest{
		VolumeID: volumeID,
		Labels:   labelMap(labels),
	}
	return c.Client.VolumeLabel(lc)
}

func (c *pc) VolumeLabelsRemove(volumeID string,
	labels []string) (*types.Volume, error) {
	lc := &types.VolumeLabelsRemoveRequest{
		VolumeID: volumeID,
		Labels:   labels,
	}
	return c.Client.VolumeLabelsRemove(lc)
}

// VolumeCreate creates a volume
func (c *pc) VolumeCreate(service, name, volumeType string,
	size, IOPS int64, availabilityZone string,
	schedulers, labels, fields []string) (*types.Volume, error) {
	lc := &types.VolumeCreateRequest{
		ServiceName:      service,
		Name:             name,
		VolumeType:       volumeType,
		Size:             size,
		IOPS:             IOPS,
		AvailabilityZone: availabilityZone,
		Schedulers:       schedulers,
		Labels:           labelMap(labels),
		Fields:           labelMap(fields),
	}
	return c.Client.VolumeCreate(lc)
}

// VolumeRemove removes a volume
func (c *pc) VolumeRemove(volumeID string) error {
	return c.Client.VolumeRemove(volumeID)
}
