package repository

// アプリケーションが必要とする Repository 群をまとめた interface。
type Repository interface {
	UserRepository
	ActionRepository
	SuggestionRepository
	PrecomputedProfileRepository
	IconRepository
	AuthSessionRepository
	OAuthStateRepository
}
