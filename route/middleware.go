package route

import (
	"github.com/labstack/echo/v4"
	"github.com/pangami/gateway-service/route/middleware"
)

var middlewareHandler = map[string]echo.MiddlewareFunc{
	"validate_token": middleware.AuthMiddleware,
}
