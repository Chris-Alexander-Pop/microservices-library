package timeseries

import (
	"context"
	"time"
)

// Point represents a single data point in a time-series database.
type Point struct {
	// Measurement is the name of the measurement or table.
	Measurement string

	// Tags are key-value pairs used for indexing and filtering.
	Tags map[string]string

	// Fields are the actual values to be stored.
	Fields map[string]interface{}

	// Time is the timestamp of the data point.
	Time time.Time
}

// Config holds configuration for the timeseries database.
type Config struct {
	// Driver specifies the timeseries backend: "timestream", "influxdb", "memory".
	Driver string `env:"TS_DRIVER" env-default:"memory"`

	// Host is the database server hostname or URL.
	Host string `env:"TS_HOST" env-default:"localhost"`

	// Port is the database server port.
	Port int `env:"TS_PORT" env-default:"8086"`

	// Database is the name of the database or bucket.
	Database string `env:"TS_DATABASE" env-default:"telemetry"`

	// User is the username for authentication.
	User string `env:"TS_USER"`

	// Password is the password or token for authentication.
	Password string `env:"TS_PASSWORD"`

	// Region is the AWS region (for Timestream).
	Region string `env:"TS_REGION"`

	// Org is the organization name (for InfluxDB).
	Org string `env:"TS_ORG"`

	// BatchSize is the number of points to write in a single batch.
	BatchSize int `env:"TS_BATCH_SIZE" env-default:"1000"`

	// BatchInterval is the maximum time to wait before writing a batch.
	BatchInterval time.Duration `env:"TS_BATCH_INTERVAL" env-default:"1s"`
}

// Timeseries defines the interface for interacting with time-series databases.
type Timeseries interface {
	// Write writes a single point to the database.
	Write(ctx context.Context, point *Point) error

	// WriteBatch writes a batch of points to the database.
	WriteBatch(ctx context.Context, points []*Point) error

	// Query executes a query against the database and returns the results.
	// The query format is driver-specific (e.g., SQL for Timestream, Flux/InfluxQL for InfluxDB).
	Query(ctx context.Context, query string) ([]*Point, error)

	// Close releases any resources associated with the connection.
	Close() error
}
