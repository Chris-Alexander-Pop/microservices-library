/*
Package sql provides a unified interface for SQL database access.

This package supports multiple SQL backends through a common interface:
  - PostgreSQL: Production-grade relational database
  - MySQL: Popular open-source relational database
  - SQLite: Embedded database for development and testing
  - SQL Server (MSSQL): Microsoft's enterprise database
  - ClickHouse: Column-oriented OLAP database

Basic usage:

	import (
		"github.com/chris-alexander-pop/system-design-library/pkg/database/sql"
		"github.com/chris-alexander-pop/system-design-library/pkg/database/sql/adapters/postgres"
	)

	cfg := sql.Config{
		Driver:   "postgres",
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "secret",
		Name:     "mydb",
	}

	db, err := postgres.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Use GORM for queries
	gormDB := db.Get(ctx)
*/
package sql
