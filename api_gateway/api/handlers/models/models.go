package models

type UserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type UpdateUserReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Id        string  `json:"id"`
}

type User struct {
	Id        string  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type IdRequest struct {
	Id int64 `json:"id"`
}

type GetAllUsersRequest struct {
	Limit int64 `json:"limit"`
	Page  int64 `json:"page"`
}

type Users struct {
	Users []User `json:"users"`
}

type SearchUsers struct {
	FirstName string `json:"first_name"`
}

type Empty struct {
}


// post

type PostRequest struct {
	UserId      int64  `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GetPostResponse struct {
	Id           int64      `json:"id"`
	UserId       string      `json:"user_id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Likes        int64      `json:"likes"`
	UserEmail    string     `json:"user_email"`
	UserName     string     `json:"user_name"`
	CreatedAt    string     `json:"created_at"`
	UpdatedAt    string     `json:"updated_at"`
	PostComments []Comments `json:"post_comments"`
}

type UpdatePostReq struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Posts struct {
	Posts []GetPostResponse `json:"posts"`
}

type Title struct {
	Title string `json:"title"`
}

type LikeRequest struct {
	PostId  int64 `json:"post_id"`
	IsLiked bool  `json:"is_liked"`
}

// Comment

type CommentRequest struct {
	UserId string  `json:"user_id"`
	PostId int64  `json:"post_id"`
	Text   string `json:"text"`
}

type CommentResponse struct {
	Id           int    `json:"id"`
	UserId       string  `json:"user_id"`
	PostId       int64  `json:"post_id"`
	Text         string `json:"text"`
	PostTitle    string `json:"post_title"`
	UserName     string `json:"user_name"`
	PostUserName string `json:"post_user_name"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type GetAllCommentsRequest struct {
	PostId int64 `json:"post_id"`
}

type Comments struct {
	Comments []CommentResponse `json:"comments"`
}

// register


type RegisterUserModel struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Code      string `json:"code"`
}

type VerifyResponse struct {
	Id           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	AccessToken  string `json:"accsee_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct{
	Id           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	AccessToken  string `json:"accsee_token"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type GetProfileByJwtRequest struct {
	Token string `header:"Authorization"`
}

// 

