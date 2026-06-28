package infrastructure

import (
	"sort"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/moznion/go-optional"
	"github.com/traP-jp/hackathon26spring_05/Qpid/domain"
)

func TestFindUserByUsernameScansProfile(t *testing.T) {
	repo, db := newUserRepositoryTestDB(t)

	mustExec(t, db, `
		INSERT INTO users
		(username, major, hometown, like_topic, like_value, dislike_topic, dislike_value, bio)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		"Arina", "math", "Tokyo", "editor", "vim", "food", "tomato", "hello",
	)
	mustExec(t, db, `INSERT INTO icons (username, icon, mime_type) VALUES (?, ?, ?)`, "Arina", []byte("icon"), "image/png")
	mustExec(t, db, `INSERT INTO tags (username, name) VALUES (?, ?)`, "Arina", "backend")
	mustExec(t, db, `INSERT INTO tags (username, name) VALUES (?, ?)`, "Arina", "go")
	mustExec(t, db, `INSERT INTO tools (username, name) VALUES (?, ?)`, "Arina", "docker")
	mustExec(t, db, `INSERT INTO tools (username, name) VALUES (?, ?)`, "Arina", "mysql")

	user, err := repo.FindUserByUsername("Arina")
	if err != nil {
		t.Fatalf("FindUserByUsername returned error: %v", err)
	}
	if user == nil {
		t.Fatal("FindUserByUsername returned nil")
	}
	if user.Username != "Arina" {
		t.Fatalf("Username = %q, want %q", user.Username, "Arina")
	}
	if !user.HasIcon {
		t.Fatal("HasIcon = false, want true")
	}
	assertOptionalString(t, user.Major, "math")
	assertOptionalString(t, user.Hometown, "Tokyo")
	assertOptionalString(t, user.Bio, "hello")
	assertOptionalTopic(t, user.FavoriteTopic, "editor", "vim")
	assertOptionalTopic(t, user.DislikedTopic, "food", "tomato")
	assertStrings(t, user.Tags, []string{"backend", "go"})
	assertStrings(t, user.Technologies, []string{"docker", "mysql"})
}

func TestIsUserExists(t *testing.T) {
	repo, db := newUserRepositoryTestDB(t)
	mustExec(t, db, `INSERT INTO users (username) VALUES (?)`, "Arina")

	exists, err := repo.IsUserExists("Arina")
	if err != nil {
		t.Fatalf("IsUserExists returned error: %v", err)
	}
	if !exists {
		t.Fatal("IsUserExists = false, want true")
	}

	exists, err = repo.IsUserExists("missing")
	if err != nil {
		t.Fatalf("IsUserExists returned error for missing user: %v", err)
	}
	if exists {
		t.Fatal("IsUserExists = true for missing user, want false")
	}
}

func TestCreateUserPersistsProfile(t *testing.T) {
	repo, _ := newUserRepositoryTestDB(t)

	err := repo.CreateUser(domain.User{
		Username:      "Arina",
		Major:         optional.Some("math"),
		Hometown:      optional.Some("Tokyo"),
		Tags:          []string{"backend", "go"},
		Technologies:  []string{"docker", "mysql"},
		Bio:           optional.Some("hello"),
		FavoriteTopic: optional.Some(domain.TopicAndValue{Topic: "editor", Value: "vim"}),
		DislikedTopic: optional.Some(domain.TopicAndValue{Topic: "food", Value: "tomato"}),
	})
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	user, err := repo.FindUserByUsername("Arina")
	if err != nil {
		t.Fatalf("FindUserByUsername returned error: %v", err)
	}
	if user == nil {
		t.Fatal("FindUserByUsername returned nil")
	}
	assertOptionalString(t, user.Major, "math")
	assertOptionalString(t, user.Hometown, "Tokyo")
	assertOptionalString(t, user.Bio, "hello")
	assertOptionalTopic(t, user.FavoriteTopic, "editor", "vim")
	assertOptionalTopic(t, user.DislikedTopic, "food", "tomato")
	assertStrings(t, user.Tags, []string{"backend", "go"})
	assertStrings(t, user.Technologies, []string{"docker", "mysql"})
}

func TestUpdateUserReplacesTagsAndTools(t *testing.T) {
	repo, _ := newUserRepositoryTestDB(t)

	err := repo.CreateUser(domain.User{
		Username:     "Arina",
		Tags:         []string{"old-tag"},
		Technologies: []string{"old-tool"},
	})
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	err = repo.UpdateUser("Arina", domain.User{
		Tags:         []string{"new-tag"},
		Technologies: []string{"new-tool"},
	})
	if err != nil {
		t.Fatalf("UpdateUser returned error: %v", err)
	}

	user, err := repo.FindUserByUsername("Arina")
	if err != nil {
		t.Fatalf("FindUserByUsername returned error: %v", err)
	}
	if user == nil {
		t.Fatal("FindUserByUsername returned nil")
	}
	assertStrings(t, user.Tags, []string{"new-tag"})
	assertStrings(t, user.Technologies, []string{"new-tool"})
}

func newUserRepositoryTestDB(t *testing.T) (*repositoryImpl, *sqlx.DB) {
	t.Helper()

	db, err := sqlx.Open("sqlite3", ":memory:")
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
			username TEXT NOT NULL PRIMARY KEY,
			created_at TIMESTAMP,
			major TEXT,
			hometown TEXT,
			like_topic TEXT,
			like_value TEXT,
			dislike_topic TEXT,
			dislike_value TEXT,
			usual_situation TEXT,
			bio TEXT
		)`,
		`CREATE TABLE icons (
			username TEXT NOT NULL PRIMARY KEY,
			icon BLOB NOT NULL,
			mime_type TEXT NOT NULL
		)`,
		`CREATE TABLE tags (
			username TEXT NOT NULL,
			name TEXT NOT NULL,
			PRIMARY KEY (username, name)
		)`,
		`CREATE TABLE tools (
			username TEXT NOT NULL,
			name TEXT NOT NULL,
			PRIMARY KEY (username, name)
		)`,
	}

	for _, query := range schema {
		mustExec(t, db, query)
	}

	return NewRepository(db), db
}

func mustExec(t *testing.T, db *sqlx.DB, query string, args ...any) {
	t.Helper()

	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("failed to exec query: %v", err)
	}
}

func assertOptionalString(t *testing.T, got optional.Option[string], want string) {
	t.Helper()

	if got.IsNone() {
		t.Fatalf("got none, want %q", want)
	}
	if got.Unwrap() != want {
		t.Fatalf("got %q, want %q", got.Unwrap(), want)
	}
}

func assertOptionalTopic(t *testing.T, got optional.Option[domain.TopicAndValue], wantTopic, wantValue string) {
	t.Helper()

	if got.IsNone() {
		t.Fatalf("got none, want topic=%q value=%q", wantTopic, wantValue)
	}
	topicAndValue := got.Unwrap()
	if topicAndValue.Topic != wantTopic || topicAndValue.Value != wantValue {
		t.Fatalf("got topic=%q value=%q, want topic=%q value=%q", topicAndValue.Topic, topicAndValue.Value, wantTopic, wantValue)
	}
}

func assertStrings(t *testing.T, got, want []string) {
	t.Helper()

	gotCopy := append([]string(nil), got...)
	wantCopy := append([]string(nil), want...)
	sort.Strings(gotCopy)
	sort.Strings(wantCopy)
	if len(gotCopy) != len(wantCopy) {
		t.Fatalf("got %v, want %v", gotCopy, wantCopy)
	}
	for i := range gotCopy {
		if gotCopy[i] != wantCopy[i] {
			t.Fatalf("got %v, want %v", gotCopy, wantCopy)
		}
	}
}
