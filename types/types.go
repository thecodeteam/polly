package types

import (
	lstypes "github.com/emccode/libstorage/api/types"
)

// Volume is a storage libStorage Volume with Polly annotations
type Volume struct {
	*lstypes.Volume // anonymous field from libstorage

	// ServiceName comes from libstorage
	ServiceName string `json:"serviceName,omitempty"`

	// Scheduler is the exclusive owner if specified
	Scheduler string `json:"scheduler,omitempty"`

	// StorageProviderName identifies the Storage provider
	StorageProviderName string `json:"storageProviderName,omitempty"`

	// Labels are (admin)user applied via API
	Labels map[string]string `json:"labels,omitempty"`
}

// Snapshot is a libStorage Volume snap with Polly annotations
type Snapshot struct {
	*lstypes.Snapshot // anonymous field

	// the Storage provider identifier
	ServiceName string `json:"serviceName,omitempty"`

	// Scheduler is the exclusive owner if specifier
	Scheduler string `json:"scheduler,omitempty"`
}

// Provider is a storage libStorage Provider with Polly annotations
type Provider struct {

	// ServiceName comes from libstorage
	ProviderName string `json:"providerName,omitempty"`

	// StorageProviderName identifies the Storage provider
	DriverName string `json:"driverName,omitempty"`

	// Labels are (admin)user applied via API
	Labels map[string]string `json:"labels,omitempty"`
}

// PollyVersion is Polly Version
type PollyVersion struct {
	PollyAdminAPIVersion     string
	PollySchedulerAPIVersion string
	PollyBuildVersion        string
}
