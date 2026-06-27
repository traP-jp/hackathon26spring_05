package infrastructure

import (
	"database/sql"
	"errors"

	"github.com/moznion/go-optional"
	"github.com/traP-jp/hackathon26spring_05/Qpid/domain"
)

type precomputedUserRow struct {
	Username     string         `db:"username"`
	Major        string         `db:"major"`
	Hometown     sql.NullString `db:"hometown"`
	LikeTopic    sql.NullString `db:"like_topic"`
	LikeValue    sql.NullString `db:"like_value"`
	DislikeTopic sql.NullString `db:"dislike_topic"`
	DislikeValue sql.NullString `db:"dislike_value"`
	Tool         sql.NullString `db:"tool"`
	Bio          sql.NullString `db:"bio"`
}

// 事前計算済みプロフィールを取得する。
func (r *repositoryImpl) FindPrecomputedProfileByUsername(username string) (*domain.User, error) {
	var row precomputedUserRow
	err := r.db.Get(&row, "SELECT username, major, hometown, like_topic, like_value, dislike_topic, dislike_value, tool, bio FROM users WHERE username = ?", username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// タグを取得する
	tags := []string{}
	if err := r.db.Select(&tags, "SELECT name FROM tags WHERE username = ?", username); err != nil {
		return nil, err
	}

	// アイコンの有無を確認する
	var iconCount int
	if err := r.db.Get(&iconCount, "SELECT COUNT(*) FROM icons WHERE username = ?", username); err != nil {
		return nil, err
	}

	major := optional.None[string]()
	if row.Major != "" {
		major = optional.Some(row.Major)
	}

	user := &domain.User{
		Username:      row.Username,
		HasIcon:       iconCount > 0,
		Major:         major,
		Affiliations:  []domain.UserAffiliation{},
		Hometown:      optional.None[string](),
		Tags:          tags,
		Technologies:  []string{},
		Bio:           optional.None[string](),
		FavoriteTopic: optional.None[domain.TopicAndValue](),
		DislikedTopic: optional.None[domain.TopicAndValue](),
	}

	if row.Hometown.Valid {
		user.Hometown = optional.Some(row.Hometown.String)
	}
	if row.Bio.Valid {
		user.Bio = optional.Some(row.Bio.String)
	}
	if row.LikeTopic.Valid && row.LikeValue.Valid {
		user.FavoriteTopic = optional.Some(domain.TopicAndValue{
			Topic: row.LikeTopic.String,
			Value: row.LikeValue.String,
		})
	}
	if row.DislikeTopic.Valid && row.DislikeValue.Valid {
		user.DislikedTopic = optional.Some(domain.TopicAndValue{
			Topic: row.DislikeTopic.String,
			Value: row.DislikeValue.String,
		})
	}
	if row.Tool.Valid && row.Tool.String != "" {
		user.Technologies = []string{row.Tool.String}
	}

	return user, nil
}

// 事前計算済みユーザー名を一覧する。
func (r *repositoryImpl) ListPrecomputedUsernames() ([]string, error) {
	return nil, nil
}
