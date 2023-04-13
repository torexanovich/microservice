package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gitlab.com/micro/api_gateway/api/token"
	"gitlab.com/micro/api_gateway/config"
)

type JWTRoleAuthorizer struct {
	enforcer   *casbin.Enforcer
	cfg        config.Config
	jwtHandler token.JWTHandler
}

func NewAuthorizer(e *casbin.Enforcer, jwtHandler token.JWTHandler, cfg config.Config) gin.HandlerFunc {
	a := &JWTRoleAuthorizer{
		enforcer:   e,
		cfg:        cfg,
		jwtHandler: jwtHandler,
	}

	return func(c *gin.Context) {
		allow, err := a.CheckPermission(c.Request)
		fmt.Println(allow, err, "<><>")
		if err != nil {
			fmt.Println(err, "<<CHECK")
			v, _ := err.(*jwt.ValidationError)
			if v.Errors == jwt.ValidationErrorExpired {
				a.ReqireRefresh(c)
			} else {
				a.RequirePermission(c)
			}
		} else if !allow {
			a.RequirePermission(c)
		}
	}
}

func (a *JWTRoleAuthorizer) CheckPermission(r *http.Request) (bool, error) {
	user, err := a.GetRole(r)
	fmt.Println("<>ERR", user, err)
	if err != nil {
		return false, err
	}

	method := r.Method
	path := r.URL.Path

	allowed, err := a.enforcer.Enforce(user, path, method)
	if err != nil {
		panic(err)
	}
	fmt.Println(allowed)
	return allowed, nil
}

func (a *JWTRoleAuthorizer) GetRole(r *http.Request) (string, error) {
	var (
		role   string
		claims jwt.MapClaims
		err    error
	)

	jwtToken := r.Header.Get("Authorization")
	if jwtToken == "" {
		return "unauthorized", nil
	} else if strings.Contains(jwtToken, "Basic") {
		return "unauthorized", nil
	}

	a.jwtHandler.Token = jwtToken

	claims, err = a.jwtHandler.ExtractClaims()
	if err != nil {
		return "", err
	}

	if claims["role"].(string) == "authorized" {
		role = "authorized"
	} else {
		role = "unknown"
	}

	return role, nil
}

func (a *JWTRoleAuthorizer) ReqireRefresh(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"Status":  "UNAUTHORIZED",
		"Message": "Token is expired",
	})
	c.AbortWithStatus(401)
}

func (a *JWTRoleAuthorizer) RequirePermission(c *gin.Context) {
	c.AbortWithStatus(403)
}
