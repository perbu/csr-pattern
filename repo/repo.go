package repo

import (
	"database/sql"
	"fmt"
	"log/slog"
	_ "modernc.org/sqlite"
)

// Repo represents the key-value store repository.
type Repo struct {
	conn *sql.DB
}

type KeyNotFoundError struct {
	key string
}

func (e KeyNotFoundError) Error() string {
	return fmt.Sprintf("key not found: '%s'", e.key)
}

func NewKeyNotFoundError(key string) KeyNotFoundError {
	return KeyNotFoundError{key: key}
}

type KeyExistsError struct {
	key string
}

func (e KeyExistsError) Error() string {
	return fmt.Sprintf("key already exists: '%s'", e.key)
}

func NewKeyExistsError(key string) KeyExistsError {
	return KeyExistsError{key: key}
}

const (
	// DefaultDBPath is the default path to the SQLite database.
	DefaultDBPath = "kv.db"
)

// New initializes a new Repo.
func New(dbPath string, logger *slog.Logger) (*Repo, error) {
	if dbPath == "" {
		dbPath = DefaultDBPath
	}
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create the key-value table if it doesn't exist.
	_, err = conn.Exec("CREATE TABLE IF NOT EXISTS kv (key TEXT PRIMARY KEY, value TEXT)", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}
	return &Repo{conn: conn}, nil
}

// Close closes the repository and the underlying database connection.
func (r *Repo) Close() error {
	return r.conn.Close()
}

// Get retrieves the value associated with the given key.
func (r *Repo) Get(key string) (string, error) {

	rows, err := r.conn.Query("SELECT value FROM kv WHERE key = ?", key)
	if err != nil {
		return "", fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer rows.Close()

	var value string
	if !rows.Next() {
		return "", NewKeyNotFoundError(key)
	}
	err = rows.Scan(&value)
	if err != nil {
		return "", fmt.Errorf("failed to query key-value pair: %w", err)
	}
	return value, nil
}

// Create creates a new key-value pair.
func (r *Repo) Create(key, value string) error {

	_, err := r.conn.Exec("INSERT INTO kv (key, value) VALUES (?, ?)", key, value)
	if err != nil {
		return fmt.Errorf("failed to insert key-value pair: %w", err)
	}
	return nil
}

// Update updates the value associated with the given key.
func (r *Repo) Update(key, value string) error {
	res, err := r.conn.Exec("UPDATE kv SET value = ? WHERE key = ?", value, key)
	if err != nil {
		return fmt.Errorf("failed to update key-value pair: %w", err)
	}
	// return an error if the key doesn't exist
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return NewKeyNotFoundError(key)
	}

	return nil
}

// Delete deletes the key-value pair associated with the given key.
// It returns an error if the key doesn't exist.
func (r *Repo) Delete(key string) error {
	res, err := r.conn.Exec("DELETE FROM kv WHERE key = ?", key)
	if err != nil {
		return fmt.Errorf("failed to delete key-value pair: %w", err)
	}
	// return an error if the key doesn't exist
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return NewKeyNotFoundError(key)
	}
	return nil
}
