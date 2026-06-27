package repository

import "github.com/traP-jp/hackathon26spring_05/Qpid/domain"

// 計算済みのおすすめユーザーを扱う
type SuggestionRepository interface {
	// おすすめユーザーを取得する
	ListSuggestions(username string, limit int) ([]domain.Suggestion, error)
	// 類似度を保存する
	UpsertSimilarity(usernameA, usernameB string, similarity float64) error
}
