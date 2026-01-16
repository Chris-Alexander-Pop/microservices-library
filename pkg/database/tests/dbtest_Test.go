package tests

import (
	"context"
	"testing"
	"time"

	"github.com/chris-alexander-pop/system-design-library/pkg/database"
	"github.com/chris-alexander-pop/system-design-library/pkg/database/dbtest"
	"github.com/chris-alexander-pop/system-design-library/pkg/test"
	"github.com/redis/go-redis/v9"
)

type DBTestSuite struct {
	*test.Suite
}

func TestDBTestSuite(t *testing.T) {
	test.Run(t, &DBTestSuite{Suite: test.NewSuite()})
}

func (s *DBTestSuite) TestInMemoryFactory() {
	ctx := context.Background()

	// 1. SQL (SQLite Memory)
	sqlCfg := database.ManagerConfig{
		Primary: database.Config{
			Driver: database.DriverSQLite,
			Type:   database.StoreTypeSQL,
			Name:   "mem_sql_test",
		},
	}
	mgrSQL, err := database.NewManager(sqlCfg, dbtest.NewConnectionFactory())
	s.Require().NoError(err)

	db := mgrSQL.Get(ctx)
	s.NotNil(db)
	// Validate it works (create table)
	type Toy struct {
		ID   uint `gorm:"primaryKey"`
		Name string
	}
	s.NoError(db.AutoMigrate(&Toy{}))
	s.NoError(db.Create(&Toy{Name: "Buzz"}).Error)
	var count int64
	db.Model(&Toy{}).Count(&count)
	s.Equal(int64(1), count)

	// 2. KV (Miniredis)
	kvCfg := database.ManagerConfig{
		Primary: database.Config{
			Driver: database.DriverRedis,
			Type:   database.StoreTypeKeyValue,
			Name:   "mem_kv_test", // Not strictly needed for redis
		},
	}
	mgrKV, err := database.NewManager(kvCfg, dbtest.NewConnectionFactory())
	s.Require().NoError(err)

	rdb := mgrKV.GetKV(ctx).(*redis.Client)
	s.NotNil(rdb)
	s.NoError(rdb.Set(ctx, "foo", "bar", time.Minute).Err())
	val, err := rdb.Get(ctx, "foo").Result()
	s.NoError(err)
	s.Equal("bar", val)

	// 3. Document (Mock)
	docCfg := database.ManagerConfig{
		Primary: database.Config{
			Driver: database.DriverMongoDB,
			Type:   database.StoreTypeDocument,
		},
	}
	mgrDoc, err := database.NewManager(docCfg, dbtest.NewConnectionFactory())
	s.Require().NoError(err)

	// Assert to Generic Interface
	docStore, ok := mgrDoc.GetDocument(ctx).(database.DocumentStore)
	s.True(ok, "should implement DocumentStore interface")
	s.NoError(docStore.Insert(ctx, "users", map[string]interface{}{"name": "woody"}))

	// 4. Vector (Mock)
	vecCfg := database.ManagerConfig{
		Primary: database.Config{
			Driver: database.DriverPinecone,
			Type:   database.StoreTypeVector,
		},
	}
	mgrVec, err := database.NewManager(vecCfg, dbtest.NewConnectionFactory())
	s.Require().NoError(err)
	// Assert to Store Interface (if imported, or generic)
	// Since we didn't import vector package here to avoid cycles if any, we rely on GetVector return.
	// Actually we can import vector package in tests safely.
	// But database package imports vector... tests is separate package.
	// wait database -> vector. Tests -> database + vector. Cyclic? No.
	// database_test -> database. database imports vector. database_test imports vector. Fine.
	// BUT dbtest imports vector.
	// dbtest -> vector.
	// It's fine.

	// Check type
	vecStore := mgrVec.GetVector(ctx)
	s.NotNil(vecStore)
}
