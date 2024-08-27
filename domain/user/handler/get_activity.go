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

type UserGetActivity struct{}

func (h *UserGetActivity) Handle(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.QueryParam("id")

	// Convert the id string to an int32
	userId, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.ErrBadRequest("Invalid user ID"))
	}

	req := &pb.DetailUserRequest{
		Id: int32(userId), // Explicitly convert to int32
	}

	grpcResp, err := client.UserGetActivity(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		log.Println("response", err.Error())

		// Handle not found error
		if st.Code() == codes.NotFound {
			return c.JSON(http.StatusNotFound, errors.ErrNotFound("User not found"))
		}

		return c.JSON(http.StatusInternalServerError, errors.Wrap(err, st.Code(), st.Message()))
	}

	resp, err := h.buildResponse(grpcResp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.Wrap(err, codes.Internal, "Failed to build response"))
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *UserGetActivity) buildResponse(response *pb.UserActivitiesResponse) (*util.Response, error) {
	resp := &util.Response{
		Status:  "true",
		Code:    util.Success,
		Message: util.StatusMessage[util.Success],
		Data:    response,
	}
	return resp, nil
}

func NewUserGetActivity() *UserGetActivity {
	return &UserGetActivity{}
}
