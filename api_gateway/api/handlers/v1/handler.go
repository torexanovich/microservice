package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gitlab.com/micro/api_gateway/api/handlers/models"
	t "gitlab.com/micro/api_gateway/api/token"
	"gitlab.com/micro/api_gateway/config"
	"gitlab.com/micro/api_gateway/pkg/logger"
	"gitlab.com/micro/api_gateway/services"
	"gitlab.com/micro/api_gateway/storage/repo"
)

type handlerV1 struct {
	log            logger.Logger
	serviceManager services.IServiceManager
	cfg            config.Config
	redis          repo.RedisRepo
	jwthandler     t.JWTHandler
}

type HandlerV1Config struct {
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Cfg            config.Config
	Redis          repo.RedisRepo
	JWTHandler     t.JWTHandler
}

func New(c *HandlerV1Config) handlerV1 {
	return handlerV1{
		log:            c.Logger,
		serviceManager: c.ServiceManager,
		cfg:            c.Cfg,
		redis:          c.Redis,
		jwthandler:     c.JWTHandler,
	}
}

func GetClaims(h *handlerV1, c *gin.Context) jwt.MapClaims {
	var (
		ErrUnauthorized = errors.New("unauthorized")
		authorization   models.GetProfileByJwtRequest
		claims          jwt.MapClaims
		err             error
	)

	authorization.Token = c.GetHeader("Authorization")
	if c.Request.Header.Get("Authorization") == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Status":  "ErrorCodeUnauthorized",
			"Message": "Unauthorized request",
		})
		h.log.Error("Unauthorized request: ", logger.Error(ErrUnauthorized))
		return nil
	}

	authorization.Token = strings.TrimSpace(strings.Trim(authorization.Token, "Bearer"))

	h.jwthandler.Token = authorization.Token
	claims, err = h.jwthandler.ExtractClaims()
	fmt.Println(claims)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Status":  "extractClaimsFail",
			"Message": "Unauthorized request",
		})
		h.log.Error("Unauthorized request: ", logger.Error(err))
		return nil
	}

	return claims
}
