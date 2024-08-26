package route

import (
	"github.com/labstack/echo/v4"
	// auth "github.com/pangami/domain/auth/handler"
	user "github.com/pangami/gateway-service/domain/user/handler"
)

// Handler endpoint to use it later
type Handler interface {
	Handle(c echo.Context) (err error)
}

var endpoint = map[string]Handler{
	//auth
	// "login":  auth.NewLogin(),
	// "logout": auth.NewLogout(),

	//user
	"user_list":   user.NewUserList(),
	"user_create": user.NewUserCreate(),
	"user_detail": user.NewUserDetail(),
	"user_update": user.NewUserUpdate(),
	"user_delete": user.NewUserDelete(),

	//userActivity
	"user_activity": user.NewUserGetActivity(),
}
