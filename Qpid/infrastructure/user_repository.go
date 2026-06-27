package infrastructure

import (
	"database/sql"
	"errors"

	//"github.com/jmoiron/sqlx"
	"github.com/moznion/go-optional"
	"github.com/traP-jp/hackathon26spring_05/Qpid/domain"
)

type User struct{
	Username string `db:"username"`
	CreatedAt sql.NullTime `db:"created_at"`
	Major sql.NullString `db:"major"`
	Hometown sql.NullString `db:"hometown"`
	LikeTopic sql.NullString `db:"like_topic"`
	LikeValue sql.NullString `db:"like_value"`
	DislikeTopic sql.NullString `db:"dislike_topic"`
	DislikeValue sql.NullString `db:"dislike_value"`
	UsualSituation sql.NullString `db:"usual_situation"`
	Bio sql.NullString `db:"bio"`
}

// ユーザーを作成する。
func (r *repositoryImpl) CreateUser(user domain.User) error {
	query := `INSERT INTO users (username) VALUES (?)`
	_, err := r.db.Exec(query,user.Username)
	return err
}

// ユーザーを取得する。
func (r *repositoryImpl) FindUserByUsername(username string) (*domain.User, error) {
	var row User
	query := `SELECT * FROM users WHERE username = ?`
	if err := r.db.Get(&row,query,username); err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return nil, nil
		}
		return nil,err
	}

	var hasIcon bool
	queryIcon := `SELECT 1 FROM icons WHERE username = ?`
	if err := r.db.Get(&hasIcon,queryIcon, username);err !=nil{
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		hasIcon = false
	}else {
		hasIcon = true
	}

	var tags []string
	queryTags := `SELECT name FROM tags WHERE username = ?`
	if err := r.db.Get(&tags,queryTags,username);err!=nil{
		return nil,err
	}

	user := &domain.User{
		Username: row.Username,
		Major:    convertNullString(row.Major),
		Hometown: convertNullString(row.Hometown),
		Bio:      convertNullString(row.Bio),
		FavoriteTopic: convertTopic(row.LikeTopic, row.LikeValue),
		DislikedTopic: convertTopic(row.DislikeTopic, row.DislikeValue),
	}
	return user,nil
}

// プロフィールを更新する。
func (r *repositoryImpl) UpdateUser(username string, user domain.User) error {
	return nil
}

// ユーザーの存在を確認する。
func (r *repositoryImpl) IsUserExists(username string) (bool, error) {
	return false, nil
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