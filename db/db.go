package db

import (
	"context"
	"database/sql"
	"time"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (db *Store) CreateJobsTable() error {
	_, err := db.db.Exec(`CREATE TABLE IF NOT EXISTS jobs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		slug TEXT NOT NULL,
		commit_sha TEXT NOT NULL UNIQUE
	)`)
	return err
}

func (db *Store) CreateFilesTable() error {
	_, err := db.db.Exec(`CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		job_id INTEGER NOT NULL,
    	name TEXT NOT NULL,
		content TEXT NOT NULL,

		FOREIGN KEY(job_id) REFERENCES jobs(id)
	)`)

	return err
}

func (s *Store) CreateJob(ctx context.Context, slug, commitSHA string) (string, error) {
	query := `
	INSERT INTO jobs (slug, commit_sha)
	VALUES (?, ?)
	RETURNING id`

	var jobID string
	if err := s.db.QueryRow(query,
		slug, commitSHA,
	).Scan(&jobID); err != nil {
		return "", err
	}

	return jobID, nil
}

func (s *Store) CreateFile(ctx context.Context, jobID, name string, content []byte) error {
	query := `
	INSERT INTO files (job_id, name, content)
	VALUES (?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		jobID, name, string(content),
	)

	return err
}

func (s *Store) GetLatestJobBySlug(ctx context.Context, slug string) (string, error) {
	query := `
	SELECT id
	FROM jobs
	WHERE slug = ?
	ORDER BY created_at DESC
	LIMIT 1
	`
	var jobID string
	if err := s.db.QueryRowContext(ctx, query, slug).Scan(&jobID); err != nil {
		return "", err
	}

	return jobID, nil
}

type CooklangFile struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
}

func (s *Store) GetFilesByJob(ctx context.Context, jobID string) ([]CooklangFile, error) {
	query := `
	SELECT created_at, name, content
	FROM files
	WHERE job_id = ?
	`
	rows, err := s.db.QueryContext(ctx, query, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []CooklangFile
	for rows.Next() {
		var entry CooklangFile
		if err := rows.Scan(&entry.CreatedAt, &entry.Name, &entry.Content); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
