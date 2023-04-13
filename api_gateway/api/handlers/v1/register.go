package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/gin-gonic/gin"
	"gitlab.com/micro/api_gateway/api/handlers/models"
	"gitlab.com/micro/api_gateway/email"
	"gitlab.com/micro/api_gateway/genproto/user"
	"gitlab.com/micro/api_gateway/pkg/etc"
	"gitlab.com/micro/api_gateway/pkg/logger"
)

// register user
// @Summary		register user
// @Description	this registers user
// @Tags		Register
// @Accept		json
// @Produce 	json
// @Param 		body	body  	 models.RegisterUserModel true "Register user"
// @Success		201 	{object} models.Error
// @Failure		500 	{object} models.Error
// @Router		/v1/user/register 	[post]
func (h *handlerV1) RegisterUser(c *gin.Context) {
	var (
		body models.RegisterUserModel
	)

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error while bind json: ", logger.Error(err))
		return
	}

	body.Email = strings.TrimSpace(body.Email)
	body.UserName = strings.TrimSpace(body.UserName)
	body.Email = strings.ToLower(body.Email)
	body.UserName = strings.ToLower(body.UserName)
	body.Password, err = etc.HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("couldn't hash the password")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	existsEmail, err := h.serviceManager.UserService().CheckField(ctx, &user.CheckFieldReq{
		Field: "email",
		Value: body.Email,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed check email unique", logger.Error(err))
		return
	}

	if existsEmail.Exists {
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
				"info":  "please enter another email",
			})
			h.log.Error("This email already exists", logger.Error(err))
			return
		}
	}

	userToBeSaved := user.UserRequest{
		Id:        uuid.New().String(),
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  body.Password,
	}

	userToBeSaved.Code = etc.GenerateCode(6)
	msg := "Subject: Exam email verification\n Your verification code: " + userToBeSaved.Code
	err = email.SendEmail([]string{body.Email}, []byte(msg))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error":       nil,
			"Code":        http.StatusAccepted,
			"Description": "Your Email is not valid, Please recheck it",
		})
		return
	}

	userBodyByte, err := json.Marshal(userToBeSaved)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Error while marshar user ", logger.Error(err))
		return
	}

	err = h.redis.SetWithTTL(body.Email, string(userBodyByte), 300)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Error set to redis user body", logger.Error(err))
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"Error":       nil,
		"Code":        http.StatusAccepted,
		"Description": "Your request successfuly accepted we have send code to your email, Your code is : " + userToBeSaved.Code,
	})
}

// Verify user
// @Summary      Verify user
// @Description  Verifys user
// @Tags         Register
// @Accept       json
// @Produce      json
// @Param        email  path string true "email"
// @Param        code   path string true "code"
// @Success      200  {object}  models.VerifyResponse
// @Router      /v1/verify/{email}/{code} [get]
func (h *handlerV1) VerifyUser(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	var (
		email = c.Param("email")
		code  = c.Param("code")
	)
	fmt.Println(email)
	userBody, err := h.redis.Get(email)
	if err != nil {
		c.JSON(http.StatusGatewayTimeout, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Error while getting user from redis", logger.Any("redis", err))
	}
	if userBody == nil {
		c.JSON(http.StatusGatewayTimeout, gin.H{
			"info": "Your time has expired",
		})
		return
	}

	userBodys := cast.ToString(userBody)
	body := user.UserRequest{}

	err = json.Unmarshal([]byte(userBodys), &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Error while unmarshaling from json to user body", logger.Any("json", err))
		return
	}

	if body.Code != code {
		fmt.Println(body.Code)
		c.JSON(http.StatusConflict, gin.H{
			"info": "Wrong code",
		})
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		panic("Can't generate uuid")
	}

	body.Id = id.String()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	// Genrating refresh and jwt tokens
	h.jwthandler.Iss = "user"
	h.jwthandler.Sub = body.Id
	h.jwthandler.Role = "authorized"
	h.jwthandler.Aud = []string{"some-app-name"}
	h.jwthandler.SigninKey = h.cfg.SignInKey
	h.jwthandler.Log = h.log
	tokens, err := h.jwthandler.GenerateAuthJWT()
	accessToken := tokens[0]
	refreshToken := tokens[1]

	if err != nil {
		h.log.Error("error occured while generating tokens")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong,please try again",
		})
		return
	}

	body.RefreshToken = refreshToken
	res, err := h.serviceManager.UserService().CreateUser(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Error while creating user", logger.Any("post", err))
		return
	}
	response := &models.VerifyResponse{
		Id:           res.Id,
		FirstName:    body.FirstName,
		LastName:     body.LastName,
		Email:        res.Email,
		Password:     res.Password,
		RefreshToken: res.RefreshToken,
		AccessToken:  res.AccessToken,
	}

	response.AccessToken = accessToken
	response.RefreshToken = refreshToken

	c.JSON(http.StatusOK, response)
}
