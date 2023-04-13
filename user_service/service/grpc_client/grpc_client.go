package grpcclient

import (
	"fmt"

	"gitlab.com/micro/user_service/config"
	cc "gitlab.com/micro/user_service/genproto/comment"
	cu "gitlab.com/micro/user_service/genproto/post"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients interface {
	Post() cu.PostServiceClient
	Comment() cc.CommentServiceClient
}

type ServiceManager struct {
	Config         config.Config
	postService    cu.PostServiceClient
	commentService cc.CommentServiceClient
}

func New(cfg config.Config) (*ServiceManager, error) {
	connPost, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.PostServiceHost, cfg.PostServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("post service dial host:%s, port:%s", cfg.PostServiceHost, cfg.PostServicePort)
	}

	connComment, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.CommentServiceHost, cfg.CommentServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("comment service dial host:%s, port:%s", cfg.CommentServiceHost, cfg.CommentServicePort)
	}

	return &ServiceManager{
		Config:         cfg,
		postService:    cu.NewPostServiceClient(connPost),
		commentService: cc.NewCommentServiceClient(connComment),
	}, nil
}

func (s *ServiceManager) Post() cu.PostServiceClient {
	return s.postService
}

func (s *ServiceManager) Comment() cc.CommentServiceClient {
	return s.commentService
}
