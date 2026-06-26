package repository

import "github.com/jmoiron/sqlx"

type Repository any

type repositoryImpl struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repositoryImpl{db: db}
}
