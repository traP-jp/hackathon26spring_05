package repository

import "github.com/traP-jp/hackathon26spring_05/Qpid/domain"

// DB に保存するユーザーアイコン画像を扱う
type IconRepository interface {
	// アイコン画像を保存する
	SaveIcon(username string, icon domain.Icon) error
	// アイコン画像を取得する
	FindIconByUsername(username string) (*domain.Icon, error)
	// アイコン画像を削除する
	DeleteIcon(username string) error
}
