package main

import (
	"database/sql"
	"fmt"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/edgesynth/edgesynth/pkg/models"
	"github.com/google/uuid"
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
		`CREATE TABLE IF NOT EXISTS results (id String, task_id String, timestamp DateTime, success Bool, latency_ms Float64, error_msg String, bytes_sent Int64 DEFAULT 0, bytes_recv Int64 DEFAULT 0, ttfb_ms Float64 DEFAULT 0.0, dns_time_ms Float64 DEFAULT 0.0, connect_time_ms Float64 DEFAULT 0.0, hops Int32 DEFAULT 0, jitter_ms Float64 DEFAULT 0.0, bit_error_rate Float64 DEFAULT 0.0, packet_loss Float64 DEFAULT 0.0) ENGINE = MergeTree() ORDER BY (timestamp, task_id);`,
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
			{ID: uuid.New().String(), Name: "HTTPBin GET", Target: "https://httpbin.org/get", IsDefault: true},
			{ID: uuid.New().String(), Name: "Random User API", Target: "https://randomuser.me/api/", IsDefault: true},
			{ID: uuid.New().String(), Name: "GitHub API Status", Target: "https://api.github.com/zen", IsDefault: true},
			{ID: uuid.New().String(), Name: "CoinDesk Bitcoin Price", Target: "https://api.coindesk.com/v1/bpi/currentprice.json", IsDefault: true},
			{ID: uuid.New().String(), Name: "Sample PCAP (FlowTest)", Target: "https://raw.githubusercontent.com/CESNET/FlowTest/master/test/testbed/generator/templates/pcap/dns.pcap", IsDefault: true},
		}
		tx, err := db.Begin()
		if err != nil {
			return nil, fmt.Errorf("failed to begin transaction for seeding sources: %w", err)
		}
		stmt, err := tx.Prepare(`INSERT INTO sources (id, name, target, is_default) VALUES (?, ?, ?, ?)`)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to prepare statement for seeding sources: %w", err)
		}
		defer stmt.Close()

		for _, s := range defaultSources {
			_, err := stmt.Exec(s.ID, s.Name, s.Target, s.IsDefault)
			if err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to execute insert for source %s: %w", s.Name, err)
			}
		}

		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("failed to commit transaction for seeding sources: %w", err)
		}
	}

	return &DB{db}, nil
}

func (db *DB) GetEnabledTasks() ([]models.Task, error) {
	rows, err := db.Query(`SELECT id, type, target, interval, config, enabled FROM tasks FINAL WHERE enabled = 1`)
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
	_, err := db.Exec(`INSERT INTO results (id, task_id, timestamp, success, latency_ms, error_msg, bytes_sent, bytes_recv, ttfb_ms, dns_time_ms, connect_time_ms, hops, jitter_ms, bit_error_rate, packet_loss) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, r.TaskID, r.Timestamp, r.Success, r.LatencyMs, r.ErrorMsg, r.BytesSent, r.BytesRecv, r.TTFBMs, r.DNSTimeMs, r.ConnectTimeMs, r.Hops, r.JitterMs, r.BitErrorRate, r.PacketLoss)
	return err
}

func (db *DB) GetRecentResults(limit int) ([]models.TaskResult, error) {
	rows, err := db.Query(`SELECT id, task_id, timestamp, success, latency_ms, error_msg, bytes_sent, bytes_recv, ttfb_ms, dns_time_ms, connect_time_ms, hops, jitter_ms, bit_error_rate, packet_loss FROM results ORDER BY timestamp DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.TaskResult
	for rows.Next() {
		var r models.TaskResult
		err := rows.Scan(&r.ID, &r.TaskID, &r.Timestamp, &r.Success, &r.LatencyMs, &r.ErrorMsg, &r.BytesSent, &r.BytesRecv, &r.TTFBMs, &r.DNSTimeMs, &r.ConnectTimeMs, &r.Hops, &r.JitterMs, &r.BitErrorRate, &r.PacketLoss)
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

func (db *DB) DeleteTask(id string) error {
	// For ReplacingMergeTree, deleting is tricky. We can insert a disabled version to "delete" it
	// from the active runner's perspective.
	_, err := db.Exec(`INSERT INTO tasks (id, type, target, interval, config, enabled) VALUES (?, '', '', 0, '', 0)`, id)
	return err
}

func (db *DB) GetAggregatedMetrics() (models.AggregatedMetrics, error) {
	var m models.AggregatedMetrics
	err := db.QueryRow(`
		SELECT
			count(),
			countIf(success = 1),
			countIf(success = 0),
			avg(latency_ms),
			sum(bytes_sent),
			sum(bytes_recv)
		FROM results
	`).Scan(&m.TotalTests, &m.TotalSuccesses, &m.TotalFailures, &m.AvgLatencyMs, &m.TotalBytesSent, &m.TotalBytesRecv)

	if err != nil {
		// Ignore null errors on empty db
		return m, nil
	}
	return m, nil
}

func (db *DB) GetSources() ([]models.Source, error) {
	rows, err := db.Query(`SELECT id, name, target, is_default FROM sources FINAL`)
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
