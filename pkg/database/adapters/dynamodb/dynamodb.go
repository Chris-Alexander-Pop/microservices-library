package dynamodb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/chris-alexander-pop/system-design-library/pkg/database"
	"github.com/chris-alexander-pop/system-design-library/pkg/errors"
)

// New creates a new DynamoDB client
func New(cfg database.Config) (*dynamodb.Client, error) {
	if cfg.Driver != database.DriverDynamoDB {
		return nil, errors.New(errors.CodeInvalidArgument, fmt.Sprintf("invalid driver %s for dynamodb adapter", cfg.Driver), nil)
	}

	// Load AWS Config (uses env vars by default, but we can override region)
	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to load aws config")
	}

	if cfg.Region != "" {
		awsCfg.Region = cfg.Region
	}

	opts := []func(*dynamodb.Options){}

	// Cloud Support: Custom Endpoint for LocalStack or other compatible APIs
	if cfg.Host != "" {
		endpoint := cfg.Host
		if cfg.Port != "" {
			endpoint = fmt.Sprintf("http://%s:%s", cfg.Host, cfg.Port)
		}

		opts = append(opts, func(o *dynamodb.Options) {
			o.BaseEndpoint = &endpoint
		})
	}

	client := dynamodb.NewFromConfig(awsCfg, opts...)
	return client, nil
}
