package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pangami/gateway-service/domain/user"
	"github.com/pangami/gateway-service/domain/user/client"

	pb "github.com/pangami/gateway-service/proto/user"
	// "github.com/pangami/gateway-service/route/middleware"
	"github.com/pangami/gateway-service/util"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserDelete struct{}

func (h *UserDelete) Handle(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	r := new(user.UserDetailRequest)
	err := h.validate(r, c)
	if err != nil {
		log.Println("validate error : ", err.Error())
		resp := &user.Response{
			Code:    400,
			Message: util.StatusMessage[util.InvalidArgument],
			Status:  false,
			Data:    nil,
		}
		return c.JSON(http.StatusBadRequest, &resp)
	}

	// Populate the protobuf request with data from the user.User struct
	req := &pb.DetailUserRequest{
		Id: r.ID,
	}

	grpcResp, err := client.UserDelete(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		log.Println("response", err.Error())
		resp, err := h.buildErrorResponse(st.Code(), st.Message())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusInternalServerError, resp)
	}

	resp, err := h.buildResponse(grpcResp)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *UserDelete) buildResponse(response *pb.NoResponse) (*util.Response, error) {
	resp := &util.Response{
		Status:  "true",
		Code:    util.Success,
		Message: util.StatusMessage[util.Success],
		Data:    response,
	}
	return resp, nil
}

func (h *UserDelete) buildErrorResponse(errorCode codes.Code, message string) (*util.Response, error) {
	resp := &util.Response{
		Status:  "false",
		Code:    errorCode,
		Message: message,
		Data:    nil,
	}
	return resp, nil
}

func (h *UserDelete) validate(r *user.UserDetailRequest, c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	return c.Validate(r)
}

func NewUserDelete() *UserDelete {
	return &UserDelete{}
}
