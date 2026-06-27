package infrastructure

import "github.com/traP-jp/hackathon26spring_05/Qpid/domain"

const (
	actionStatusNope = 0
	actionStatusLike = 1
)

// ユーザーを LIKE する。
func (r *repositoryImpl) LikeUser(fromUsername, toUsername string) error {
	_, err := r.db.Exec(
		`INSERT INTO actions (id, from_username, to_username, status) VALUES (UUID(), ?, ?, ?)`,
		fromUsername, toUsername, actionStatusLike,
	)
	return err
}

// ユーザーを NOPE する。
func (r *repositoryImpl) NopeUser(fromUsername, toUsername string) error {
	_, err := r.db.Exec(
		`INSERT INTO actions (id, from_username, to_username, status) VALUES (UUID(), ?, ?, ?)`,
		fromUsername, toUsername, actionStatusNope,
	)
	return err
}

// LIKE したユーザーを取得する。
func (r *repositoryImpl) ListLikedUsers(username string) ([]domain.UserSummary, error) {
	var users []domain.UserSummary
	err := r.db.Select(&users,
		`SELECT to_username AS username FROM actions WHERE from_username = ? AND status = ?`,
		username, actionStatusLike,
	)
	return users, err
}

// LIKE してくれたユーザーを取得する。
func (r *repositoryImpl) ListUsersWhoLiked(username string) ([]domain.UserSummary, error) {
	var users []domain.UserSummary
	err := r.db.Select(&users,
		"SELECT from_username AS username FROM actions WHERE to_username = ? AND status = ?",
		username, actionStatusLike,
	)
	return users, err
}

// アクション済みか確認する。
func (r *repositoryImpl) IsActionExists(fromUsername, toUsername string) (bool, error) {
	return false, nil
}
