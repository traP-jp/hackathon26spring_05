package infrastructure

import "github.com/traP-jp/hackathon26spring_05/Qpid/domain"

// ユーザーを作成する。
func (r *repositoryImpl) Create(user domain.User) error {
	return nil
}

// ユーザーを取得する。
func (r *repositoryImpl) FindByUsername(username string) (*domain.User, error) {
	return nil, nil
}

// プロフィールを更新する。
func (r *repositoryImpl) UpdateProfile(username string, user domain.User) error {
	return nil
}

// ユーザーの存在を確認する。
func (r *repositoryImpl) Exists(username string) (bool, error) {
	return false, nil
}
