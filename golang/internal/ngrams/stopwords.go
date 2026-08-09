package ngrams

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/heussd/nats-news-analysis/pkg/utils"
)

var (
	db          *sql.DB = nil
	lastRefresh         = time.Time{}
)

var (
	host     = utils.GetEnv("POSTGRES_HOST", "postgresql")
	port     = 5432
	user     = utils.GetEnv("POSTGRES_USER", "postgres")
	password = utils.GetEnv("POSTGRES_PASSWORD", "mysecretpassword")
	dbname   = utils.GetEnv("POSTGRES_DB", "ngrams2")
	sslmode  = utils.GetEnv("PGSSLMODE", "disable")
)

func initConnection() (db *sql.DB, err error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	if db, err = sql.Open("postgres", psqlInfo); err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS stopwords (
		id SERIAL PRIMARY KEY,
		stopword TEXT NOT NULL,
		exact BOOLEAN NOT NULL DEFAULT TRUE
	)`); err != nil {
		return nil, fmt.Errorf("failed to create stopwords table: %w", err)
	}

	return db, nil
}

func retrieveStopwords() (stopWordsExact []string, stopWordsAny []string, err error) {
	if db == nil {
		var err error
		db, err = initConnection()
		if err != nil {
			panic(err)
		}
	}

	rows, err := db.Query(`SELECT stopword, exact FROM stopwords`)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve stopwords: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var stopword string
		var exact bool
		if err := rows.Scan(&stopword, &exact); err != nil {
			return nil, nil, fmt.Errorf("failed to scan stopword: %w", err)
		}
		if exact {
			stopWordsExact = append(stopWordsExact, stopword)
		} else {
			stopWordsAny = append(stopWordsAny, stopword)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error iterating over stopwords: %w", err)
	}

	return stopWordsExact, stopWordsAny, nil
}

var (
	cachedStopWordsExact []string
	cachedStopWordsAny   []string
)

func CachedStopwords() ([]string, []string) {
	if time.Since(lastRefresh) > 5*time.Minute {
		stopWordsExact, stopWordsAny, err := retrieveStopwords()
		if err != nil {
			panic(err)
		}
		cachedStopWordsExact = stopWordsExact
		cachedStopWordsAny = stopWordsAny
		lastRefresh = time.Now()
	}

	return cachedStopWordsExact, cachedStopWordsAny
}

func AnyStopwordsRegex() *regexp.Regexp {
	_, stopWordsAny := CachedStopwords()
	if len(stopWordsAny) == 0 {
		return nil
	}
	pattern := `(^|\s)(` + strings.Join(stopWordsAny, "|") + `)(\s|$)`
	return regexp.MustCompile(pattern)
}
