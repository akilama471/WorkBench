package database

import (
	"database/sql"
	"fmt"
)

type Migration struct {
	Version     int
	Description string
	Up          func(tx *sql.Tx) error
}

func (d *Database) runMigrations() error {
	_, err := d.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_version (
			version INTEGER PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create schema_version table: %w", err)
	}

	migrations := getMigrations()

	for _, m := range migrations {
		applied, err := d.isMigrationApplied(m.Version)
		if err != nil {
			return fmt.Errorf("failed to check migration %d: %w", m.Version, err)
		}

		if applied {
			continue
		}

		tx, err := d.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %d: %w", m.Version, err)
		}

		if err := m.Up(tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration %d (%s) failed: %w", m.Version, m.Description, err)
		}

		_, err = tx.Exec(`INSERT INTO schema_version (version) VALUES (?)`, m.Version)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %d: %w", m.Version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d: %w", m.Version, err)
		}
	}

	return nil
}

func (d *Database) isMigrationApplied(version int) (bool, error) {
	var count int
	err := d.db.QueryRow(`SELECT COUNT(*) FROM schema_version WHERE version = ?`, version).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func getMigrations() []Migration {
	return []Migration{
		{
			Version:     1,
			Description: "create settings table",
			Up: func(tx *sql.Tx) error {
				_, err := tx.Exec(`
					CREATE TABLE IF NOT EXISTS settings (
						key TEXT PRIMARY KEY,
						value TEXT NOT NULL,
						created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
						updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
					)
				`)
				return err
			},
		},
		{
			Version:     2,
			Description: "create installed_packages table",
			Up: func(tx *sql.Tx) error {
				_, err := tx.Exec(`
					CREATE TABLE IF NOT EXISTS installed_packages (
						id TEXT NOT NULL,
						name TEXT NOT NULL,
						version TEXT NOT NULL,
						platform TEXT NOT NULL,
						architecture TEXT NOT NULL,
						type TEXT NOT NULL,
						install_path TEXT NOT NULL,
						installed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
						PRIMARY KEY (id, version, platform, architecture)
					)
				`)
				return err
			},
		},
		{
			Version:     3,
			Description: "create projects table",
			Up: func(tx *sql.Tx) error {
				_, err := tx.Exec(`
					CREATE TABLE IF NOT EXISTS projects (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						name TEXT NOT NULL,
						path TEXT NOT NULL UNIQUE,
						type TEXT NOT NULL,
						created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
						updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
					)
				`)
				return err
			},
		},
		{
			Version:     4,
			Description: "create active_runtime table",
			Up: func(tx *sql.Tx) error {
				_, err := tx.Exec(`
					CREATE TABLE IF NOT EXISTS active_runtimes (
						runtime_id TEXT PRIMARY KEY,
						version TEXT NOT NULL,
						activated_at DATETIME DEFAULT CURRENT_TIMESTAMP
					)
				`)
				return err
			},
		},
	}
}
