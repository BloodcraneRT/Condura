package main

import (
	"database/sql"
	"fmt"

	"github.com/edgesynth/edgesynth/pkg/models"
	"github.com/google/uuid"
	_ "github.com/ClickHouse/clickhouse-go/v2"
)

type DB struct {
	*sql.DB
}

func initDB(connStr string) (*DB, error) {
	db, err := sql.Open("clickhouse", connStr)
	if err != nil {
		return nil, err
	}

	// Clickhouse go driver doesn't support executing multiple statements at once in a single db.Exec.
	// Split by semicolon.
	statements := []string{
		`CREATE TABLE IF NOT EXISTS tasks (id String, type String, target String, interval Int32, config String, enabled Bool) ENGINE = ReplacingMergeTree() ORDER BY id;`,
		`CREATE TABLE IF NOT EXISTS results (id String, task_id String, timestamp DateTime, success Bool, latency_ms Float64, error_msg String, bytes_sent Int64 DEFAULT 0, bytes_recv Int64 DEFAULT 0) ENGINE = MergeTree() ORDER BY (timestamp, task_id);`,
		`CREATE TABLE IF NOT EXISTS sources (id String, name String, target String, is_default Bool) ENGINE = ReplacingMergeTree() ORDER BY id;`,
	}

	for _, stmt := range statements {
		_, err = db.Exec(stmt)
		if err != nil {
			return nil, fmt.Errorf("failed creating schema with stmt %s: %w", stmt, err)
		}
	}

	// Seed default sources if table is empty
	var count uint64
	db.QueryRow(`SELECT count() FROM sources`).Scan(&count)
	if count == 0 {
		defaultSources := []models.Source{
			{ID: uuid.New().String(), Name: "Google DNS", Target: "8.8.8.8", IsDefault: true},
			{ID: uuid.New().String(), Name: "Cloudflare DNS", Target: "1.1.1.1", IsDefault: true},
			{ID: uuid.New().String(), Name: "Google Web", Target: "https://google.com", IsDefault: true},
			{ID: uuid.New().String(), Name: "Local 10MB Download", Target: "http://localhost:8080/api/v1/testdata/download?sizeBytes=10485760", IsDefault: true},
			{ID: uuid.New().String(), Name: "Local Data Upload", Target: "http://localhost:8080/api/v1/testdata/upload", IsDefault: true},
			{ID: uuid.New().String(), Name: "JSONPlaceholder API", Target: "https://jsonplaceholder.typicode.com/posts/1", IsDefault: true},
			{ID: uuid.New().String(), Name: "Public PokeAPI", Target: "https://pokeapi.co/api/v2/pokemon/ditto", IsDefault: true},
		}
		for _, s := range defaultSources {
			db.Exec(`INSERT INTO sources (id, name, target, is_default) VALUES (?, ?, ?, ?)`,
				s.ID, s.Name, s.Target, s.IsDefault)
		}
	}

	return &DB{db}, nil
}

func (db *DB) GetEnabledTasks() ([]models.Task, error) {
	rows, err := db.Query(`SELECT id, type, target, interval, config, enabled FROM tasks WHERE enabled = 1`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		err := rows.Scan(&t.ID, &t.Type, &t.Target, &t.Interval, &t.Config, &t.Enabled)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (db *DB) InsertResult(r models.TaskResult) error {
	id := uuid.New().String()
	_, err := db.Exec(`INSERT INTO results (id, task_id, timestamp, success, latency_ms, error_msg, bytes_sent, bytes_recv) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id, r.TaskID, r.Timestamp, r.Success, r.LatencyMs, r.ErrorMsg, r.BytesSent, r.BytesRecv)
	return err
}

func (db *DB) GetRecentResults(limit int) ([]models.TaskResult, error) {
	rows, err := db.Query(`SELECT id, task_id, timestamp, success, latency_ms, error_msg, bytes_sent, bytes_recv FROM results ORDER BY timestamp DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.TaskResult
	for rows.Next() {
		var r models.TaskResult
		err := rows.Scan(&r.ID, &r.TaskID, &r.Timestamp, &r.Success, &r.LatencyMs, &r.ErrorMsg, &r.BytesSent, &r.BytesRecv)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

func (db *DB) CreateTask(t models.Task) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	_, err := db.Exec(`INSERT INTO tasks (id, type, target, interval, config, enabled) VALUES (?, ?, ?, ?, ?, ?)`,
		t.ID, t.Type, t.Target, t.Interval, t.Config, t.Enabled)
	return err
}

func (db *DB) GetSources() ([]models.Source, error) {
	rows, err := db.Query(`SELECT id, name, target, is_default FROM sources`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []models.Source
	for rows.Next() {
		var s models.Source
		err := rows.Scan(&s.ID, &s.Name, &s.Target, &s.IsDefault)
		if err != nil {
			return nil, err
		}
		sources = append(sources, s)
	}
	return sources, nil
}

func (db *DB) CreateSource(s models.Source) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	_, err := db.Exec(`INSERT INTO sources (id, name, target, is_default) VALUES (?, ?, ?, ?)`,
		s.ID, s.Name, s.Target, s.IsDefault)
	return err
}
