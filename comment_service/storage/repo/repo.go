package repo

import (
	c "gitlab.com/micro/comment_service/genproto/comment"
)

type CommentStorageI interface {
	WriteComment(*c.CommentRequest) (*c.CommentResponse, error)
	GetComments(*c.GetAllCommentsRequest) (*c.Comments, error)
	GetCommentsForPost(*c.GetAllCommentsRequest) (*c.Comments, error)
	DeleteComment(*c.IdRequest) (*c.CommentResponse, error)
}
