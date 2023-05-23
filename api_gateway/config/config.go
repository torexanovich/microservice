package config

import (
	"os"

	"github.com/spf13/cast"
)
 
type Config struct {
	Environment string // develop, staging, production

	// context timeout in seconds
	CtxTimeout int

	LogLevel string
	HTTPPort string

	UserServiceHost    string
	UserServicePort    string
	PostServiceHost    string
	PostServicePort    string
	CommentServiceHost string
	CommentServicePort string
	RedisHost          string
	RedisPort          int
	SignInKey          string
	CSVFilePath        string
	AuthConfigPath     string
	PostgresUser       string
	PostgresPassword   string
	PostgresHost       string
	PostgresDatabase   string
	PostgresPort       string
}

func Load() Config {
	c := Config{}

	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))
	c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "smn"))
	c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "123"))
	c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "database"))
	c.PostgresPort = cast.ToString(getOrReturnDefault("POSTGRES_PORT", "5434"))
	c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "smn"))
	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.HTTPPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":5050"))

	c.UserServiceHost = cast.ToString(getOrReturnDefault("USER_SERVICE_HOST", "user_service"))
	c.UserServicePort = cast.ToString(getOrReturnDefault("USER_SERVICE_PORT", "8000"))

	c.PostServiceHost = cast.ToString(getOrReturnDefault("POST_SERVICE_HOST", "post_service"))
	c.PostServicePort = cast.ToString(getOrReturnDefault("POST_SERVICE_PORT", "8010"))

	c.CommentServiceHost = cast.ToString(getOrReturnDefault("COMMENT_SERVICE_HOST", "comment_service"))
	c.CommentServicePort = cast.ToString(getOrReturnDefault("COMMENT_SERVICE_PORT", "8020"))

	c.RedisHost = cast.ToString(getOrReturnDefault("REDIS_HOST", "redisdb"))
	c.RedisPort = cast.ToInt(getOrReturnDefault("REDIS_PORT", 6379))
	c.SignInKey = cast.ToString(getOrReturnDefault("SIGNING_KEY", "123"))
	c.CtxTimeout = cast.ToInt(getOrReturnDefault("CTX_TIMEOUT", 7))
	c.CSVFilePath = cast.ToString(getOrReturnDefault("CSV_FILE_PATH", "./config/casbin.csv"))
	c.AuthConfigPath = cast.ToString(getOrReturnDefault("AUTH_PATH", "./config/auth.conf"))
	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
