package infrastructure

import (
	"github.com/jmoiron/sqlx"
)

// var _ repository.Repository = (*repositoryImpl)(nil)

type repositoryImpl struct {
	db *sqlx.DB
}

// DB を使う Repository 実装を作成する。
func NewRepository(db *sqlx.DB) *repositoryImpl {
	return &repositoryImpl{db: db}
}
