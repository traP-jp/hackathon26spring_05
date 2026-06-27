package infrastructure

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
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
	if err := r.db.Select(&tags,queryTags,username);err!=nil{
		return nil,err
	}

	var tools []string
	queryTools :=`SELECT name FROM tools WHERE username = ?`
	if err := r.db.Select(&tools,queryTools,username);err!=nil{
		return nil,err
	}

	user := &domain.User{
		Username: row.Username,
		Major:    convertNullString(row.Major),
		Hometown: convertNullString(row.Hometown),
		Bio:      convertNullString(row.Bio),
		FavoriteTopic: convertTopic(row.LikeTopic, row.LikeValue),
		DislikedTopic: convertTopic(row.DislikeTopic, row.DislikeValue),
		Tags: tags,
		Technologies: tools,
		Affiliations:  []domain.UserAffiliation{},
	}
	return user,nil
}

// プロフィールを更新する。
func (r *repositoryImpl) UpdateUser(username string, user domain.User) error {
	major := sql.NullString{Valid: false}
	if user.Major.IsSome() {
		major = sql.NullString{String: user.Major.Unwrap(), Valid: true}
	}

	hometown := sql.NullString{Valid: false}
	if user.Hometown.IsSome() {
		hometown = sql.NullString{String: user.Hometown.Unwrap(), Valid: true}
	}

	bio := sql.NullString{Valid: false}
	if user.Bio.IsSome() {
		bio = sql.NullString{String: user.Bio.Unwrap(), Valid: true}
	}

	// Topic系は構造体の中身を確認してセット
	likeTopic, likeValue := sql.NullString{Valid: false}, sql.NullString{Valid: false}
	if user.FavoriteTopic.IsSome() {
		t := user.FavoriteTopic.Unwrap()
		likeTopic = sql.NullString{String: t.Topic, Valid: true}
		likeValue = sql.NullString{String: t.Value, Valid: true}
	}

	dislikeTopic, dislikeValue := sql.NullString{Valid: false}, sql.NullString{Valid: false}
	if user.DislikedTopic.IsSome() {
		t := user.DislikedTopic.Unwrap()
		dislikeTopic = sql.NullString{String: t.Topic, Valid: true}
		dislikeValue = sql.NullString{String: t.Value, Valid: true}
	}

	// 更新クエリの実行
	tx, err := r.db.Beginx()
	if err != nil {
			return err
	}
	defer tx.Rollback()

	query := `
		UPDATE users 
		SET major = ?, hometown = ?, like_topic = ?, like_value = ?, 
		    dislike_topic = ?, dislike_value = ?, bio = ?
		WHERE username = ?`

	_, err = tx.Exec(query, major, hometown, likeTopic, likeValue, dislikeTopic, dislikeValue, bio, username)
	if err != nil {
		return err
	}

	// 関連テーブルの更新
	if err := r.updateUserTags(tx,username, user.Tags); err != nil {
		return err
	}
	if err := r.updateUserTechnologies(tx,username, user.Technologies); err != nil {
		return err
	}

	return nil
}

// ユーザーの存在を確認する。
func (r *repositoryImpl) IsUserExists(username string) (bool, error) {
	var row User
	query := `SELECT 1 FROM users WHERE username = ?`
	if err := r.db.Get(&row,query,username); err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return false, nil
		}
		return false,err
	}
	return true,nil

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
	
	// 3. 一括実行
	_, err := tx.Exec(query, valueArgs...)
	return err
}

// technologiesテーブルを更新する。
func (r *repositoryImpl) updateUserTechnologies(tx *sqlx.Tx, username string, techs []string) error {
	if _, err := tx.Exec(`DELETE FROM technologies WHERE username = ?`, username); err != nil {
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

	query := `INSERT INTO technologies (username, name) VALUES ` + strings.Join(valueStrings, ",")
	
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