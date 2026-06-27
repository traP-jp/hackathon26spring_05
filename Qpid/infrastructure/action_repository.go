package infrastructure

import "github.com/traP-jp/hackathon26spring_05/Qpid/domain"

// ユーザーを LIKE する。
func (r *repositoryImpl) LikeUser(fromUsername, toUsername string) error {
	return nil
}

// ユーザーを NOPE する。
func (r *repositoryImpl) NopeUser(fromUsername, toUsername string) error {
	return nil
}

// LIKE したユーザーを取得する。
func (r *repositoryImpl) ListLikedUsers(username string) ([]domain.UserSummary, error) {
	return nil, nil
}

// LIKE してくれたユーザーを取得する。
func (r *repositoryImpl) ListUsersWhoLiked(username string) ([]domain.UserSummary, error) {
	return nil, nil
}

// アクション済みか確認する。
func (r *repositoryImpl) IsActionExists(fromUsername, toUsername string) (bool, error) {
	return false, nil
}
