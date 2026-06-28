package infrastructure

import "github.com/traP-jp/hackathon26spring_05/Qpid/domain"

// おすすめユーザーを取得する。
func (r *repositoryImpl) ListSuggestions(username string, limit int) ([]domain.Suggestion, error) {
	var users []domain.Suggestion
	err := r.db.Select(&users, `
		SELECT username
		FROM users
		WHERE username <> ?
			AND NOT EXISTS (
				SELECT 1
				FROM actions
				WHERE actions.from_username = ?
					AND actions.to_username = users.username
			)
		ORDER BY RAND()
		LIMIT ?
	`, username, username, limit)
	if err != nil {
		return nil, err
	}
	for user := range users {
		users[user].Similarity = 0.5
	}
	return users, nil
}

// 類似度を保存する。
func (r *repositoryImpl) UpsertSimilarity(usernameA, usernameB string, similarity float64) error {
	return nil
}
