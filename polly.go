package polly

import (
	"github.com/akutz/gofig"

	core "github.com/emccode/polly/core"
	ctypes "github.com/emccode/polly/core/types"
)

// NewWithConfigFile creates a new Polly instance and configures it with a
// custom configuration file.
func NewWithConfigFile(path string) (*ctypes.Polly, error) {
	return core.NewWithConfigFile(path)
}

// NewWithConfig creates a new Polly instance and configures it with a
// custom configuration stream.
func NewWithConfig(config gofig.Config) *ctypes.Polly {
	return core.NewWithConfig(config)
}
