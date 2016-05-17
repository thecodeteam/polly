package context

import "github.com/emccode/libstorage/api/context"

// PollyHeaderKey is the type for Polly HTTP header keys.
type PollyHeaderKey int

const (
	// RequestPathHeaderKey is the header key for the Polly-Requestpath header.
	RequestPathHeaderKey PollyHeaderKey = iota
)

func (k PollyHeaderKey) String() string {
	switch k {
	case RequestPathHeaderKey:
		return "Polly-Requestpath"
	}
	return ""
}

func init() {
	context.RegisterCustomKey(RequestPathHeaderKey, context.CustomHeaderKey)
}
