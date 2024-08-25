package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pangami/gateway-service/domain/user/client"

	pb "github.com/pangami/gateway-service/proto/user"
	"github.com/pangami/gateway-service/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserList struct{}

func (h *UserList) Handle(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	pageStr := c.QueryParam("page")
	pageSizeStr := c.QueryParam("page_size")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1 // Default to page 1 if not provided or invalid
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10 // Default to 10 if not provided or invalid
	}

	req := &pb.ListUsersRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	grpcResp, err := client.UserList(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		log.Println("ListUsers gRPC error:", err.Error())
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

func (h *UserList) buildResponse(response *pb.ListUsersResponse) (*util.Response, error) {
	resp := &util.Response{
		Status:  "true",
		Code:    util.Success,
		Message: util.StatusMessage[util.Success],
		Data:    response,
	}
	return resp, nil
}

func (h *UserList) buildErrorResponse(errorCode codes.Code, message string) (*util.Response, error) {
	resp := &util.Response{
		Status:  "false",
		Code:    errorCode,
		Message: message,
		Data:    nil,
	}
	return resp, nil
}

func NewUserList() *UserList {
	return &UserList{}
}
