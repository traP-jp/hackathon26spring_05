package mock

import (
	"time"

	"github.com/traP-jp/hackathon26spring_05/Qpid/domain"
	"github.com/traP-jp/hackathon26spring_05/Qpid/repository"
)

var _ repository.Repository = (*MockRepository)(nil)

// MockRepository の mock 実装。
type MockRepository struct{}

// mock Repository を作成する。
func NewMockRepository() *MockRepository {
	return &MockRepository{}
}

// ユーザーを作成する。
func (r *MockRepository) Create(user domain.User) error {
	return nil
}

// ユーザーを取得する。
func (r *MockRepository) FindByUsername(username string) (*domain.User, error) {
	return mockUser(username), nil
}

// プロフィールを更新する。
func (r *MockRepository) UpdateProfile(username string, user domain.User) error {
	return nil
}

// ユーザーの存在を確認する。
func (r *MockRepository) Exists(username string) (bool, error) {
	return true, nil
}

// ユーザーを LIKE する。
func (r *MockRepository) Like(fromUsername, toUsername string) error {
	return nil
}

// ユーザーを NOPE する。
func (r *MockRepository) Nope(fromUsername, toUsername string) error {
	return nil
}

// LIKE したユーザーを取得する。
func (r *MockRepository) ListLikedUsers(username string) ([]domain.UserSummary, error) {
	return []domain.UserSummary{{Username: "liked-user"}}, nil
}

// LIKE してくれたユーザーを取得する。
func (r *MockRepository) ListUsersWhoLiked(username string) ([]domain.UserSummary, error) {
	return []domain.UserSummary{{Username: "liked-by-user"}}, nil
}

// アクション済みか確認する。
func (r *MockRepository) HasAction(fromUsername, toUsername string) (bool, error) {
	return false, nil
}

// おすすめユーザーを取得する。
func (r *MockRepository) ListSuggestions(username string, limit int) ([]domain.Suggestion, error) {
	return []domain.Suggestion{{Username: "suggested-user", Similarity: 0.5}}, nil
}

// 類似度を保存する。
func (r *MockRepository) UpsertSimilarity(usernameA, usernameB string, similarity float64) error {
	return nil
}

// 事前計算済みプロフィールを取得する。
func (r *MockRepository) FindPrecomputedProfileByUsername(username string) (*domain.User, error) {
	return mockUser(username), nil
}

// 事前計算済みユーザー名を一覧する。
func (r *MockRepository) ListPrecomputedUsernames() ([]string, error) {
	return []string{"precomputed-user"}, nil
}

// アイコン画像を保存する。
func (r *MockRepository) SaveIcon(username string, icon domain.Icon) error {
	return nil
}

// アイコン画像を取得する。
func (r *MockRepository) FindIconByUsername(username string) (*domain.Icon, error) {
	return &domain.Icon{
		Blob:     []byte("mock-icon"),
		MimeType: domain.IconMimeTypePNG,
	}, nil
}

// アイコン画像を削除する。
func (r *MockRepository) DeleteIcon(username string) error {
	return nil
}

// 認証セッションを作成する。
func (r *MockRepository) CreateAuthSession(sessionID string, session domain.AuthSession) error {
	return nil
}

// 認証セッションを取得する。
func (r *MockRepository) FindAuthSession(sessionID string) (*domain.AuthSession, error) {
	return &domain.AuthSession{
		Username:    "mock-user",
		AccessToken: "mock-access-token",
		ExpiresAt:   time.Now().Add(time.Hour),
	}, nil
}

// 認証セッションを削除する。
func (r *MockRepository) DeleteAuthSession(sessionID string) error {
	return nil
}

// OAuth2 state を作成する。
func (r *MockRepository) CreateOAuthState(state string, expiresAt time.Time) error {
	return nil
}

// OAuth2 state を消費する。
func (r *MockRepository) ConsumeOAuthState(state string) error {
	return nil
}

func mockUser(username string) *domain.User {
	return &domain.User{
		Username:     username,
		Affiliations: []string{"sysad"},
		Tags: map[string]domain.Tag{
			"go": {
				Label:    "programmingLanguage",
				Affinity: domain.TagAffinityPositive,
				Strength: 0.8,
			},
		},
		Bio: "mock user",
	}
}
