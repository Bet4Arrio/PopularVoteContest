package sqlrepo

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Migrate runs every *.sql file found in the embedded migrations directory in
// lexicographic order. Each file is executed as a single statement batch so you
// can place multiple statements separated by semicolons inside one file.
func Migrate(ctx context.Context, db *DB) error {
	// Ensure the migrations tracking table exists.
	const createTable = `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version     VARCHAR PRIMARY KEY,
			applied_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`
	if _, err := db.Pool.Exec(ctx, createTable); err != nil {
		return fmt.Errorf("migrate: create schema_migrations: %w", err)
	}

	// Collect already-applied versions.
	rows, err := db.Pool.Query(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return fmt.Errorf("migrate: query applied migrations: %w", err)
	}
	applied := make(map[string]bool)
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			rows.Close()
			return fmt.Errorf("migrate: scan version: %w", err)
		}
		applied[v] = true
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("migrate: rows error: %w", err)
	}

	// Read migration files.
	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("migrate: read migrations dir: %w", err)
	}

	// Sort files lexicographically (001_..., 002_..., etc.).
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		version := entry.Name()
		if applied[version] {
			continue
		}

		data, err := migrationsFS.ReadFile("migrations/" + version)
		if err != nil {
			return fmt.Errorf("migrate: read file %s: %w", version, err)
		}

		// Execute the whole file as one batch.
		if _, err := db.Pool.Exec(ctx, string(data)); err != nil {
			return fmt.Errorf("migrate: execute %s: %w", version, err)
		}

		// Record the applied migration.
		if _, err := db.Pool.Exec(ctx,
			`INSERT INTO schema_migrations (version) VALUES ($1)`, version,
		); err != nil {
			return fmt.Errorf("migrate: record %s: %w", version, err)
		}
	}

	return nil
}
