package repository

import "github.com/traP-jp/hackathon26spring_05/Qpid/domain"

// アプリセッションと OAuth2 トークンの対応を管理する
type AuthSessionRepository interface {
	// 認証セッションを作成する
	CreateAuthSession(sessionID string, session domain.AuthSession) error
	// 認証セッションを取得する
	FindAuthSession(sessionID string) (*domain.AuthSession, error)
	// 認証セッションを削除する
	DeleteAuthSession(sessionID string) error
}
