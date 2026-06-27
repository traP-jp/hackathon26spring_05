package repository

import "github.com/traP-jp/hackathon26spring_05/Qpid/domain"

// ユーザープロフィールを永続化する
type UserRepository interface {
	// ユーザーを作成する
	CreateUser(user domain.User) error
	// ユーザーを取得する
	FindUserByUsername(username string) (*domain.User, error)
	// プロフィールを更新する
	UpdateUser(username string, user domain.User) error
	// ユーザーの存在を確認する
	IsUserExists(username string) (bool, error)
}
