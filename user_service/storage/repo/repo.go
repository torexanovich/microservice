package repo

import (
	u "gitlab.com/micro/user_service/genproto/user"
)
 
type UserStoreI interface {
	CreateUser(*u.UserRequest) (*u.UserResponse, error)
	GetUserById(*u.IdRequest) (*u.UserResponse, error)
	GetUserForClient(*u.IdRequest) (*u.UserResponse, error)
	GetAllUsers(*u.AllUsersRequest) (*u.Users, error)
	SearchUsersByName(*u.SearchUsers) (*u.Users, error)
	UpdateUser(*u.UpdateUserRequest) error
	DeleteUser(*u.IdRequest) (*u.UserResponse, error)
	CheckField(*u.CheckFieldReq) (*u.CheckFieldResp, error)
	GetByEmail(*u.EmailReq) (*u.UserResponse, error)
	CreateMod(*u.IdRequest) (*u.Empty, error)
}