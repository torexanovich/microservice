package main

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	defaultrolemanager "github.com/casbin/casbin/v2/rbac/default-role-manager"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/gomodule/redigo/redis"
	"gitlab.com/micro/api_gateway/api"
	"gitlab.com/micro/api_gateway/config"
	"gitlab.com/micro/api_gateway/pkg/logger"
	"gitlab.com/micro/api_gateway/services"
	r "gitlab.com/micro/api_gateway/storage/redis"
)

func main() {
	var casbinEnforcer *casbin.Enforcer
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "api_gateway")

	psqlString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
	)

	a, err := gormadapter.NewAdapter("postgres", psqlString, true)
	if err != nil {
		log.Error("new adapter error", logger.Error(err))
		return
	}
	fmt.Println(a)

	casbinEnforcer, err = casbin.NewEnforcer("./config/auth.conf", a)
	if err != nil {
		log.Error("casbin enforcer error", logger.Error(err))
		return
	}

	err = casbinEnforcer.LoadPolicy()
	if err != nil {
		log.Error("casbin error load policy", logger.Error(err))
		return
	}

	serviceManager, err := services.NewServiceManager(&cfg)
	if err != nil {
		log.Error("gRPC dial error: ", logger.Error(err))
	}

	casbinEnforcer.GetRoleManager().(*defaultrolemanager.RoleManager).AddMatchingFunc("keyMatch", util.KeyMatch)
	casbinEnforcer.GetRoleManager().(*defaultrolemanager.RoleManager).AddMatchingFunc("keyMatch3", util.KeyMatch3)

	pool := redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort))
			if err != nil {
				panic(err)
			}
			return c, err
		},
	}

	server := api.New(api.Option{
		Conf:            cfg,
		ServiceManager:  serviceManager,
		Logger:          log,
		InMemoryStorage: r.NewRedisRepo(&pool),
		CasbinEnforcer:  casbinEnforcer,
	})

	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("failed to run HTTP server: ", logger.Error(err))
		panic(err)
	}
}
