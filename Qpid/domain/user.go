package domain

// ユーザープロフィール
type User struct {
	// ユーザー ID
	Username string
	// アイコン画像 ID; nil の場合は traQ のアイコンを参照する
	IconFileID string
	// 学部または系
	Major string
	// 所属班の一覧
	Affiliations []string
	// 出身地
	Hometown string
	// タグ名をキーにしたタグ情報
	Tags map[string]Tag
	// 自己紹介文
	Bio string
}

// プロフィールから抽出された興味や属性の情報
type Tag struct {
	// 好きなものや嫌いなものの分類
	Label string
	// タグに対する好き嫌い
	Affinity TagAffinity
	// 言及回数に対して単調増加する強度
	Strength float64
}

// タグに対する好き嫌い
type TagAffinity string

const (
	// 好意的なタグ
	TagAffinityPositive TagAffinity = "positive"
	// 中立的なタグ
	TagAffinityNeutral TagAffinity = "neutral"
	// 否定的なタグ
	TagAffinityNegative TagAffinity = "negative"
)
