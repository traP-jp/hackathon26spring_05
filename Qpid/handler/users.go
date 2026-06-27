package handler

import "github.com/labstack/echo/v4"

// タグに対する好き嫌いの型（domainの定義をそのまま使えるなら domain.TagAffinity に置き換えてください）
type TagResponse struct {
	Label    string  `json:"label"`
	Affinity string  `json:"affinity"` // "positive", "neutral", "negative"
	Strength float64 `json:"strength"`
}

// domain.User の定義に沿った、最終的なJSONレスポンスの型
type UserResponse struct {
	Username     string                 `json:"username"`
	IconFileID   *string                `json:"iconFileId"` // nil の可能性あり
	Major        *string                `json:"major"`      // nil の可能性あり
	Affiliations []string               `json:"affiliations"`
	Hometown     *string                `json:"hometown"`   // nil の可能性あり
	Tags         map[string]TagResponse `json:"tags"`
	Bio          string                 `json:"bio"`
}

// GET /api/users/:id
func (h *handler) getUser(c echo.Context) error {
	// ① パスパラメータから「id」を取得
	userID := c.Param("id")

	// ポインタ型に値を代入するための補助変数
	major := "情報理工学系"
	hometown := "東京"

	// ② ドメインモデルに沿った詳細なモックデータ（仮データ）を作成
	mockUser := UserResponse{
		Username:     userID,
		IconFileID:   nil, // 初期状態としてnilをテスト
		Major:        &major,
		Affiliations: []string{"開発班", "グラフィック班"},
		Hometown:     &hometown,
		Tags: map[string]TagResponse{
			"Go言語": {
				Label:    "プログラミング言語",
				Affinity: "positive",
				Strength: 4.5,
			},
			"早起き": {
				Label:    "生活習慣",
				Affinity: "negative",
				Strength: 1.2,
			},
		},
		Bio: "こんにちは！ハッカソン頑張りましょう！",
	}

	// ③ ステータス200（成功）と、詳細なJSONデータを返却
	return c.JSON(200, mockUser)
}
