package grpcclient

import (
	"fmt"

	"gitlab.com/micro/post_service/config"
	cc "gitlab.com/micro/post_service/genproto/comment"
	cu "gitlab.com/micro/post_service/genproto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients interface {
	User() cu.UserServiceClient
	Comment() cc.CommentServiceClient
}

type ServiceManager struct {
	Config         config.Config
	userService    cu.UserServiceClient
	commentService cc.CommentServiceClient
}

func New(cfg config.Config) (*ServiceManager, error) {
	connUser, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.UserServiceHost, cfg.UserServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("user service dial host:%s, port:%s", cfg.UserServiceHost, cfg.UserServicePort)
	}

	connComment, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.CommentServiceHost, cfg.CommentServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("comment service dial host:%s, port:%s", cfg.CommentServiceHost, cfg.CommentServicePort)
	}

	return &ServiceManager{
		Config:         cfg,
		userService:    cu.NewUserServiceClient(connUser),
		commentService: cc.NewCommentServiceClient(connComment),
	}, nil
}

func (s *ServiceManager) User() cu.UserServiceClient {
	return s.userService
}

func (s *ServiceManager) Comment() cc.CommentServiceClient {
	return s.commentService
}
