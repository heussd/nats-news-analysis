package timeseries

import (
	"database/sql"
	"fmt"

	"github.com/heussd/nats-news-analysis/internal/ngrams"
	"github.com/heussd/nats-news-analysis/pkg/utils"
	_ "github.com/lib/pq"
)

var (
	host     = utils.GetEnv("POSTGRES_HOST", "postgresql")
	port     = 5432
	user     = utils.GetEnv("POSTGRES_USER", "postgres")
	password = utils.GetEnv("POSTGRES_PASSWORD", "mysecretpassword")
	dbname   = utils.GetEnv("POSTGRES_DB", "ngrams")
	sslmode  = utils.GetEnv("PGSSLMODE", "disable")
)

var db *sql.DB

func ensureDatabaseExists() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, password, "postgres", sslmode)
	adminDB, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	defer adminDB.Close()

	if err := adminDB.Ping(); err != nil {
		return err
	}

	var exists bool
	err = adminDB.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbname).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		if _, err := adminDB.Exec(fmt.Sprintf("CREATE DATABASE %q", dbname)); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	if err := ensureDatabaseExists(); err != nil {
		panic(err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS ngrams (
		id SERIAL PRIMARY KEY,
		words TEXT NOT NULL,
		source TEXT NOT NULL,
		count INTEGER NOT NULL,
		language TEXT NOT NULL,
		n_gram INTEGER NOT NULL,
		timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`)
	if err != nil {
		panic(err)
	}

	// Index the timestamp for the time-window range scans used by trend queries.
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_ngrams_timestamp ON ngrams (timestamp)`)
	if err != nil {
		panic(err)
	}

	// Composite index to support grouping/lookups of phrases by (words, n_gram).
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_ngrams_words_ngram ON ngrams (words, n_gram)`)
	if err != nil {
		panic(err)
	}
}

// AddTimeSeriesData bulk-loads n-grams using PostgreSQL's COPY protocol, which
// is the fastest way to insert large volumes of rows. All rows are streamed
// inside a single transaction, so the load is atomic: either every row is
// committed or none are.
func AddTimeSeriesData(data []ngrams.NGram) error {
	if len(data) == 0 {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // no-op once Commit succeeds

	stmt, err := tx.Prepare("COPY ngrams (words, count, n_gram, source, language, timestamp) FROM STDIN")
	if err != nil {
		return err
	}

	for _, ngram := range data {
		if _, err := stmt.Exec(
			ngram.Words,
			ngram.Count,
			ngram.NGram,
			ngram.Source,
			ngram.Language,
			ngram.Timestamp,
		); err != nil {
			return err
		}
	}

	// A final Exec with no arguments flushes the buffered COPY data.
	if _, err := stmt.Exec(); err != nil {
		return err
	}

	if err := stmt.Close(); err != nil {
		return err
	}

	return tx.Commit()
}
