package repository

import "github.com/traP-jp/hackathon26spring_05/Qpid/domain"

// ユーザープロフィールを永続化する
type UserRepository interface {
	// ユーザーを作成する
	Create(user domain.User) error
	// ユーザーを取得する
	FindByUsername(username string) (*domain.User, error)
	// プロフィールを更新する
	UpdateProfile(username string, user domain.User) error
	// ユーザーの存在を確認する
	Exists(username string) (bool, error)
}
