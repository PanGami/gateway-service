package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pangami/gateway-service/domain/auth"
	"github.com/pangami/gateway-service/domain/auth/client"
	util "github.com/pangami/gateway-service/util"
)

type Logout struct{}

func (h *Logout) Handle(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	headerToken := c.Request().Header.Get(util.ApiKey)

	r := new(auth.LogoutRequest)
	err := h.validate(r, c)
	if err != nil {
		log.Println("validate error : ", err.Error())
		return err
	}

	r.Token = headerToken
	log.Println("UserId ", r.UserId)
	log.Println("Token ", r.Token)

	result, err := client.Logout(ctx, r)
	if err != nil {
		log.Println("response", err.Error())
		resp, err := h.buildResponse(result)
		if err != nil {
			return err
		}
		// return resp.JSON(c)
		return c.JSON(http.StatusInternalServerError, resp)
	}
	resp, err := h.buildResponse(result)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Logout) buildResponse(response map[string]interface{}) (*auth.ClientResponse, error) {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshaling map to JSON:", err)
		return nil, err
	}

	var resp auth.ClientResponse
	err = json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return nil, err
	}

	return &resp, nil
}

func (h *Logout) validate(r *auth.LogoutRequest, c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	return c.Validate(r)
}

func NewLogout() *Logout {
	return &Logout{}
}
