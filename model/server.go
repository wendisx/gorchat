package model

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/wendisx/gorchat/internal/log"
)

type Dependency struct {
	Echo        *echo.Echo
	Database    *sql.DB
	Logger      log.Logger
	RedisClient *redis.Client
	Response    Response
}
