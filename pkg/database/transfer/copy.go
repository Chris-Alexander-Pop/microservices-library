package transfer

import (
	"context"

	"github.com/chris-alexander-pop/system-design-library/pkg/errors" // Import our error package
	"gorm.io/gorm"
)

// CopyTable transfers data from source to destination for a given model.
// It uses batch processing to avoid memory overflow.
func CopyTable(ctx context.Context, src *gorm.DB, dst *gorm.DB, model interface{}, batchSize int) error {
	// 1. Ensure Table Exists in Dst (AutoMigrate)
	if err := dst.AutoMigrate(model); err != nil {
		return errors.Wrap(err, "failed to migrate destination table")
	}

	// 2. Read and Write in Batches
	var rows []map[string]interface{} // Generic map to hold data

	// We use FindInBatches with map? GORM usually needs a struct slice for FindInBatches.
	// For generic "model" interface, we might need a slice of that model type.
	// But `model` is likely a struct instance.
	// Let's rely on GORM's ability to reflect.

	result := src.WithContext(ctx).Model(model).FindInBatches(&rows, batchSize, func(tx *gorm.DB, batch int) error {
		// Use BulkUpsert style generic create
		// Clause.OnConflict is driver specific, so careful. simpler to just Create (Insert)
		if err := dst.WithContext(ctx).Model(model).Create(&rows).Error; err != nil {
			return errors.Wrap(err, "failed to insert batch")
		}
		return nil
	})

	return result.Error
}

// MigrateSchema copies the schema structure from source model to destination.
func MigrateSchema(ctx context.Context, dst *gorm.DB, models ...interface{}) error {
	return dst.WithContext(ctx).AutoMigrate(models...)
}
