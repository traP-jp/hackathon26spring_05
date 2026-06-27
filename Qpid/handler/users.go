package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traP-jp/hackathon26spring_05/Qpid/domain"
	"github.com/traP-jp/hackathon26spring_05/Qpid/handler/middleware"
)

// タグに対する好き嫌いの型
type TagResponse struct {
	Label    string             `json:"label"`
	Affinity domain.TagAffinity `json:"affinity"` // "positive", "neutral", "negative"
	Strength float64            `json:"strength"`
}

// domain.User の定義に沿った、最終的なJSONレスポンスの型
type UserResponse struct {
	Username     string                 `json:"username"`
	IconFileID   string                 `json:"iconFileId"`
	Major        string                 `json:"major"`
	Affiliations []string               `json:"affiliations"`
	Hometown     string                 `json:"hometown"`
	Tags         map[string]TagResponse `json:"tags"`
	Bio          string                 `json:"bio"`
}

// GET /api/users/:id
func (h *handler) getUser(c echo.Context) error {
	// ① ログインチェックを追加
	loginUserRetriever := middleware.GetLoginUserRetriever(c)

	if !loginUserRetriever.IsUserLoggedIn() {
		return unauthorized(c)
	}

	// ② パスパラメータから「id」を取得
	userID := c.Param("id")

	// ③ データベース等からユーザーを取得
	user, err := h.repository.FindByUsername(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to load user"})
	}
	if user == nil {
		//調べようとしてるユーザーが見つからないケース
		return notFound(c)
	}

	// ④ ドメインモデルを UserResponse に詰め替えてステータス200で返却
	return c.JSON(http.StatusOK, toUserResponse(*user))
}

// domain.User から UserResponse へ変換するヘルパー関数
func toUserResponse(user domain.User) UserResponse {
	tags := make(map[string]TagResponse, len(user.Tags))
	for name, userTag := range user.Tags {
		tags[name] = TagResponse{
			Label:    userTag.Label,
			Affinity: userTag.Affinity, // domain.TagAffinity 型のまま代入
			Strength: userTag.Strength,
		}
	}

	return UserResponse{
		Username:     user.Username,
		IconFileID:   user.IconFileID,
		Major:        user.Major,
		Affiliations: user.Affiliations,
		Hometown:     user.Hometown,
		Tags:         tags,
		Bio:          user.Bio,
	}
}
