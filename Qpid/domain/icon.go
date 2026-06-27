package domain

type IconMimeType string

const (
	IconMimeTypePNG  IconMimeType = "image/png"
	IconMimeTypeJPEG IconMimeType = "image/jpeg"
	IconMimeTypeWebp IconMimeType = "image/webp"
)

// ユーザーアイコン画像
type Icon struct {
	// 画像データ
	Blob []byte
	// 画像の MIME type
	MimeType IconMimeType
}
