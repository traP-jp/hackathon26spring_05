package infrastructure

import "github.com/traP-jp/hackathon26spring_05/Qpid/domain"

// おすすめユーザーを取得する。
func (r *repositoryImpl) ListSuggestions(username string, limit int) ([]domain.Suggestion, error) {
	return nil, nil
}

// 類似度を保存する。
func (r *repositoryImpl) UpsertSimilarity(usernameA, usernameB string, similarity float64) error {
	return nil
}
