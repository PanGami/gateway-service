package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pangami/gateway-service/domain/user"
	"github.com/pangami/gateway-service/domain/user/client"
	"github.com/pangami/gateway-service/util/errors"

	pb "github.com/pangami/gateway-service/proto/user"
	"github.com/pangami/gateway-service/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserCreate struct{}

func (h *UserCreate) Handle(c echo.Context) error {
    ctx := c.Request().Context()
    if ctx == nil {
        ctx = context.Background()
    }

    r := new(user.User)
    err := h.validate(r, c)
    if err != nil {
        log.Println("validate error : ", err.Error())
        return errors.ErrBadRequest("Invalid request data")
    }

    // Populate the protobuf request with data from the user.User struct
    req := &pb.CreateUserRequest{
        Username: r.Username,
        FullName: r.FullName,
        Password: r.Password,
    }

    grpcResp, err := client.UserCreate(ctx, req)
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

func (h *UserCreate) buildResponse(response *pb.NoResponse) (*util.Response, error) {
    resp := &util.Response{
        Status:  "true",
        Code:    util.Success,
        Message: util.StatusMessage[util.Success],
        Data:    response,
    }
    return resp, nil
}

func (h *UserCreate) validate(r *user.User, c echo.Context) error {
    if err := c.Bind(r); err != nil {
        return err
    }

    return c.Validate(r)
}

func NewUserCreate() *UserCreate {
    return &UserCreate{}
}