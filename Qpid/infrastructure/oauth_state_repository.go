package infrastructure

import "time"

// OAuth2 state を作成する。
func (r *repositoryImpl) CreateOAuthState(state string, expiresAt time.Time) error {
	return nil
}

// OAuth2 state を消費する。
func (r *repositoryImpl) ConsumeOAuthState(state string) error {
	return nil
}
