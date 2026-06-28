package infrastructure

import (
	"database/sql"
	"errors"

	"github.com/traP-jp/hackathon26spring_05/Qpid/domain"
)

type iconRow struct {
	Blob     []byte `db:"icon"`
	MimeType string `db:"mime_type"`
}

// アイコン画像を保存する。
func (r *repositoryImpl) SaveIcon(username string, icon domain.Icon) error {
	var count int

	//ユーザーが存在するかを確認
	err := r.db.Get(
		&count,
		`SELECT COUNT(*) FROM users WHERE username = ?`,
		username,
	)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("user not found")
	}

	_, err = r.db.Exec(
		`INSERT INTO icons (username, icon, mime_type)
		 VALUES (?, ?, ?)
		 ON DUPLICATE KEY UPDATE
			icon = VALUES(icon),
			mime_type = VALUES(mime_type)`,
		username,
		icon.Blob,
		icon.MimeType,
	)
	return err
}

// アイコン画像を取得する。
func (r *repositoryImpl) FindIconByUsername(username string) (*domain.Icon, error) {
	var icon iconRow
	err := r.db.Get(&icon,
		"SELECT icon, mime_type FROM icons WHERE username = ? LIMIT 1",
		username,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
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
	_, err := r.db.Exec(
		"DELETE FROM icons WHERE username = ?",
		username,
	)
	return err
}
