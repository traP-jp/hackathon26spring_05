package repository

import "github.com/traP-jp/hackathon26spring_05/Qpid/domain"

// 事前計算済みプロフィールを読み込む
type PrecomputedProfileRepository interface {
	// 事前計算済みプロフィールを取得する
	FindPrecomputedProfileByUsername(username string) (*domain.User, error)
	// 事前計算済みユーザー名を一覧する
	ListPrecomputedUsernames() ([]string, error)
}
