package infrastructure

import (
	"database/sql"
	"sync"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

var registerSuggestionRepositoryTestDBOnce sync.Once

func TestListSuggestionsExcludesSelfAndActionedUsers(t *testing.T) {
	repo, db := newSuggestionRepositoryTestDB(t)
	mustExec(t, db, `INSERT INTO users (username) VALUES (?)`, "arina")
	mustExec(t, db, `INSERT INTO users (username) VALUES (?)`, "bob")
	mustExec(t, db, `INSERT INTO users (username) VALUES (?)`, "charlie")
	mustExec(t, db, `INSERT INTO actions (id, from_username, to_username, status) VALUES (?, ?, ?, ?)`, "action-1", "arina", "bob", actionStatusLike)

	suggestions, err := repo.ListSuggestions("arina", 20)
	if err != nil {
		t.Fatalf("ListSuggestions returned error: %v", err)
	}

	if len(suggestions) != 1 {
		t.Fatalf("got %d suggestions, want 1: %v", len(suggestions), suggestions)
	}
	if suggestions[0].Username != "charlie" {
		t.Fatalf("got suggestion %q, want %q", suggestions[0].Username, "charlie")
	}
	if suggestions[0].Similarity != 0.5 {
		t.Fatalf("got similarity %v, want 0.5", suggestions[0].Similarity)
	}
}

func newSuggestionRepositoryTestDB(t *testing.T) (*repositoryImpl, *sqlx.DB) {
	t.Helper()

	registerSuggestionRepositoryTestDBOnce.Do(func() {
		sql.Register("sqlite3_suggestion_repository_test", &sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				return conn.RegisterFunc("RAND", func() float64 { return 0.5 }, true)
			},
		})
	})

	db, err := sqlx.Open("sqlite3_suggestion_repository_test", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Fatalf("failed to close test db: %v", err)
		}
	})

	schema := []string{
		`CREATE TABLE users (
			username TEXT NOT NULL PRIMARY KEY
		)`,
		`CREATE TABLE actions (
			id TEXT NOT NULL PRIMARY KEY,
			from_username TEXT NOT NULL,
			to_username TEXT NOT NULL,
			status INTEGER NOT NULL
		)`,
	}

	for _, query := range schema {
		mustExec(t, db, query)
	}

	return NewRepository(db), db
}
