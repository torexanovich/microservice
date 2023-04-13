package service

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"

	c "gitlab.com/micro/post_service/genproto/comment"
	p "gitlab.com/micro/post_service/genproto/post"
	u "gitlab.com/micro/post_service/genproto/user"
	"gitlab.com/micro/post_service/pkg/logger"
	grpcclient "gitlab.com/micro/post_service/service/grpc_client"
	"gitlab.com/micro/post_service/storage"

	"google.golang.org/protobuf/types/known/emptypb"
)

type PostService struct {
	storage storage.IStorage
	Logger  logger.Logger
	Client  grpcclient.Clients
}

func NewPostService(db *sqlx.DB, log logger.Logger, client grpcclient.Clients) *PostService {
	return &PostService{
		storage: storage.NewStoragePg(db),
		Logger:  log,
		Client:  client,
	}
}

func (s *PostService) CreatePost(ctx context.Context, req *p.PostRequest) (*p.GetPostResponse, error) {
	res, err := s.storage.Post().CreatePost(req)
	if err != nil {
		log.Println("failed to create post: ", err)
		return &p.GetPostResponse{}, err
	}

	user, err := s.Client.User().GetUserForClient(ctx, &u.IdRequest{Id: res.UserId})
	if err != nil {
		log.Println("failed to getting user for create post: ", err)
		return &p.GetPostResponse{}, err
	}
	res.UserName = user.FirstName + " " + user.LastName
	res.UserEmail = user.Email

	return res, nil
}

func (s *PostService) GetPostById(ctx context.Context, req *p.IdRequest) (*p.GetPostResponse, error) {
	res, err := s.storage.Post().GetPostById(req)
	if err != nil {
		log.Println("failed to get post by id: ", err)
		return &p.GetPostResponse{}, err
	}

	postUser, err := s.Client.User().GetUserForClient(ctx, &u.IdRequest{Id: res.UserId})
	if err != nil {
		log.Println("failed to getting user for get post by id: ", err)
		return &p.GetPostResponse{}, err
	}
	res.UserName = postUser.FirstName + " " + postUser.LastName
	res.UserEmail = postUser.Email

	comments, err := s.Client.Comment().GetCommentsForPost(ctx, &c.GetAllCommentsRequest{PostId: res.Id})
	if err != nil {
		log.Println("failed to getting comments for get post by id: ", err)
		return &p.GetPostResponse{}, err
	}

	for _, comment := range comments.Comments {
		comUser, err := s.Client.User().GetUserForClient(ctx, &u.IdRequest{Id: postUser.Id})
		if err != nil {
			log.Println("failed to get user for comment: ", err)
			return &p.GetPostResponse{}, err
		}

		com := p.Comment{}
		com.UserId = comUser.Id
		com.UserName = comUser.FirstName + " " + comUser.LastName
		com.PostId = comment.PostId
		com.PostTitle = res.Title
		com.PostUserName = postUser.FirstName + " " + postUser.LastName
		com.Text = comment.Text
		com.CreatedAt = comment.CreatedAt

		res.Comments = append(res.Comments, &com)
	}

	return res, nil
}

func (s *PostService) GetPostByUserId(ctx context.Context, req *p.IdUser) (*p.Posts, error) {
	res, err := s.storage.Post().GetPostByUserId(req)
	if err != nil {
		log.Println("failed to get post by user id: ", err)
		return &p.Posts{}, err
	}

	postUser, err := s.Client.User().GetUserForClient(ctx, &u.IdRequest{Id: req.Id})
	if err != nil {
		log.Println("failed to getting user for get post by user id: ", err)
		return &p.Posts{}, err
	}

	for _, p := range res.Posts {
		p.UserName = postUser.FirstName + " " + postUser.LastName
		p.UserEmail = postUser.Email
	}

	for _, post := range res.Posts {
		comments, err := s.Client.Comment().GetCommentsForPost(ctx, &c.GetAllCommentsRequest{PostId: post.Id})
		if err != nil {
			log.Println("failed to getting comments for get post by user id: ", err)
			return &p.Posts{}, err
		}

		for _, comment := range comments.Comments {
			comUser, err := s.Client.User().GetUserForClient(ctx, &u.IdRequest{Id: comment.UserId})
			if err != nil {
				log.Println("failed to get user for comment: ", err)
				return &p.Posts{}, err
			}

			com := p.Comment{}
			com.UserId = comUser.Id
			com.UserName = comUser.FirstName + " " + comUser.LastName
			com.PostId = comment.PostId
			com.PostTitle = post.Title
			com.PostUserName = postUser.FirstName + " " + postUser.LastName
			com.Text = comment.Text
			com.CreatedAt = comment.CreatedAt

			post.Comments = append(post.Comments, &com)
		}
	}

	return res, nil
}

func (s *PostService) GetPostForUser(ctx context.Context, req *p.IdUser) (*p.Posts, error) {
	res, err := s.storage.Post().GetPostForUser(req)
	if err != nil {
		log.Println("failed to get post for user: ", err)
		return &p.Posts{}, err
	}

	return res, nil
}

func (s *PostService) GetPostForComment(ctx context.Context, req *p.IdRequest) (*p.GetPostResponse, error) {
	res, err := s.storage.Post().GetPostForComment(req)
	if err != nil {
		log.Println("failed to get post for comment: ", err)
		return &p.GetPostResponse{}, err
	}

	return res, nil
}

func (s *PostService) SearchByTitle(ctx context.Context, req *p.Title) (*p.Posts, error) {
	res, err := s.storage.Post().SearchByTitle(req)
	if err != nil {
		log.Println("failed to get post by search title: ", err)
		return &p.Posts{}, err
	}

	for _, post := range res.Posts {
		postUser, err := s.Client.User().GetUserForClient(ctx, &u.IdRequest{Id: post.UserId})
		if err != nil {
			log.Println("failed to getting user for get post by search title: ", err)
			return &p.Posts{}, err
		}

		post.UserName = postUser.FirstName + " " + postUser.LastName
		post.UserEmail = postUser.Email

		comments, err := s.Client.Comment().GetCommentsForPost(ctx, &c.GetAllCommentsRequest{PostId: post.Id})
		if err != nil {
			log.Println("failed to getting comments for get post by search title: ", err)
			return &p.Posts{}, err
		}

		for _, comment := range comments.Comments {
			comUser, err := s.Client.User().GetUserForClient(ctx, &u.IdRequest{Id: comment.UserId})
			if err != nil {
				log.Println("failed to get user for comment: ", err)
				return &p.Posts{}, err
			}

			com := p.Comment{}
			com.UserId = comUser.Id
			com.UserName = comUser.FirstName + " " + comUser.LastName
			com.PostId = comment.PostId
			com.PostTitle = post.Title
			com.PostUserName = postUser.FirstName + " " + postUser.LastName
			com.Text = comment.Text
			com.CreatedAt = comment.CreatedAt

			post.Comments = append(post.Comments, &com)
		}
	}

	return res, nil
}

func (s *PostService) LikePost(ctx context.Context, req *p.LikeRequest) (*p.GetPostResponse, error) {
	res, err := s.storage.Post().LikePost(req)
	if err != nil {
		log.Println("failed to like post: ", err)
		return &p.GetPostResponse{}, err
	}
	if !req.IsLiked {
		res.Likes -= 1
	}

	postUser, err := s.Client.User().GetUserForClient(ctx, &u.IdRequest{Id: res.UserId})
	if err != nil {
		log.Println("failed to getting user for like post: ", err)
		return &p.GetPostResponse{}, err
	}
	res.UserName = postUser.FirstName + " " + postUser.LastName
	res.UserEmail = postUser.Email

	comments, err := s.Client.Comment().GetCommentsForPost(ctx, &c.GetAllCommentsRequest{PostId: res.Id})
	if err != nil {
		log.Println("failed to getting comments for like post: ", err)
		return &p.GetPostResponse{}, err
	}

	for _, comment := range comments.Comments {
		comUser, err := s.Client.User().GetUserForClient(ctx, &u.IdRequest{Id: comment.UserId})
		if err != nil {
			log.Println("failed to get user for comment: ", err)
			return &p.GetPostResponse{}, err
		}

		com := p.Comment{}
		com.UserId = comUser.Id
		com.UserName = comUser.FirstName + " " + comUser.LastName
		com.PostId = comment.PostId
		com.PostTitle = res.Title
		com.PostUserName = postUser.FirstName + " " + postUser.LastName
		com.Text = comment.Text
		com.CreatedAt = comment.CreatedAt

		res.Comments = append(res.Comments, &com)
	}

	return res, nil
}

func (s *PostService) UpdatePost(ctx context.Context, req *p.UpdatePostRequest) (*emptypb.Empty, error) {
	err := s.storage.Post().UpdatePost(req)
	if err != nil {
		log.Println("failed to update post: ", err)
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (s *PostService) DeletePost(ctx context.Context, req *p.IdRequest) (*p.GetPostResponse, error) {
	res, err := s.storage.Post().DeletePost(req)
	if err != nil {
		log.Println("failed to delete post: ", err)
		return &p.GetPostResponse{}, err
	}

	postUser, err := s.Client.User().GetUserForClient(ctx, &u.IdRequest{Id: res.UserId})
	if err != nil {
		log.Println("failed to getting user for delete post: ", err)
		return &p.GetPostResponse{}, err
	}
	res.UserName = postUser.FirstName + " " + postUser.LastName
	res.UserEmail = postUser.Email

	comments, err := s.Client.Comment().GetCommentsForPost(ctx, &c.GetAllCommentsRequest{PostId: res.Id})
	if err != nil {
		log.Println("failed to getting comments for delete post: ", err)
		return &p.GetPostResponse{}, err
	}

	for _, comment := range comments.Comments {
		comUser, err := s.Client.User().GetUserForClient(ctx, &u.IdRequest{Id: comment.UserId})
		if err != nil {
			log.Println("failed to get user for comment: ", err)
			return &p.GetPostResponse{}, err
		}

		com := p.Comment{}
		com.UserId = comUser.Id
		com.UserName = comUser.FirstName + " " + comUser.LastName
		com.PostId = comment.PostId
		com.PostTitle = res.Title
		com.PostUserName = postUser.FirstName + " " + postUser.LastName
		com.Text = comment.Text
		com.CreatedAt = comment.CreatedAt

		res.Comments = append(res.Comments, &com)
	}

	return res, err
}
