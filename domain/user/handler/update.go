package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pangami/gateway-service/domain/user"
	"github.com/pangami/gateway-service/domain/user/client"
	// "github.com/pangami/gateway-service/route/middleware"
	// "github.com/pangami/gateway-service/util"
)

type UserUpdate struct{}

func (h *UserUpdate) Handle(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	// _, ok := c.Get(util.ContextTokenValueKey).(middleware.ValidateTokenResponse)
	// if !ok {
	// 	resp := &user.Response{
	// 		Code:    "500",
	// 		Message: util.StatusMessage[http.StatusInternalServerError],
	// 		Status:  false,
	// 		Data:    nil,
	// 	}
	// 	return c.JSON(http.StatusInternalServerError, &resp)
	// }

	r := new(user.User)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("Error converting ID to integer:", err.Error())
		return err
	}
	r.ID = id

	err = h.validate(r, c)
	if err != nil {
		log.Println("validate error : ", err.Error())
		return err
	}

	result, err := client.UserUpdate(ctx, r)
	if err != nil {
		log.Println("response", err.Error())
		resp, err := h.buildResponse(result)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusInternalServerError, resp)
	}
	resp, err := h.buildResponse(result)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *UserUpdate) buildResponse(response map[string]interface{}) (*user.Response, error) {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshaling map to JSON:", err)
		return nil, err
	}

	var resp user.Response
	err = json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return nil, err
	}

	return &resp, nil
}

func (h *UserUpdate) validate(r *user.User, c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	return c.Validate(r)
}

func NewUserUpdate() *UserUpdate {
	return &UserUpdate{}
}
