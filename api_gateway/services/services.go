package services

import (
	"fmt"

	"gitlab.com/micro/api_gateway/config"
	pc "gitlab.com/micro/api_gateway/genproto/comment"
	pp "gitlab.com/micro/api_gateway/genproto/post"
	pu "gitlab.com/micro/api_gateway/genproto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

type IServiceManager interface {
	UserService() pu.UserServiceClient
	PostService() pp.PostServiceClient
	CommentService() pc.CommentServiceClient
}

type serviceManager struct {
	userService    pu.UserServiceClient
	postService    pp.PostServiceClient
	commentService pc.CommentServiceClient
}

func NewServiceManager(conf *config.Config) (IServiceManager, error) {
	resolver.SetDefaultScheme("dns")

	connUser, err := grpc.Dial(
		fmt.Sprintf("%s:%s", conf.UserServiceHost, conf.UserServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	connPost, err := grpc.Dial(
		fmt.Sprintf("%s:%s", conf.PostServiceHost, conf.PostServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	connComment, err := grpc.Dial(
		fmt.Sprintf("%s:%s", conf.CommentServiceHost, conf.CommentServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	serviceManager := &serviceManager{
		userService:    pu.NewUserServiceClient(connUser),
		postService:    pp.NewPostServiceClient(connPost),
		commentService: pc.NewCommentServiceClient(connComment),
	}

	return serviceManager, nil
}

func (s *serviceManager) UserService() pu.UserServiceClient {
	return s.userService
}

func (s *serviceManager) PostService() pp.PostServiceClient {
	return s.postService
}

func (s *serviceManager) CommentService() pc.CommentServiceClient {
	return s.commentService
}
