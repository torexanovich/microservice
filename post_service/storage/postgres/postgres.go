package postgres

import "github.com/jmoiron/sqlx"

type PostRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) *PostRepo {
	return &PostRepo{
		db: db,
	}
}