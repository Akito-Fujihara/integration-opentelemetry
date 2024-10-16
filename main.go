package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Akito-Fujihara/integration-opentelemetry/config"
	"github.com/Akito-Fujihara/integration-opentelemetry/model"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

func main() {
	// この関数はtracer.goに移動
	tp, err := config.InitTracer()
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer func() {
		// アプリケーション終了時にトレーサーをシャットダウン
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shut down tracer: %v", err)
		}
	}()
	db, err := config.NewDB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(model.User{})
	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		panic(err)
	}

	repo := NewRepository(db)
	userHandler := &UserHandler{repo: repo}

	e := echo.New()
	e.Use(config.TraceMiddlewire(tp))

	// シンプルなハンドラー
	e.GET("/", Hello)
	e.POST("/users", userHandler.CreateUser)

	// サーバーを起動
	e.Logger.Fatal(e.Start(":8080"))
}

func Hello(c echo.Context) error {
	_, span := otel.Tracer("example-tracer").Start(c.Request().Context(), "hello world")
	defer span.End()
	return c.String(http.StatusOK, "Hello, World!")
}

type UserHandler struct {
	repo *repository
}

func (u *UserHandler) CreateUser(c echo.Context) error {
	_, span := otel.Tracer("example-tracer").Start(c.Request().Context(), "create user")
	defer span.End()

	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return c.String(http.StatusBadRequest, "request error")
	}

	if err := u.repo.insertUser(c, user); err != nil {
		return c.String(http.StatusInternalServerError, "insert error")
	}

	return c.JSON(http.StatusCreated, user)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) insertUser(c echo.Context, user *model.User) error {
	return r.db.WithContext(c.Request().Context()).Create(user).Error
}
