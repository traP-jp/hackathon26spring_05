package infrastructure

import (
	"github.com/traP-jp/hackathon26spring_05/Qpid/domain"
)

type iconRow struct {
	Blob     []byte `db:"icon"`
	MimeType string `db:"mime_type"`
}

// アイコン画像を保存する。
func (r *repositoryImpl) SaveIcon(username string, icon domain.Icon) error {
	return nil
}

// アイコン画像を取得する。
func (r *repositoryImpl) FindIconByUsername(username string) (*domain.Icon, error) {
	var icon iconRow
	err := r.db.Get(&icon,
		"SELECT icon, mime_type FROM icons WHERE username = ? LIMIT 1",
		username,
	)
	if err != nil {
		return nil, err
	}

	return &domain.Icon{
			Blob:     icon.Blob,
			MimeType: domain.IconMimeType(icon.MimeType),
		},
		nil
}

// アイコン画像を削除する。
func (r *repositoryImpl) DeleteIcon(username string) error {
	return nil
}
