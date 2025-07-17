package model

import (
	"log"

	"github.com/labstack/echo/v4"

	"github.com/wendisx/gorchat/internal/constant"
)

type Response interface {
	Success(ctx echo.Context, httpCode int, msg string, data any) error
	Fail(ctx echo.Context, httpCode int, code int, msg string) error
}

type response struct{}

func NewResponser() Response {
	defer log.Printf("[init] -- (model/response) status: success\n")
	return &response{}
}

// 成功响应
func (res *response) Success(ctx echo.Context, httpCode int, msg string, data any) error {
	return ctx.JSON(httpCode, struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    any    `json:"data"`
	}{
		Code:    constant.SUCCESS_CODE,
		Message: msg,
		Data:    data,
	})
}

// 失败响应
func (res *response) Fail(ctx echo.Context, httpCode int, code int, msg string) error {
	return ctx.JSON(httpCode, struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    any    `json:"data"`
	}{
		Code:    code,
		Message: msg,
		Data:    nil,
	})
}
