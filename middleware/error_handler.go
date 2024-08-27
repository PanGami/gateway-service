// gateway-service/middleware/error_handler.go
package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pangami/gateway-service/util/errors"
)

func ErrorHandlerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			appErr, ok := err.(*errors.AppError)
			if ok {
				return c.JSON(int(appErr.Code), map[string]string{"error": appErr.Message})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
		return nil
	}
}
