package infrastructure

import "github.com/traP-jp/hackathon26spring_05/Qpid/domain"

// アイコン画像を保存する。
func (r *repositoryImpl) SaveIcon(username string, icon domain.Icon) error {
	return nil
}

// アイコン画像を取得する。
func (r *repositoryImpl) FindIconByUsername(username string) (*domain.Icon, error) {
	return nil, nil
}

// アイコン画像を削除する。
func (r *repositoryImpl) DeleteIcon(username string) error {
	return nil
}
