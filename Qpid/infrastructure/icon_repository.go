package infrastructure

import (
	"github.com/traP-jp/hackathon26spring_05/Qpid/domain"
)

// アイコン画像を保存する。
func (r *repositoryImpl) SaveIcon(username string, icon domain.Icon) error {
	return nil
}

// アイコン画像を取得する。
func (r *repositoryImpl) FindIconByUsername(username string) (*domain.Icon, error) {
	var icon *domain.Icon
	err := r.db.Select(&icon,
		"SELECT icon LONGBLOB, mime_type FROM icons WHERE username = ? LIMIT 1",
		username,
	)
	return icon, err
}

// アイコン画像を削除する。
func (r *repositoryImpl) DeleteIcon(username string) error {
	return nil
}
