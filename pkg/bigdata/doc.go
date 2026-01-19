/*
Package bigdata provides detailed interfaces for Big Data processing and analytics.

Features:
  - Generic Client interface for Data Warehouses (BigQuery, Redshift)
  - Compute abstractions (MapReduce, Spark)
  - Data formats (Avro, Parquet)

Usage:

	import "github.com/chris-alexander-pop/system-design-library/pkg/bigdata/adapters/bigquery"

	client := bigquery.New(cfg)
	res, err := client.Query(ctx, "SELECT * FROM users LIMIT 10")
*/
package bigdata
