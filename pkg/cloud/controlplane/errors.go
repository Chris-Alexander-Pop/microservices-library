package controlplane

import "github.com/chris-alexander-pop/system-design-library/pkg/errors"

var (
	// ErrHostNotFound is returned when a requested host ID is not registered.
	ErrHostNotFound = errors.NotFound("host not found", nil)

	// ErrHostAlreadyRegistered is returned when attempting to register a host ID that already exists.
	ErrHostAlreadyRegistered = errors.Conflict("host already registered", nil)
)
