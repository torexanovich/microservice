package postgres

// import (
// 	"testing"

// 	"github.com/stretchr/testify/suite"
// 	"gitlab.com/micro/user_service/config"
// 	pg "gitlab.com/micro/user_service/genproto/user"
// 	"gitlab.com/micro/user_service/pkg/db"
// 	"gitlab.com/micro/user_service/storage/repo"
// )

// type UserSuiteTest struct {
// 	suite.Suite
// 	cleanUpFunc func ()
// 	repo repo.UserStoreI
// }

// func (s *UserSuiteTest)	SetupSuite(){
// 	pgPool, cleanUp := db.ConnectToDBForSuite(config.Load())
// 	s.repo = NewUserRepo(pgPool)
// 	s.cleanUpFunc = cleanUp
// }

// // TEST CRUD
// func (s *UserSuiteTest) TestUserCrud() {
// 	user := &pg.UserRequest {
// 		FirstName: "John",
// 		LastName: "Doe",
// 	}

// 	// CREATE
// 	createUserResp, err := s.repo.CreateUser(user)
// 	s.Nil(err)
// 	s.NotNil(createUserResp)
// 	s.Equal(user.FirstName, createUserResp.FirstName)
// 	s.Equal(user.LastName, createUserResp.LastName)

// 	// GET
// 	getUserResp, err := s.repo.GetUserById(&pg.IdRequest{Id: createUserResp.Id})
// 	s.Nil(err)
// 	s.NotNil(getUserResp)
// 	s.Equal(getUserResp.FirstName, user.FirstName)
// 	s.Equal(getUserResp.LastName, user.LastName)

// 	// UPDATE
// 	updateBody := &pg.UpdateUserRequest{
// 		Id: createUserResp.Id,
// 		FirstName: "John",
// 		LastName: "Smith",
// 	}
// 	err = s.repo.UpdateUser(updateBody)
// 	s.Nil(err)
// 	updateUserResp, err := s.repo.GetUserById(&pg.IdRequest{Id: createUserResp.Id})
// 	s.Nil(err)
// 	s.NotNil(updateUserResp)
// 	s.NotEqual(user.LastName, updateUserResp.LastName)
// 	s.Equal(user.FirstName, updateUserResp.FirstName)

// 	// GET ALL
// 	listUsersResp, err := s.repo.GetAllUsers(&pg.AllUsersRequest{Page: 1, Limit: 100})
// 	s.Nil(err)
// 	s.NotNil(listUsersResp)

// 	// GET USERS FOR CLIENT
// 	getForClientResp, err := s.repo.GetUserForClient(&pg.IdRequest{Id: createUserResp.Id})
// 	s.Nil(err)
// 	s.NotNil(getForClientResp)

// 	// SEARCH BY NAME
// 	searchResp, err := s.repo.SearchUsersByName(&pg.SearchUsers{FirstName: createUserResp.FirstName})
// 	s.Nil(err)
// 	s.NotNil(searchResp)

// 	// DELETE USER
// 	deleteUserResp, err := s.repo.DeleteUser(&pg.IdRequest{Id: createUserResp.Id})
// 	s.Nil(err)
// 	s.NotNil(deleteUserResp)
// }

// func (suite *UserSuiteTest) TearDownSuite() {
// 	suite.cleanUpFunc()
// }

// func TestUserRepositoryTestSuite(t *testing.T) {
// 	suite.Run(t, new(UserSuiteTest))
// }
