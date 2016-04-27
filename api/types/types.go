package types

import (
	lstypes "github.com/emccode/libstorage/api/types"
)

// VolumeOfferRequest contains offer information
type VolumeOfferRequest struct {
	VolumeID   string   `json:"volumeID,omitempty"`
	Schedulers []string `json:"schedulers,omitempty"`
}

// VolumeOfferRevokeRequest contains offer revoke information
type VolumeOfferRevokeRequest struct {
	VolumeID   string   `json:"volumeID,omitempty"`
	Schedulers []string `json:"schedulers,omitempty"`
}

// VolumeLabelRequest creates labels on volumes
type VolumeLabelRequest struct {
	VolumeID string            `json:"volumeID,omitempty"`
	Labels   map[string]string `json:"labels,omitempty"`
}

// VolumeLabelsRemoveRequest removes labels on volumes
type VolumeLabelsRemoveRequest struct {
	VolumeID string   `json:"volumeID,omitempty"`
	Labels   []string `json:"labels,omitempty"`
}

// VolumeCreateRequest creates a volume
type VolumeCreateRequest struct {
	ServiceName      string            `json:"service,omitempty"`
	Name             string            `json:"name,omitempty"`
	VolumeType       string            `json:"volumeType,omitempty"`
	Size             int64             `json:"size,omitempty"`
	IOPS             int64             `json:"iops,omitempty"`
	AvailabilityZone string            `json:"availabilityZone,omitempty"`
	Schedulers       []string          `json:"schedulers,omitempty"`
	Labels           map[string]string `json:"labels,omitempty"`
	Fields           map[string]string `json:"fields,omitempty"`
}

// Volume is a storage libStorage Volume with Polly annotations
type Volume struct {
	*lstypes.Volume // anonymous field from libstorage

	// VolumeID is the Polly VolumeID
	VolumeID string `json:"volumeid,omitempty"`

	// ServiceName comes from libstorage
	ServiceName string `json:"serviceName,omitempty"`

	// Scheduler is the exclusive owner if specified
	Schedulers []string `json:"schedulers,omitempty"`

	// Labels are (admin)user applied via API
	Labels map[string]string `json:"labels,omitempty"`
}

// Snapshot is a libStorage Volume snap with Polly annotations
type Snapshot struct {
	*lstypes.Snapshot // anonymous field

	// SnapshotID is the Polly SnapshotID
	SnapshotID string `json:"snapshotid,omitempty"`

	// the Storage provider identifier
	ServiceName string `json:"serviceName,omitempty"`

	// Scheduler is the exclusive owner if specifier
	Scheduler string `json:"scheduler,omitempty"`
}
