package state

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

type ConnectionRequest struct {
	ID          int64     `json:"id"`
	ProfileURL  string    `json:"profile_url"`
	ProfileName string    `json:"profile_name"`
	Note        string    `json:"note"`
	Status      string    `json:"status"` // sent, accepted, rejected
	SentAt      time.Time `json:"sent_at"`
	AcceptedAt  *time.Time `json:"accepted_at,omitempty"`
}

type Message struct {
	ID         int64     `json:"id"`
	ProfileURL string    `json:"profile_url"`
	Content    string    `json:"content"`
	SentAt     time.Time `json:"sent_at"`
}

type ActivityLog struct {
	ID        int64     `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Type      string    `json:"type"` // info, success, warning, error
	Details   string    `json:"details"`
}

// NewStore creates a new state store with SQLite
func NewStore(dbPath string) (*Store, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	store := &Store{db: db}
	if err := store.migrate(); err != nil {
		return nil, err
	}

	return store, nil
}

// migrate creates necessary tables
func (s *Store) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS connection_requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		profile_url TEXT NOT NULL UNIQUE,
		profile_name TEXT,
		note TEXT,
		status TEXT DEFAULT 'sent',
		sent_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		accepted_at DATETIME
	);

	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		profile_url TEXT NOT NULL,
		content TEXT NOT NULL,
		sent_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS activity_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		action TEXT NOT NULL,
		type TEXT NOT NULL,
		details TEXT
	);

	CREATE TABLE IF NOT EXISTS session_data (
		key TEXT PRIMARY KEY,
		value TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_connection_status ON connection_requests(status);
	CREATE INDEX IF NOT EXISTS idx_activity_timestamp ON activity_logs(timestamp DESC);
	`

	_, err := s.db.Exec(schema)
	return err
}

// SaveConnectionRequest stores a new connection request
func (s *Store) SaveConnectionRequest(profileURL, profileName, note string) error {
	_, err := s.db.Exec(
		"INSERT INTO connection_requests (profile_url, profile_name, note) VALUES (?, ?, ?)",
		profileURL, profileName, note,
	)
	return err
}

// GetConnectionsSentToday returns count of connections sent today
func (s *Store) GetConnectionsSentToday() (int, error) {
	var count int
	today := time.Now().Format("2006-01-02")
	err := s.db.QueryRow(
		"SELECT COUNT(*) FROM connection_requests WHERE DATE(sent_at) = ?",
		today,
	).Scan(&count)
	return count, err
}

// GetMessagesSentToday returns count of messages sent today
func (s *Store) GetMessagesSentToday() (int, error) {
	var count int
	today := time.Now().Format("2006-01-02")
	err := s.db.QueryRow(
		"SELECT COUNT(*) FROM messages WHERE DATE(sent_at) = ?",
		today,
	).Scan(&count)
	return count, err
}

// SaveMessage stores a sent message
func (s *Store) SaveMessage(profileURL, content string) error {
	_, err := s.db.Exec(
		"INSERT INTO messages (profile_url, content) VALUES (?, ?)",
		profileURL, content,
	)
	return err
}

// LogActivity adds an entry to the activity log
func (s *Store) LogActivity(action, logType, details string) error {
	_, err := s.db.Exec(
		"INSERT INTO activity_logs (action, type, details) VALUES (?, ?, ?)",
		action, logType, details,
	)
	return err
}

// GetRecentActivityLogs retrieves the most recent activity logs
func (s *Store) GetRecentActivityLogs(limit int) ([]ActivityLog, error) {
	rows, err := s.db.Query(
		"SELECT id, timestamp, action, type, details FROM activity_logs ORDER BY timestamp DESC LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []ActivityLog
	for rows.Next() {
		var log ActivityLog
		if err := rows.Scan(&log.ID, &log.Timestamp, &log.Action, &log.Type, &log.Details); err != nil {
			continue
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// ProfileExists checks if a profile URL has already been contacted
func (s *Store) ProfileExists(profileURL string) (bool, error) {
	var count int
	err := s.db.QueryRow(
		"SELECT COUNT(*) FROM connection_requests WHERE profile_url = ?",
		profileURL,
	).Scan(&count)
	return count > 0, err
}

// SaveSessionData stores arbitrary session data
func (s *Store) SaveSessionData(key, value string) error {
	_, err := s.db.Exec(
		"INSERT OR REPLACE INTO session_data (key, value, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP)",
		key, value,
	)
	return err
}

// GetSessionData retrieves session data by key
func (s *Store) GetSessionData(key string) (string, error) {
	var value string
	err := s.db.QueryRow("SELECT value FROM session_data WHERE key = ?", key).Scan(&value)
	return value, err
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}
