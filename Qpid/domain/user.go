package domain

import "github.com/moznion/go-optional"

// ユーザープロフィール
type User struct {
	// ユーザー ID
	Username string
	// Qpidにアップロードされたアイコンがあるかどうか
	HasIcon bool
	// 学部または系
	Major optional.Option[string]
	// 所属班の一覧
	Affiliations []UserAffiliation
	// 出身地
	Hometown optional.Option[string]
	// タグ
	Tags []string
	// よく使う技術
	Technologies []string
	// 自己紹介文
	Bio optional.Option[string]
	// 好きな〇〇
	FavoriteTopic optional.Option[TopicAndValue]
	// 嫌いな〇〇
	DislikedTopic optional.Option[TopicAndValue]
}

type TopicAndValue struct {
	Topic string
	Value string
}

type UserAffiliation string

const (
	UserAffiliationAlgorithm UserAffiliation = "algorithm"
	UserAffiliationCtf       UserAffiliation = "ctf"
	UserAffiliationGame      UserAffiliation = "game"
	UserAffiliationGraphics  UserAffiliation = "graphics"
	UserAffiliationKaggle    UserAffiliation = "kaggle"
	UserAffiliationSound     UserAffiliation = "sound"
	UserAffiliationSysAd     UserAffiliation = "sysad"
)
