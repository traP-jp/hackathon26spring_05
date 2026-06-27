package repository

import "github.com/traP-jp/hackathon26spring_05/Qpid/domain"

// ユーザー間の LIKE/NOPE を永続化する
type ActionRepository interface {
	// ユーザーを LIKE する
	Like(fromUsername, toUsername string) error
	// ユーザーを NOPE する
	Nope(fromUsername, toUsername string) error
	// LIKE したユーザーを取得する
	ListLikedUsers(username string) ([]domain.UserSummary, error)
	// LIKE してくれたユーザーを取得する
	ListUsersWhoLiked(username string) ([]domain.UserSummary, error)
	// アクション済みか確認する
	HasAction(fromUsername, toUsername string) (bool, error)
}
