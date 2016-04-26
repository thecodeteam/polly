package volumes

import (
	"github.com/akutz/goof"
	lstypes "github.com/emccode/libstorage/api/types"
	"github.com/emccode/polly/api/types"
	ctypes "github.com/emccode/polly/core/types"
)

// Init beings the initialization process for volumes
func Init(p *ctypes.Polly) error {
	vols, err := p.LsClient.Volumes()
	if err != nil {
		return err
	}

	p.Store.EraseStore()

	for _, vol := range vols {
		exists, err := p.Store.Exists(vol)
		if err != nil {
			return err
		}
		if !exists {
			p.Store.SaveVolumeMetadata(vol)
		}
	}

	// read in volume records from Polly persistent store and cache them
	ids, err := p.Store.GetVolumeIds()
	if err != nil {
		return goof.WithError("failed to retrieve volume IDs form persistent store", err)
	}

	for _, id := range ids {
		v := &types.Volume{
			Volume: &lstypes.Volume{
				ID: id,
			},
		}
		err = p.Store.SetVolumeMetadata(v)
		if err != nil {
			return goof.WithError("failed to retrieve volume metadata from persitent store", err)
		}
	}

	return nil
}
