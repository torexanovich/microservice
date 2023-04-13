package v1

import (
	"context"

	"net/http"
	"time"

	"gitlab.com/micro/api_gateway/api/handlers/models"
	"gitlab.com/micro/api_gateway/genproto/user"
	"gitlab.com/micro/api_gateway/pkg/etc"

	"github.com/gin-gonic/gin"
	"gitlab.com/micro/api_gateway/pkg/logger"
)

// User Login
// @Summary 		Login User
// @Description 	This function get login user
// @Tags 			Register
// @Accept 			json
// @Produce			json
// @Param 			email 		path string true "email"
// @Param 			password 	path string true "password"
// @Success 		200 {object} 	models.LoginResponse
// @Failure			500 {object} 	models.Error
// @Failure			400 {object} 	models.Error
// @Router			/v1/login/{email}/{password} [get]
func (h *handlerV1) Login(c *gin.Context) {
	var (
		email    = c.Param("email")
		password = c.Param("password")
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	res, err := h.serviceManager.UserService().GetByEmail(ctx, &user.EmailReq{Email: email})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error":       err,
			"Description": "Couln't find matching information, Have you registered before?",
		})
		h.log.Error("Error while getting user by email", logger.Any("post", err))
		return
	}

	if !etc.CheckPasswordHash(password, res.Password) {
		c.JSON(http.StatusNotFound, gin.H{
			"Description": "Password or Email error",
			"Code":        http.StatusBadRequest,
		})
		return
	}

	h.jwthandler.Iss = "user"
	h.jwthandler.Sub = res.Id
	h.jwthandler.Role = "authorized"
	h.jwthandler.Aud = []string{"test-app"}
	h.jwthandler.SigninKey = h.cfg.SignInKey
	h.jwthandler.Log = h.log
	tokens, err := h.jwthandler.GenerateAuthJWT()
	accesstoken := tokens[0]
	refreshToken := tokens[1]
	if err != nil {
		h.log.Error("error occured while generating tokens")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong,please try again",
		})
		return
	}

	res.AccessToken = accesstoken
	res.RefreshToken = refreshToken
	res.Password = ""

	response := models.LoginResponse{
		Id:           res.Id,
		FirstName:    res.FirstName,
		LastName:     res.LastName,
		Email:        res.Email,
		Password:     res.Password,
		RefreshToken: res.RefreshToken,
	}
	c.JSON(http.StatusOK, response)
}
