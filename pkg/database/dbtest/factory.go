package dbtest

import (
	"context"
	"fmt"
	"sync"

	"github.com/alicebob/miniredis/v2"
	"github.com/chris-alexander-pop/system-design-library/pkg/database"
	"github.com/chris-alexander-pop/system-design-library/pkg/database/vector"
	"github.com/chris-alexander-pop/system-design-library/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MemoryFactory tracks resources that need cleanup (like miniredis)
type MemoryFactory struct {
	cleanups []func()
	mu       sync.Mutex
}

// Close cleans up all tracked resources
func (f *MemoryFactory) Close() {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, c := range f.cleanups {
		c()
	}
	f.cleanups = nil
}

// Factory conforms to database.ConnectionFactory
func (f *MemoryFactory) Factory(cfg database.Config) (interface{}, error) {
	switch cfg.Type {
	case database.StoreTypeSQL:
		// SQLite Memory
		dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", cfg.Name)
		return gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	case database.StoreTypeKeyValue:
		// Miniredis
		mr, err := miniredis.Run()
		if err != nil {
			return nil, errors.Wrap(err, "failed to start miniredis")
		}

		// Track cleanup
		f.mu.Lock()
		f.cleanups = append(f.cleanups, mr.Close)
		f.mu.Unlock()

		// Connect generic go-redis client to it
		client := redis.NewClient(&redis.Options{
			Addr: mr.Addr(),
		})
		return client, nil

	case database.StoreTypeDocument:
		return &MemoryDocumentStore{Data: make(map[string][]interface{})}, nil

	case database.StoreTypeVector:
		return &MemoryVectorStore{}, nil

	default:
		// Fallback for unset Type -> assume SQL if Driver is postgres/mysql generic
		if cfg.Driver == database.DriverPostgres || cfg.Driver == database.DriverMySQL || cfg.Driver == database.DriverSQLite {
			dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", cfg.Name)
			return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
		}
		return nil, errors.New(errors.CodeInvalidArgument, fmt.Sprintf("unsupported type %s for dbtest", cfg.Type), nil)
	}
}

// NewFactory creates a new MemoryFactory.
// Use this if you want to be able to cleanup resources via f.Close()
func NewFactory() *MemoryFactory {
	return &MemoryFactory{}
}

// NewConnectionFactory returns a factory function only.
// Note: This does NOT expose the cleanup capability easily.
// Use NewFactory() if cleanup is required.
func NewConnectionFactory() database.ConnectionFactory {
	f := NewFactory()
	return f.Factory
}

// MemoryVectorStore implements vector.Store in memory
type MemoryVectorStore struct {
	Vectors map[string][]float32
	mu      sync.RWMutex
}

func (m *MemoryVectorStore) Search(ctx context.Context, queryVector []float32, limit int) ([]vector.Result, error) {
	// Dummy implementation returning empty results
	return []vector.Result{}, nil
}

func (m *MemoryVectorStore) Upsert(ctx context.Context, id string, vec []float32, meta map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Vectors == nil {
		m.Vectors = make(map[string][]float32)
	}
	m.Vectors[id] = vec
	return nil
}

func (m *MemoryVectorStore) Delete(ctx context.Context, ids ...string) error {
	return nil
}

// MemoryDocumentStore mock
type MemoryDocumentStore struct {
	Data map[string][]interface{} // Collection -> Docs
	mu   sync.Mutex
}

func (m *MemoryDocumentStore) Insert(ctx context.Context, collection string, doc interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Data == nil {
		m.Data = make(map[string][]interface{})
	}
	m.Data[collection] = append(m.Data[collection], doc)
	return nil
}

func (m *MemoryDocumentStore) Find(ctx context.Context, collection string, query map[string]interface{}) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

func (m *MemoryDocumentStore) Update(ctx context.Context, collection string, filter map[string]interface{}, update map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	// No-op for now unless we implement full expression matching
	return nil
}

func (m *MemoryDocumentStore) Delete(ctx context.Context, collection string, filter map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	// No-op for now
	return nil
}
