package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Database struct {
	db *sql.DB
}

func Open(path string) (*Database, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database %s: %w", path, err)
	}

	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) DB() *sql.DB {
	return d.db
}

func (d *Database) Initialize() error {
	if err := d.runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

func (d *Database) SetSetting(key, value string) error {
	_, err := d.db.Exec(
		`INSERT INTO settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value`,
		key, value,
	)
	return err
}

func (d *Database) GetSetting(key string) (string, error) {
	var value string
	err := d.db.QueryRow(`SELECT value FROM settings WHERE key = ?`, key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return value, nil
}

func (d *Database) DeleteSetting(key string) error {
	_, err := d.db.Exec(`DELETE FROM settings WHERE key = ?`, key)
	return err
}

func (d *Database) GetAllSettings() (map[string]string, error) {
	rows, err := d.db.Query(`SELECT key, value FROM settings`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		result[key] = value
	}
	return result, nil
}

type ProjectRecord struct {
	ID        int
	Name      string
	Path      string
	Type      string
	CreatedAt string
	UpdatedAt string
}

func (d *Database) SaveProject(name, path, projectType string) error {
	_, err := d.db.Exec(
		`INSERT INTO projects (name, path, type, updated_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		 ON CONFLICT(path) DO UPDATE SET name = excluded.name, type = excluded.type, updated_at = CURRENT_TIMESTAMP`,
		name, path, projectType,
	)
	return err
}

func (d *Database) GetAllProjects() ([]*ProjectRecord, error) {
	rows, err := d.db.Query(`SELECT id, name, path, type, created_at, updated_at FROM projects ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*ProjectRecord
	for rows.Next() {
		var p ProjectRecord
		if err := rows.Scan(&p.ID, &p.Name, &p.Path, &p.Type, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, &p)
	}
	return projects, nil
}

func (d *Database) GetProjectByPath(path string) (*ProjectRecord, error) {
	var p ProjectRecord
	err := d.db.QueryRow(`SELECT id, name, path, type, created_at, updated_at FROM projects WHERE path = ?`, path).
		Scan(&p.ID, &p.Name, &p.Path, &p.Type, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (d *Database) DeleteProjectByPath(path string) error {
	_, err := d.db.Exec(`DELETE FROM projects WHERE path = ?`, path)
	return err
}

