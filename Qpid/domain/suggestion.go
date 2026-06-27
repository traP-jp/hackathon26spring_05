package domain

// おすすめユーザーと類似度
type Suggestion struct {
	// ユーザー ID
	Username string
	// 対象ユーザーとの類似度
	Similarity float64
}
