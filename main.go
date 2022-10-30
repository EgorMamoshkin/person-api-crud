package main

import (
	"fmt"
	"github.com/EgorMamoshkin/person-api-crud/internal/config"
	handlers "github.com/EgorMamoshkin/person-api-crud/internal/http"
	"github.com/EgorMamoshkin/person-api-crud/internal/logic"
	"github.com/EgorMamoshkin/person-api-crud/internal/postgres"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"time"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		logrus.Fatal(err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s?sslmode=disable", cfg.DBUser, cfg.DBPass, cfg.DBPath)

	db := postgres.NewPostgresRepo(dsn)

	perLogic := logic.NewPersonLogic(db, 5*time.Second)

	e := echo.New()

	handlers.NewPersonHandler(e, perLogic)

	logrus.Fatal(e.Start(cfg.ApiServAddr))
}
