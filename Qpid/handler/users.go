package handler

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/moznion/go-optional"
	"github.com/traP-jp/hackathon26spring_05/Qpid/domain"
)

// domain.User の定義に沿った、最終的なJSONレスポンスの型
type userResponse struct {
	Username      string                                 `json:"username"`
	HasIcon       bool                                   `json:"hasIcon"`
	Major         optional.Option[string]                `json:"major"`
	Affiliations  []domain.UserAffiliation               `json:"affiliations"`
	Hometown      optional.Option[string]                `json:"hometown"`
	Tags          []string                               `json:"tags"`
	Technologies  []string                               `json:"technologies"`
	Bio           optional.Option[string]                `json:"bio"`
	FavoriteTopic optional.Option[topicAndValueResponse] `json:"favoriteTopic"`
	DislikedTopic optional.Option[topicAndValueResponse] `json:"dislikedTopic"`
}

type topicAndValueResponse struct {
	Topic string `json:"topic"`
	Value string `json:"value"`
}

// GET /api/users/:id
func (h *handler) getUser(c *echo.Context) error {
	// ② パスパラメータから「id」を取得
	userID, err := echo.PathParam[string](c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "invalid path parameter"})
	}

	// ③ データベース等からユーザーを取得
	user, err := h.repository.FindUserByUsername(userID)
	if err != nil {
		c.Logger().Error(
			"failed to find user",
			slog.String("username", userID),
			slog.Any("error", err),
		)
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
func toUserResponse(user domain.User) userResponse {
	return userResponse{
		Username:      user.Username,
		HasIcon:       user.HasIcon,
		Major:         user.Major,
		Affiliations:  user.Affiliations,
		Hometown:      user.Hometown,
		Tags:          user.Tags,
		Technologies:  user.Technologies,
		Bio:           user.Bio,
		FavoriteTopic: toTopicAndValueResponse(user.FavoriteTopic),
		DislikedTopic: toTopicAndValueResponse(user.DislikedTopic),
	}
}

func toTopicAndValueResponse(topicAndValue optional.Option[domain.TopicAndValue]) optional.Option[topicAndValueResponse] {
	if topicAndValue.IsNone() {
		return optional.None[topicAndValueResponse]()
	}
	topicAndValueValue := topicAndValue.Unwrap()
	return optional.Some(topicAndValueResponse{
		Topic: topicAndValueValue.Topic,
		Value: topicAndValueValue.Value,
	})
}
