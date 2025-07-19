package middleware

import (
	"log"
	"net/http"
	"reflect"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/wendisx/gorchat/internal/constant"
	"github.com/wendisx/gorchat/internal/validator"
)

type Middleware interface {
	ValidatorMiddleware(v any) echo.MiddlewareFunc
	SessionCheckMiddleware(allowNew bool) echo.MiddlewareFunc
}

type middleware struct {
	va    *validator.Validator
	store sessions.Store
}

func NewMiddleware(va *validator.Validator, store sessions.Store) Middleware {
	defer log.Printf("[init] -- (config/middleware) status: success\n")
	return &middleware{
		va:    va,
		store: store,
	}
}

func (md *middleware) ValidatorMiddleware(v any) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			body := reflect.New(reflect.TypeOf(v).Elem()).Interface()
			if err := c.Bind(body); err != nil {
				log.Printf("[middleware] -- (validator) %v\n", err.Error())
				return echo.NewHTTPError(http.StatusBadRequest, constant.MsgBadRequest)
			}
			if errs := md.va.Check(body); errs != nil {
				log.Printf("[middleware] -- (validator) %v\n", errs.Error())
				return echo.NewHTTPError(http.StatusBadRequest, constant.MsgValidateFail)
			}
			// 设置校验对象为 key = body value = body
			c.Set("body", body)
			log.Printf("[middleware] -- (validator) %v\n", constant.MsgValidateSuccess)
			return next(c)
		}
	}
}

func (md *middleware) SessionCheckMiddleware(allowNew bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session, err := md.store.Get(c.Request(), constant.SESSION_KEY)
			if allowNew {
				err = md.store.Save(c.Request(), c.Response().Writer, session)
				if err != nil {
					log.Printf("[middleware] -- (sessioncheck) %v\n", err.Error())
					return echo.NewHTTPError(http.StatusInternalServerError, constant.MsgServerInternalErr)
				}
				log.Printf("[middleware] -- (sessioncheck) status: bypass\n")
			} else {
				if !session.IsNew && err == nil {
					log.Printf("[middleware] -- (sessioncheck) status: bypass\n")
					return next(c)
				} else {
					err = md.store.Save(c.Request(), c.Response().Writer, session)
					if err != nil {
						log.Printf("[middleware] -- (sessioncheck) %v\n", err.Error())
						return echo.NewHTTPError(http.StatusInternalServerError, constant.MsgServerInternalErr)
					}
					return echo.NewHTTPError(http.StatusUnauthorized, constant.MsgNotAuthenticate)
				}

			}
			return next(c)
		}
	}
}
