package repository

import "time"

// OAuth2 state の発行と消費を管理する
type OAuthStateRepository interface {
	// OAuth2 state を作成する
	CreateOAuthState(state string, expiresAt time.Time) error
	// OAuth2 state を消費する
	ConsumeOAuthState(state string) error
}
