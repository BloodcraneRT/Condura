package main

import (
	"database/sql"
	"time"

	"github.com/edgesynth/edgesynth/pkg/api"
	"github.com/edgesynth/edgesynth/pkg/models"
	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func initDB(filepath string) (*DB, error) {
	db, err := sql.Open("sqlite", filepath)
	if err != nil {
		return nil, err
	}

	// Create tables
	schema := `
	CREATE TABLE IF NOT EXISTS agents (
		id TEXT PRIMARY KEY,
		hostname TEXT,
		os TEXT,
		ip TEXT,
		status TEXT,
		last_seen DATETIME
	);
	CREATE TABLE IF NOT EXISTS tasks (
		id TEXT PRIMARY KEY,
		agent_id TEXT,
		type TEXT,
		target TEXT,
		interval INTEGER,
		config TEXT,
		enabled BOOLEAN
	);
	CREATE TABLE IF NOT EXISTS results (
		id TEXT PRIMARY KEY,
		task_id TEXT,
		agent_id TEXT,
		timestamp DATETIME,
		success BOOLEAN,
		latency_ms REAL,
		error_msg TEXT,
		bytes_sent INTEGER DEFAULT 0,
		bytes_recv INTEGER DEFAULT 0
	);
	CREATE TABLE IF NOT EXISTS sources (
		id TEXT PRIMARY KEY,
		name TEXT,
		target TEXT,
		is_default BOOLEAN
	);
	`
	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	// Seed default sources if table is empty
	var count int
	db.QueryRow(`SELECT COUNT(*) FROM sources`).Scan(&count)
	if count == 0 {
		defaultSources := []models.Source{
			{ID: uuid.New().String(), Name: "Google DNS", Target: "8.8.8.8", IsDefault: true},
			{ID: uuid.New().String(), Name: "Cloudflare DNS", Target: "1.1.1.1", IsDefault: true},
			{ID: uuid.New().String(), Name: "Google Web", Target: "https://google.com", IsDefault: true},
			{ID: uuid.New().String(), Name: "Local 10MB Download", Target: "http://localhost:8080/api/v1/testdata/download?sizeBytes=10485760", IsDefault: true},
			{ID: uuid.New().String(), Name: "Local Data Upload", Target: "http://localhost:8080/api/v1/testdata/upload", IsDefault: true},
		}
		for _, s := range defaultSources {
			db.Exec(`INSERT INTO sources (id, name, target, is_default) VALUES (?, ?, ?, ?)`,
				s.ID, s.Name, s.Target, s.IsDefault)
		}
	}

	return &DB{db}, nil
}

func (db *DB) RegisterAgent(req api.RegisterRequest, ip string) (string, error) {
	id := uuid.New().String()
	_, err := db.Exec(`INSERT INTO agents (id, hostname, os, ip, status, last_seen) VALUES (?, ?, ?, ?, 'online', ?)`,
		id, req.Hostname, req.OS, ip, time.Now())
	return id, err
}

func (db *DB) UpdateHeartbeat(agentID string) error {
	_, err := db.Exec(`UPDATE agents SET last_seen = ?, status = 'online' WHERE id = ?`, time.Now(), agentID)
	return err
}

func (db *DB) GetTasksForAgent(agentID string) ([]models.Task, error) {
	rows, err := db.Query(`SELECT id, agent_id, type, target, interval, config, enabled FROM tasks WHERE agent_id = ?`, agentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		err := rows.Scan(&t.ID, &t.AgentID, &t.Type, &t.Target, &t.Interval, &t.Config, &t.Enabled)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (db *DB) InsertResult(r models.TaskResult) error {
	id := uuid.New().String()
	_, err := db.Exec(`INSERT INTO results (id, task_id, agent_id, timestamp, success, latency_ms, error_msg, bytes_sent, bytes_recv) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, r.TaskID, r.AgentID, r.Timestamp, r.Success, r.LatencyMs, r.ErrorMsg, r.BytesSent, r.BytesRecv)
	return err
}

// Additional helpers for UI
func (db *DB) GetAllAgents() ([]models.Agent, error) {
	rows, err := db.Query(`SELECT id, hostname, os, ip, status, last_seen FROM agents`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []models.Agent
	for rows.Next() {
		var a models.Agent
		err := rows.Scan(&a.ID, &a.Hostname, &a.OS, &a.IP, &a.Status, &a.LastSeen)
		if err != nil {
			return nil, err
		}
		agents = append(agents, a)
	}
	return agents, nil
}

func (db *DB) GetRecentResults(limit int) ([]models.TaskResult, error) {
	rows, err := db.Query(`SELECT id, task_id, agent_id, timestamp, success, latency_ms, error_msg, bytes_sent, bytes_recv FROM results ORDER BY timestamp DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.TaskResult
	for rows.Next() {
		var r models.TaskResult
		err := rows.Scan(&r.ID, &r.TaskID, &r.AgentID, &r.Timestamp, &r.Success, &r.LatencyMs, &r.ErrorMsg, &r.BytesSent, &r.BytesRecv)
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
	_, err := db.Exec(`INSERT INTO tasks (id, agent_id, type, target, interval, config, enabled) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		t.ID, t.AgentID, t.Type, t.Target, t.Interval, t.Config, t.Enabled)
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
