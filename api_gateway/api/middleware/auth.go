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

func NewAuth(enforce *casbin.Enforcer, jwtHandler token.JWTHandler, cfg config.Config) gin.HandlerFunc {
	a := &JWTRoleAuthorizer{
		enforcer:   enforce,
		cfg:        cfg,
		jwtHandler: jwtHandler,
	}

	return func(c *gin.Context) {
		allow, err := a.CheckPermission(c.Request)
		fmt.Println(allow)
		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			if v.Errors == jwt.ValidationErrorExpired {
				a.RequireRefresh(c)
			} else {
				a.RequirePermission(c)
			}
		} else if !allow {
			a.RequirePermission(c)
		}
	}
}

// CheckPermission checks whether user is allowed to use certain endpoint
func (a *JWTRoleAuthorizer) CheckPermission(r *http.Request) (bool, error) {
	user, err := a.GetRole(r)
	if err != nil {
		return false, err
	}
	fmt.Println(user)
	method := r.Method
	path := r.URL.Path
	fmt.Println(r.Method)
	fmt.Println(r.URL.Path)

	allowed, err := a.enforcer.Enforce(user, path, method)
	if err != nil {
		panic(err)
	}

	return allowed, nil
}

// GetRole gets role from Authorization header if there is a token then it is
// parsed and in role got from role claim. If there is no token then role is
// unauthorized
func (a *JWTRoleAuthorizer) GetRole(r *http.Request) (string, error) {
	var (
		role   string
		claims jwt.MapClaims
		err    error
	)

	jwtToken := r.Header.Get("Authorization")
	fmt.Println("200:", jwtToken)
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
	} else if claims["role"].(string) == "sudo" {
		role = "sudo"
	} else if claims["role"].(string) == "admin" {
		role = "admin"
	} else {
		role = "unknown"
	}
	return role, nil
}

// RequirePermission aborts request with 403 status
func (a *JWTRoleAuthorizer) RequirePermission(c *gin.Context) {
	c.AbortWithStatus(403)
}

// RequireRefresh aborts request with 401 status
func (a *JWTRoleAuthorizer) RequireRefresh(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"Error":   "smth went wrong",
		"Status":  "UNAUTHORIZED",
		"Message": "Token is expired",
	})
	c.AbortWithStatus(401)
}
