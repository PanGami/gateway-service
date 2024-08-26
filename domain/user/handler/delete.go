package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pangami/gateway-service/domain/user/client"
	"github.com/pangami/gateway-service/util/errors"

	pb "github.com/pangami/gateway-service/proto/user"
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

	id := c.QueryParam("id")

	// Convert the id string to an int32
	userId, err := strconv.Atoi(id)
	if err != nil {
		return errors.ErrBadRequest("Invalid user ID")
	}

	// Populate the protobuf request with data from the user.User struct
	req := &pb.DetailUserRequest{
		Id: int32(userId), // Explicitly convert to int32
	}

	grpcResp, err := client.UserDelete(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		log.Println("response", err.Error())
		return errors.Wrap(err, st.Code(), st.Message())
	}

	resp, err := h.buildResponse(grpcResp)
	if err != nil {
		return errors.Wrap(err, codes.Internal, "Failed to build response")
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

func NewUserDelete() *UserDelete {
	return &UserDelete{}
}
