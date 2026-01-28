package metering

import "github.com/chris-alexander-pop/system-design-library/pkg/errors"

var (
	// ErrRateNotFound is returned when pricing information is missing for a resource type.
	ErrRateNotFound = errors.NotFound("rate not found", nil)

	// ErrInvalidUsage is returned when usage data is malformed.
	ErrInvalidUsage = errors.InvalidArgument("invalid usage data", nil)
)
