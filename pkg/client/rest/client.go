package rest

import (
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Config struct {
	Timeout   time.Duration `env:"CLIENT_TIMEOUT" env-default:"30s"`
	Retries   int           `env:"CLIENT_RETRIES" env-default:"3"`
	UserAgent string        `env:"CLIENT_USER_AGENT" env-default:"system-design-library-client"`
}

// New creates a robust HTTP client with Retries and OTel Tracing
func New(cfg Config) *http.Client {
	// 1. Retryable Client
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = cfg.Retries
	retryClient.HTTPClient.Timeout = cfg.Timeout
	retryClient.Logger = nil // Use our own logger or silences generic logs? For now silence.

	// 2. Wrap Transport with OTel
	// retryablehttp uses proper RoundTripper
	baseTransport := retryClient.HTTPClient.Transport
	if baseTransport == nil {
		baseTransport = http.DefaultTransport
	}

	otelTransport := otelhttp.NewTransport(baseTransport)

	// update the client to use the instrumented transport
	// Note: retryablehttp wraps the standard client. We need to inject OTel INSIDE or OUTSIDE?
	// If we wrap `retryClient.StandardClient()`, the retry logic (which is inside retryClient) might be outside OTel?
	// retryablehttp logic is: Client.Do -> Retries -> RoundTripper -> Network.
	// If we wrap the RoundTripper, OTel sees each attempt? OR one span for the whole op?
	// Usually one span per attempt is better for debugging, or one parent span.
	// otelhttp.NewTransport wraps the *outgoing* request.

	retryClient.HTTPClient.Transport = otelTransport

	// 3. Return standard client interface
	stdClient := retryClient.StandardClient()

	return stdClient
}
