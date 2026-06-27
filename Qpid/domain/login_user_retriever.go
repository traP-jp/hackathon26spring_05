package domain

type LoginUserRetriever interface {
	// ユーザーがログイン中かどうかを判定する
	IsUserLoggedIn() bool
	// ログイン中のユーザー名を取得する
	GetLoginUser() (string, error)
}
