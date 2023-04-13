package storage

import (
	"gitlab.com/micro/comment_service/storage/postgres"
	"gitlab.com/micro/comment_service/storage/repo"

	"github.com/jmoiron/sqlx"
)

type IStorage interface {
	Comment() repo.CommentStorageI
}

type storagePg struct {
	db          *sqlx.DB
	commentRepo repo.CommentStorageI
}

func NewStoragePg(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:          db,
		commentRepo: postgres.NewCommentRepo(db),
	}
}

func (s storagePg) Comment() repo.CommentStorageI {
	return s.commentRepo
}
