package v1

import (
	"context"
	"net/http"
	"time"

	pu "gitlab.com/micro/api_gateway/genproto/user"
	l "gitlab.com/micro/api_gateway/pkg/logger"
	"gitlab.com/micro/api_gateway/pkg/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// // CreateUser...
// // @Summary CreateUser
// // @Description This API for creating a new user
// // @Tags user
// // @Accept  json
// // @Produce  json
// // @Param body body models.UserRequest true "UserRequest"
// // @Success 200 {object} models.User
// // @Failure 400 {object} models.Error
// // @Router /v1/user/ [post]
// func (h *handlerV1) CreateUser(c *gin.Context) {
// 	var (
// 		body        pu.UserRequest
// 		jspbMarshal protojson.MarshalOptions
// 	)
// 	jspbMarshal.UseProtoNames = true

// 	err := c.ShouldBindJSON(&body)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to bind json", l.Error(err))
// 		return
// 	}

// 	response, err := h.serviceManager.UserService().CreateUser(context.Background(), &body)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to create user", l.Error(err))
// 		return
// 	}

// 	c.JSON(http.StatusCreated, response)
// }

// GetUser ...
// @Summary GetUser
// @Description This API for getting user detail
// @Tags user
// @Accept  json
// @Produce  json
// @Param uuid path string true "uuid"
// @Success 200 {object} models.User
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/user [get]
func (h *handlerV1) GetUser(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	claims := GetClaims(h, c)

	id := claims["sub"].(string)

	response, err := h.serviceManager.UserService().GetUserById(context.Background(), &pu.IdRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user by id", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetUsers ...
// @Summary GetAllUsers
// @Description This API for getting all users
// @Tags user
// @Accept  json
// @Produce  json
// @Param limit path int true "limit"
// @Param page path int true "page"
// @Success 200 {object} models.Users
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/users/ [get]
func (h *handlerV1) GetAllUsers(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	params, errStr := utils.ParseQueryParams(queryParams)
	if errStr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errStr[0],
		})
		h.log.Error("failed to parse query params to json: " + errStr[0])
		return
	}

	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	response, err := h.serviceManager.UserService().GetAllUsers(context.Background(), &pu.AllUsersRequest{Limit: params.Limit, Page: params.Page})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get all users", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Update User
// @Description This API for updating user
// @Tags user
// @Accept  json
// @Produce  json
// @Param body body models.UpdateUserReq true "UpdateUsersReq"
// @Success 200 {object} models.Empty
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/user/{id} [put]
func (h *handlerV1) UpdateUser(c *gin.Context) {
	var (
		body        pu.UpdateUserRequest
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind JSON", l.Error(err))
		return
	}

	response, err := h.serviceManager.UserService().UpdateUser(context.Background(), &pu.UpdateUserRequest{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Id:        body.Id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Delete User
// @Description This API for deleting user
// @Tags user
// @Accept  json
// @Produce  json
// @Param id path string true "Id"
// @Success 200 {object} models.User
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/user/{id} [delete]
func (h *handlerV1) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	jspbMarshal := protojson.MarshalOptions{}
	jspbMarshal.UseProtoNames = true

	response, err := h.serviceManager.UserService().DeleteUser(context.Background(), &pu.IdRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}


// create mod
// @Summary create mod
// @Description create mod by user id
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "user id" 
// @Succes 200 {object} models.Empty
// @Router /v1/admin/create_mod/{id} [patch]	
func (h *handlerV1) CreateMod(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	uuid := c.Param("id")
	ctx, cancel :=  context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	resp, err := h.serviceManager.UserService().CreateMod(ctx, &pu.IdRequest{Id: uuid})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create moderator", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}