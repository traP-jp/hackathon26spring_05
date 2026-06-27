package domain

import "time"

// アプリケーションのセッションに紐づく OAuth2 認証情報。
type AuthSession struct {
	Username     string
	AccessToken  string
	RefreshToken *string
	ExpiresAt    time.Time
}
