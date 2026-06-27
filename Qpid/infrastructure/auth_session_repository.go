package infrastructure

import "github.com/traP-jp/hackathon26spring_05/Qpid/domain"

// 認証セッションを作成する。
func (r *repositoryImpl) CreateAuthSession(sessionID string, session domain.AuthSession) error {
	return nil
}

// 認証セッションを取得する。
func (r *repositoryImpl) FindAuthSession(sessionID string) (*domain.AuthSession, error) {
	return nil, nil
}

// 認証セッションを削除する。
func (r *repositoryImpl) DeleteAuthSession(sessionID string) error {
	return nil
}
