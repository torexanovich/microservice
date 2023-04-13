package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	c "gitlab.com/micro/user_service/genproto/comment"
	p "gitlab.com/micro/user_service/genproto/post"
	u "gitlab.com/micro/user_service/genproto/user"
	"gitlab.com/micro/user_service/pkg/logger"
	grpcclient "gitlab.com/micro/user_service/service/grpc_client"
	"gitlab.com/micro/user_service/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService struct {
	storage storage.IStorage
	Logger  logger.Logger
	Client  grpcclient.Clients
}

func NewUserService(db *sqlx.DB, log logger.Logger, client grpcclient.Clients) *UserService {
	return &UserService{
		storage: storage.NewStoragePg(db),
		Logger:  log,
		Client:  client,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *u.UserRequest) (*u.UserResponse, error) {
	res, err := s.storage.User().CreateUser(req)
	fmt.Println(res)
	if err != nil {
		log.Println("failed to creating user: ", err)
		return &u.UserResponse{}, err
	}

	return res, nil
}

func (s *UserService) GetUserById(ctx context.Context, req *u.IdRequest) (*u.UserResponse, error) {
	res, err := s.storage.User().GetUserById(req)
	if err != nil {
		log.Println("failed to getting user: ", err)
		return &u.UserResponse{}, err
	}

	post, err := s.Client.Post().GetPostForUser(ctx, &p.IdUser{Id: req.Id})
	if err != nil {
		log.Println("failed to getting post for get user: ", err)
		return &u.UserResponse{}, err
	}

	for _, pt := range post.Posts {
		comments, err := s.Client.Comment().GetCommentsForPost(ctx, &c.GetAllCommentsRequest{PostId: pt.Id})
		if err != nil {
			log.Println("failed to get comments for post in user service: ", err)
			return &u.UserResponse{}, err
		}

		pst := u.Post{}

		for _, comment := range comments.Comments {
			comUser, err := s.storage.User().GetUserById(&u.IdRequest{Id: comment.UserId})
			if err != nil {
				log.Println("failed to get comment user: ", err)
				return &u.UserResponse{}, err
			}

			com := u.Comment{}
			com.UserId = comUser.Id
			com.UserName = comUser.FirstName + " " + comUser.LastName
			com.PostId = comment.PostId
			com.PostTitle = pt.Title
			com.PostUserName = res.FirstName + " " + res.LastName
			com.Text = comment.Text
			com.CreatedAt = comment.CreatedAt

			pst.Comments = append(pst.Comments, &com)
		}
		pst.Id = pt.Id
		pst.Title = pt.Title
		pst.Description = pt.Description
		pst.Likes = pt.Likes
		pst.CreatedAt = pt.CreatedAt
		pst.UpdatedAt = pt.UpdatedAt

		res.Posts = append(res.Posts, &pst)
	}

	return res, nil
}

func (s *UserService) GetUserForClient(ctx context.Context, req *u.IdRequest) (*u.UserResponse, error) {
	res, err := s.storage.User().GetUserById(req)
	if err != nil {
		log.Println("failed to getting user for clients: ", err)
		return &u.UserResponse{}, err
	}

	return res, nil
}

func (s *UserService) GetAllUsers(ctx context.Context, req *u.AllUsersRequest) (*u.Users, error) {
	res, err := s.storage.User().GetAllUsers(req)
	if err != nil {
		log.Println("failed to getting all users: ", err)
		return &u.Users{}, err
	}

	for _, user := range res.Users {
		post, err := s.Client.Post().GetPostForUser(ctx, &p.IdUser{Id: user.Id})
		if err != nil {
			log.Println("failed to getting post for get all users: ", err)
			return &u.Users{}, err
		}

		for _, pt := range post.Posts {
			comments, err := s.Client.Comment().GetCommentsForPost(ctx, &c.GetAllCommentsRequest{PostId: pt.Id})
			if err != nil {
				log.Println("failed to get comments for post in user service: ", err)
				return &u.Users{}, err
			}

			pst := u.Post{}

			for _, comment := range comments.Comments {
				comUser, err := s.storage.User().GetUserById(&u.IdRequest{Id: comment.UserId})
				if err != nil {
					log.Println("failed to get user comment: ", err)
					return &u.Users{}, err
				}

				com := u.Comment{}
				com.UserId = comment.UserId
				com.UserName = comUser.FirstName + " " + comUser.LastName
				com.PostId = comment.PostId
				com.PostTitle = pt.Title
				com.PostUserName = user.FirstName + " " + user.LastName
				com.Text = comment.Text
				com.CreatedAt = comment.CreatedAt

				pst.Comments = append(pst.Comments, &com)
			}
			pst.Id = pt.Id
			pst.Title = pt.Title
			pst.Description = pt.Description
			pst.Likes = pt.Likes
			pst.CreatedAt = pt.CreatedAt
			pst.UpdatedAt = pt.UpdatedAt

			user.Posts = append(user.Posts, &pst)
		}
	}

	return res, nil
}

func (s *UserService) SearchUsersByName(ctx context.Context, req *u.SearchUsers) (*u.Users, error) {
	res, err := s.storage.User().SearchUsersByName(req)
	if err != nil {
		log.Println("failed to searching user by name: ", err)
		return &u.Users{}, err
	}

	for _, user := range res.Users {
		post, err := s.Client.Post().GetPostForUser(ctx, &p.IdUser{Id: user.Id})
		if err != nil {
			log.Println("failed to getting post for searching user by name: ", err)
			return &u.Users{}, err
		}

		for _, pt := range post.Posts {
			comments, err := s.Client.Comment().GetCommentsForPost(ctx, &c.GetAllCommentsRequest{PostId: pt.Id})
			if err != nil {
				log.Println("failed to get comments for post in user service: ", err)
				return &u.Users{}, err
			}

			pst := u.Post{}

			for _, comment := range comments.Comments {
				comUser, err := s.storage.User().GetUserById(&u.IdRequest{Id: comment.UserId})
				if err != nil {
					log.Println("failed to get user comment: ", err)
					return &u.Users{}, err
				}

				com := u.Comment{}
				com.UserId = comment.UserId
				com.UserName = comUser.FirstName + " " + comUser.LastName
				com.PostId = comment.PostId
				com.PostTitle = pt.Title
				com.PostUserName = user.FirstName + " " + user.LastName
				com.Text = comment.Text
				com.CreatedAt = comment.CreatedAt

				pst.Comments = append(pst.Comments, &com)
			}
			pst.Id = pt.Id
			pst.Title = pt.Title
			pst.Description = pt.Description
			pst.Likes = pt.Likes
			pst.CreatedAt = pt.CreatedAt
			pst.UpdatedAt = pt.UpdatedAt

			user.Posts = append(user.Posts, &pst)
		}
	}

	return res, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *u.UpdateUserRequest) (*emptypb.Empty, error) {
	err := s.storage.User().UpdateUser(req)
	if err != nil {
		log.Println("failed to updating user: ", err)
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *u.IdRequest) (*u.UserResponse, error) {
	res, err := s.storage.User().DeleteUser(req)
	if err != nil {
		log.Println("failed to deleting user: ", err)
		return &u.UserResponse{}, err
	}

	post, err := s.Client.Post().GetPostForUser(ctx, &p.IdUser{Id: req.Id})
	if err != nil {
		log.Println("failed to getting post for deleting user: ", err)
		return &u.UserResponse{}, err
	}

	for _, pt := range post.Posts {
		comments, err := s.Client.Comment().GetCommentsForPost(ctx, &c.GetAllCommentsRequest{PostId: pt.Id})
		if err != nil {
			log.Println("failed to get comments for post in user service: ", err)
			return &u.UserResponse{}, err
		}

		pst := u.Post{}

		for _, comment := range comments.Comments {
			comUser, err := s.storage.User().GetUserById(&u.IdRequest{Id: comment.UserId})
			if err != nil {
				log.Println("failed to get user comment: ", err)
				return &u.UserResponse{}, err
			}

			com := u.Comment{}
			com.UserId = comment.UserId
			com.UserName = comUser.FirstName + " " + comUser.LastName
			com.PostId = comment.PostId
			com.PostTitle = pt.Title
			com.UserName = res.FirstName + " " + res.LastName
			com.Text = comment.Text
			com.CreatedAt = comment.CreatedAt

			pst.Comments = append(pst.Comments, &com)
		}
		pst.Id = pt.Id
		pst.Title = pt.Title
		pst.Description = pt.Description
		pst.Likes = pt.Likes
		pst.CreatedAt = pt.CreatedAt
		pst.UpdatedAt = pt.UpdatedAt

		res.Posts = append(res.Posts, &pst)
	}

	return res, nil
}

func (s *UserService) CheckField(ctx context.Context, req *u.CheckFieldReq) (*u.CheckFieldResp, error) {
	boolean, err := s.storage.User().CheckField(req)
	if err != nil {
		s.Logger.Error("error checking field", logger.Any("error while checking field", err))
		return &u.CheckFieldResp{}, status.Error(codes.Internal, "Internal server error")
	}

	return &u.CheckFieldResp{Exists: boolean.Exists}, nil
}

func (s *UserService) Login(ctx context.Context, req *u.LoginReq) (*u.LoginResp, error) {
	req.Email = strings.TrimSpace(req.Email)
	req.Email = strings.ToLower(req.Email)

	user, err := s.storage.User().GetByEmail(&u.EmailReq{Email: req.Email})
	if err == sql.ErrNoRows {
		s.Logger.Error("error while getting user by email, Not Found", logger.Any("req", req))
		return nil, status.Error(codes.NotFound, "Not Found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		s.Logger.Error("error while comparing hashed password, Invalid credentials", logger.Any("req", req))
		return nil, status.Error(codes.InvalidArgument, "Invalid credentials")
	}

	return &u.LoginResp{
		Id:        user.Id,
		Email:     req.Email,
		Password:  req.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Posts:     user.Posts,
	}, nil
}

func (s *UserService) GetByEmail(ctx context.Context, req *u.EmailReq) (*u.UserResponse, error) {
	res, err := s.storage.User().GetByEmail(req)
	if err != nil {
		log.Println("failed to getting user by email: ", err)
		return &u.UserResponse{}, err
	}

	return res, nil
}
