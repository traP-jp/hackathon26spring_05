package infrastructure

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/moznion/go-optional"
	"github.com/traP-jp/hackathon26spring_05/Qpid/domain"
)

type userRow struct {
	Username     string         `db:"username"`
	Major        sql.NullString `db:"major"`
	Hometown     sql.NullString `db:"hometown"`
	LikeTopic    sql.NullString `db:"like_topic"`
	LikeValue    sql.NullString `db:"like_value"`
	DislikeTopic sql.NullString `db:"dislike_topic"`
	DislikeValue sql.NullString `db:"dislike_value"`
	Bio          sql.NullString `db:"bio"`
}

// ユーザーを作成する。
func (r *repositoryImpl) CreateUser(user domain.User) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := insertUserProfile(tx, user.Username, user); err != nil {
		return err
	}
	if err := r.updateUserTags(tx, user.Username, user.Tags); err != nil {
		return err
	}
	if err := r.updateUserTechnologies(tx, user.Username, user.Technologies); err != nil {
		return err
	}

	return tx.Commit()
}

// ユーザーを取得する。
func (r *repositoryImpl) FindUserByUsername(username string) (*domain.User, error) {
	var row userRow
	query := `
		SELECT username, major, hometown, like_topic, like_value, dislike_topic, dislike_value, bio
		FROM users
		WHERE username = ?`
	if err := r.db.Get(&row, query, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	var hasIcon bool
	queryIcon := `SELECT EXISTS (SELECT 1 FROM icons WHERE username = ?)`
	if err := r.db.Get(&hasIcon, queryIcon, username); err != nil {
		return nil, err
	}

	var tags []string
	queryTags := `SELECT name FROM tags WHERE username = ?`
	if err := r.db.Select(&tags, queryTags, username); err != nil {
		return nil, err
	}

	var tools []string
	queryTools := `SELECT name FROM tools WHERE username = ?`
	if err := r.db.Select(&tools, queryTools, username); err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:      row.Username,
		HasIcon:       hasIcon,
		Major:         convertNullString(row.Major),
		Hometown:      convertNullString(row.Hometown),
		Bio:           convertNullString(row.Bio),
		FavoriteTopic: convertTopic(row.LikeTopic, row.LikeValue),
		DislikedTopic: convertTopic(row.DislikeTopic, row.DislikeValue),
		Tags:          tags,
		Technologies:  tools,
		Affiliations:  []domain.UserAffiliation{},
	}
	return user, nil
}

// プロフィールを更新する。
func (r *repositoryImpl) UpdateUser(username string, user domain.User) error {
	// 更新クエリの実行
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM users WHERE username = ?`, username); err != nil {
		return err
	}

	if err := insertUserProfile(tx, username, user); err != nil {
		return err
	}

	// 関連テーブルの更新
	if err := r.updateUserTags(tx, username, user.Tags); err != nil {
		return err
	}
	if err := r.updateUserTechnologies(tx, username, user.Technologies); err != nil {
		return err
	}

	return tx.Commit()
}

func insertUserProfile(tx *sqlx.Tx, username string, user domain.User) error {
	likeTopic, likeValue := toNullTopic(user.FavoriteTopic)
	dislikeTopic, dislikeValue := toNullTopic(user.DislikedTopic)

	query := `
		INSERT INTO users
		(username, major, hometown, like_topic, like_value, dislike_topic, dislike_value, bio)
		VALUES(?,?,?,?,?,?,?,?)`

	_, err := tx.Exec(
		query,
		username,
		toNullString(user.Major),
		toNullString(user.Hometown),
		likeTopic,
		likeValue,
		dislikeTopic,
		dislikeValue,
		toNullString(user.Bio),
	)
	return err
}

// ユーザーの存在を確認する。
func (r *repositoryImpl) IsUserExists(username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE username = ?)`
	if err := r.db.Get(&exists, query, username); err != nil {
		return false, err
	}
	return exists, nil
}

// tagsテーブルを更新する。
func (r *repositoryImpl) updateUserTags(tx *sqlx.Tx, username string, tags []string) error {
	if _, err := tx.Exec(`DELETE FROM tags WHERE username = ?`, username); err != nil {
		return err
	}
	if len(tags) == 0 {
		return nil
	}

	//バルクインサート用のクエリを構築
	valueStrings := make([]string, 0, len(tags))
	valueArgs := make([]interface{}, 0, len(tags)*2)
	for _, tag := range tags {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, username, tag)
	}

	query := `INSERT INTO tags (username, name) VALUES ` + strings.Join(valueStrings, ",")

	// 一括実行
	_, err := tx.Exec(query, valueArgs...)
	return err
}

// toolsテーブルを更新する。
func (r *repositoryImpl) updateUserTechnologies(tx *sqlx.Tx, username string, techs []string) error {
	if _, err := tx.Exec(`DELETE FROM tools WHERE username = ?`, username); err != nil {
		return err
	}
	if len(techs) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(techs))
	valueArgs := make([]interface{}, 0, len(techs)*2)
	for _, tech := range techs {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, username, tech)
	}

	query := `INSERT INTO tools (username, name) VALUES ` + strings.Join(valueStrings, ",")

	_, err := tx.Exec(query, valueArgs...)
	return err
}

// 以下ヘルパー関数
func convertNullString(ns sql.NullString) optional.Option[string] {
	if ns.Valid {
		return optional.Some(ns.String)
	}
	return optional.None[string]()
}

func convertTopic(topic, value sql.NullString) optional.Option[domain.TopicAndValue] {
	if topic.Valid && value.Valid {
		return optional.Some(domain.TopicAndValue{
			Topic: topic.String,
			Value: value.String,
		})
	}
	return optional.None[domain.TopicAndValue]()
}

func toNullString(opt optional.Option[string]) sql.NullString {
	if opt.IsSome() {
		return sql.NullString{
			String: opt.Unwrap(),
			Valid:  true,
		}
	}

	return sql.NullString{
		Valid: false,
	}
}

func toNullTopic(opt optional.Option[domain.TopicAndValue]) (sql.NullString, sql.NullString) {
	if opt.IsSome() {
		topicAndValue := opt.Unwrap()
		return sql.NullString{String: topicAndValue.Topic, Valid: true},
			sql.NullString{String: topicAndValue.Value, Valid: true}
	}

	return sql.NullString{Valid: false}, sql.NullString{Valid: false}
}
