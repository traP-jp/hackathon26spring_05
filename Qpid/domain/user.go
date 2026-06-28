package domain

import "github.com/moznion/go-optional"

// ユーザープロフィール
type User struct {
	// ユーザー ID
	Username string
	// 表示名
	DisplayName string
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

// UserOverride は事前計算済みプロフィールによるユーザー情報の上書き値。
// optional.Option は「未指定」と「空文字での指定」を区別するために使う。
type UserOverride struct {
	Major         optional.Option[string]
	Affiliations  []UserAffiliation
	Hometown      optional.Option[string]
	Tags          []string
	Technologies  []string
	Bio           optional.Option[string]
	FavoriteTopic optional.Option[TopicAndValue]
	DislikedTopic optional.Option[TopicAndValue]
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
