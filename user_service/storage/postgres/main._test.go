package postgres

import (
	"log"
	"os"
	"testing"

	"gitlab.com/micro/user_service/config"
	"gitlab.com/micro/user_service/pkg/db"
	"gitlab.com/micro/user_service/pkg/logger"
)

var pgRepo *UserRepo

func TestMain(m *testing.M) {
	conf := config.Load()

	connDB, err := db.ConnectToDB(conf)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}

	pgRepo = NewUserRepo(connDB)

	os.Exit(m.Run())
}
