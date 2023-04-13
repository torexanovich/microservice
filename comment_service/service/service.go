package service

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"

	c "gitlab.com/micro/comment_service/genproto/comment"
	p "gitlab.com/micro/comment_service/genproto/post"
	u "gitlab.com/micro/comment_service/genproto/user"
	"gitlab.com/micro/comment_service/pkg/logger"
	grpcclient "gitlab.com/micro/comment_service/service/grpc_client"
	"gitlab.com/micro/comment_service/storage"
)

type CommentService struct {
	storage storage.IStorage
	Logger  logger.Logger
	Client  grpcclient.Clients
}

func NewCommentService(db *sqlx.DB, log logger.Logger, client grpcclient.Clients) *CommentService {
	return &CommentService{
		storage: storage.NewStoragePg(db),
		Logger:  log,
		Client:  client,
	}
}

func (s *CommentService) WriteComment(ctx context.Context, req *c.CommentRequest) (*c.CommentResponse, error) {
	res, err := s.storage.Comment().WriteComment(req)
	if err != nil {
		log.Println("failed to write comment: ", err)
		return &c.CommentResponse{}, err
	}

	post, err := s.Client.Post().GetPostForComment(ctx, &p.IdRequest{Id: res.PostId})
	if err != nil {
		log.Println("failed to getting post for write comment: ", err)
		return &c.CommentResponse{}, err
	}
	res.PostTitle = post.Title

	user, err := s.Client.User().GetUserForClient(ctx, &u.IdRequest{Id: res.UserId})
	if err != nil {
		log.Println("failed to get user for comment")
		return &c.CommentResponse{}, err
	}
	res.UserName = user.FirstName + " " + user.LastName

	postUser, err := s.Client.User().GetUserForClient(ctx, &u.IdRequest{Id: post.UserId})
	if err != nil {
		log.Println("failed to get post's user for write comment: ", err)
		return &c.CommentResponse{}, err
	}
	res.PostUserName = postUser.FirstName + " " + postUser.LastName

	return res, nil
}

func (s *CommentService) GetComments(ctx context.Context, req *c.GetAllCommentsRequest) (*c.Comments, error) {
	res, err := s.storage.Comment().GetComments(req)
	if err != nil {
		log.Println("failed to get comments: ", err)
		return &c.Comments{}, err
	}

	post, err := s.Client.Post().GetPostForComment(ctx, &p.IdRequest{Id: req.PostId})
	if err != nil {
		log.Println("failed to getting post for get comments: ", err)
		return &c.Comments{}, err
	}

	for _, comment := range res.Comments {
		user, err := s.Client.User().GetUserForClient(ctx, &u.IdRequest{Id: comment.UserId})
		if err != nil {
			log.Println("failed to get user for get comments: ", err)
			return &c.Comments{}, err
		}
		comment.UserName = user.FirstName + " " + user.LastName
	}

	postUser, err := s.Client.User().GetUserForClient(ctx, &u.IdRequest{Id: post.UserId})
	if err != nil {
		log.Println("failed to get post user for get comments: ", err)
		return &c.Comments{}, err
	}

	for _, comment := range res.Comments {
		comment.PostTitle = post.Title
		comment.PostUserName = postUser.FirstName + " " + postUser.LastName
	}

	return res, nil
}

func (s *CommentService) GetCommentsForPost(ctx context.Context, req *c.GetAllCommentsRequest) (*c.Comments, error) {
	res, err := s.storage.Comment().GetComments(req)
	if err != nil {
		log.Println("failed to get comments for client: ", err)
		return &c.Comments{}, err
	}

	return res, nil
}

func (s *CommentService) DeleteComment(ctx context.Context, id *c.IdRequest) (*c.CommentResponse, error) {
	res, err := s.storage.Comment().DeleteComment(id)
	if err != nil {
		log.Println("failed to delete comment: ", err)
		return &c.CommentResponse{}, err
	}

	post, err := s.Client.Post().GetPostForComment(ctx, &p.IdRequest{Id: res.PostId})
	if err != nil {
		log.Println("failed to getting post for delete comment: ", err)
		return &c.CommentResponse{}, err
	}
	res.PostTitle = post.Title

	user, err := s.Client.User().GetUserForClient(ctx, &u.IdRequest{Id: res.UserId})
	if err != nil {
		log.Println("failed to get user for delete comment: ", err)
		return &c.CommentResponse{}, err
	}
	res.UserName = user.FirstName + " " + user.LastName

	postUser, err := s.Client.User().GetUserForClient(ctx, &u.IdRequest{Id: post.UserId})
	if err != nil {
		log.Println("failed to get post user for delete comment: ", err)
		return &c.CommentResponse{}, err
	}
	res.PostUserName = postUser.FirstName + " " + postUser.LastName

	return res, nil
}
