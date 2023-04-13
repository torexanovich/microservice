package repo

import (
	p "gitlab.com/micro/post_service/genproto/post"
)

type PostStorageI interface {
	CreatePost(*p.PostRequest) (*p.GetPostResponse, error)
	GetPostById(*p.IdRequest) (*p.GetPostResponse, error)
	GetPostByUserId(*p.IdUser) (*p.Posts, error)
	GetPostForUser(*p.IdUser) (*p.Posts, error)
	GetPostForComment(*p.IdRequest) (*p.GetPostResponse, error)
	SearchByTitle(*p.Title) (*p.Posts, error)
	LikePost(*p.LikeRequest) (*p.GetPostResponse, error)
	UpdatePost(*p.UpdatePostRequest) error
	DeletePost(*p.IdRequest) (*p.GetPostResponse, error)
}
