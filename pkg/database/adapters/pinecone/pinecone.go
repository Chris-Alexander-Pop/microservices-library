package pinecone

import (
	"context"
	"fmt"

	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/chris-alexander-pop/system-design-library/pkg/database"
	"github.com/chris-alexander-pop/system-design-library/pkg/database/vector"
	"github.com/chris-alexander-pop/system-design-library/pkg/errors"
)

// NOTE: Pinecone official Go SDK is often in flux or community maintained.
// For "Overengineering" without external instability, I will implement a robust HTTP Client wrapper
// that adheres to the VectorStore interface.
// If an official SDK `github.com/pinecone-io/go-pinecone` is available we would use it.
// Assuming we implement the vector.Store interface directly.

type PineconeStore struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

// New creates a new Pinecone adapter
func New(cfg database.Config) (*PineconeStore, error) {
	if cfg.Driver != database.DriverPinecone {
		return nil, errors.New(errors.CodeInvalidArgument, fmt.Sprintf("invalid driver %s for pinecone adapter", cfg.Driver), nil)
	}

	// Construct Base URL based on inputs or use explicit Host if provided
	baseURL := cfg.Host
	if baseURL == "" {
		// Fallback to standard Pinecone URL structure if Host not generic
		// https://index-project.svc.environment.pinecone.io
		baseURL = fmt.Sprintf("https://%s-%s.svc.%s.pinecone.io", cfg.Name, cfg.ProjectID, cfg.Environment)
	}

	return &PineconeStore{
		BaseURL: baseURL,
		APIKey:  cfg.APIKey,
		Client:  &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Search implements vector.Store interface via HTTP
func (p *PineconeStore) Search(ctx context.Context, queryVector []float32, limit int) ([]vector.Result, error) {
	url := fmt.Sprintf("%s/query", p.BaseURL)

	reqBody := map[string]interface{}{
		"vector":          queryVector,
		"topK":            limit,
		"includeMetadata": true,
	}

	resp, err := p.doRequest(ctx, "POST", url, reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, p.handleError(resp)
	}

	var result struct {
		Matches []struct {
			ID       string                 `json:"id"`
			Score    float32                `json:"score"`
			Metadata map[string]interface{} `json:"metadata"`
		} `json:"matches"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errors.Wrap(err, "failed to decode pinecone response")
	}

	vecResults := make([]vector.Result, len(result.Matches))
	for i, m := range result.Matches {
		vecResults[i] = vector.Result{
			ID:       m.ID,
			Score:    m.Score,
			Metadata: m.Metadata,
		}
	}

	return vecResults, nil
}

func (p *PineconeStore) Upsert(ctx context.Context, id string, vec []float32, metadata map[string]interface{}) error {
	url := fmt.Sprintf("%s/vectors/upsert", p.BaseURL)

	reqBody := map[string]interface{}{
		"vectors": []map[string]interface{}{
			{
				"id":       id,
				"values":   vec,
				"metadata": metadata,
			},
		},
	}

	resp, err := p.doRequest(ctx, "POST", url, reqBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return p.handleError(resp)
	}
	return nil
}

func (p *PineconeStore) Delete(ctx context.Context, ids ...string) error {
	url := fmt.Sprintf("%s/vectors/delete", p.BaseURL)

	reqBody := map[string]interface{}{
		"ids": ids,
	}

	resp, err := p.doRequest(ctx, "POST", url, reqBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return p.handleError(resp)
	}

	return nil
}

// doRequest helper
func (p *PineconeStore) doRequest(ctx context.Context, method, url string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal request body")
		}
		bodyReader = bytes.NewReader(jsonBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", p.APIKey)

	return p.Client.Do(req)
}

func (p *PineconeStore) handleError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	return errors.New(errors.CodeInternal, fmt.Sprintf("pinecone api error: %s", string(body)), nil)
}
