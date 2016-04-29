package client

import (
	"github.com/emccode/polly/api/types"
)

// Client is the libStorage client.
type Client interface {
	// Volumes returns the registered volumes
	Volumes() ([]*types.Volume, error)

	// VolumesAll returns all volumes
	VolumesAll() ([]*types.Volume, error)

	// VolumeInspect will retrieve details about a volume
	VolumeInspect(volumeID string) (*types.Volume, error)

	// VolumeOffer will advertise a volume to	schedulers
	VolumeOffer(volumeID string, schedulers []string) (*types.Volume, error)

	// VolumeOfferRevoke will revoke a volume offer from schedulers
	VolumeOfferRevoke(volumeID string, schedulers []string) (*types.Volume, error)

	// VolumeLabel creates labels on a volume
	VolumeLabel(volumeID string, labels []string) (*types.Volume, error)

	// VolumeLabelsRemove removes labels from a volume
	VolumeLabelsRemove(volumeID string, labels []string) (*types.Volume, error)

	// VolumeCreate creates a volume
	VolumeCreate(service, name, volumeType string, size, IOPS int64, availabilityZone string, schedulers, labels, fields []string) (*types.Volume, error)

	// VolumeRemove removes a volume
	VolumeRemove(volumeID string) error
}
