package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/pangami/gateway-service/domain/auth/client"
	"github.com/pangami/gateway-service/util"

	"github.com/labstack/echo/v4"
)

// Define your data structures
type AuthAccessToken struct {
	ID        string    `json:"id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type EmptyObject struct{}

type ValidateTokenResponse struct {
	Status  bool   `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    Data   `json:"data"`
}

type Data struct {
	ID           string      `json:"id"`
	Username     string      `json:"username"`
	FullName     string      `json:"full_name"`
	MerchantID   string      `json:"merchant_id"`
	MerchantName string      `json:"merchant_name"`
	RegisterDate string      `json:"register_date"`
	FcmToken     string      `json:"fcm_token"`
	SecretKey    interface{} `json:"secret_key"`
	IsFirstLogin int         `json:"is_first_login"`
	IsLogin      int         `json:"is_login"`
	IsBod        string      `json:"is_bod"`
	StatusID     string      `json:"status_id"`
	Note         string      `json:"note"`
	EsignID      interface{} `json:"esign_id"`
}

// AuthMiddleware checks if the user is authenticated and decodes the JWT token
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		if ctx == nil {
			ctx = context.Background()
		}
		apiToken := c.Request().Header.Get("x-api-token")
		if apiToken != "" {
			c.Request().Header.Set("Authorization", "Bearer "+apiToken)
		}

		validateTokenResult, err := client.ValidateToken(ctx, apiToken)

		if err != nil {
			log.Printf("JWT decoding error: %v", err)
			response := ErrorResponseWithTraceID("", "006", "Invalid Access Token")
			return c.JSON(http.StatusUnauthorized, response)
		}

		validateTokenResponse, _ := convertToValidateTokenResponse(validateTokenResult)
		if validateTokenResponse.Code != "200" {
			log.Printf("JWT decoding error: %v", err)
			response := ErrorResponseWithTraceID("", "006", "Invalid Access Token")
			return c.JSON(http.StatusUnauthorized, response)
		}

		jwtClaims, err := DecodeJWT(apiToken)
		if err != nil {
			log.Printf("JWT decoding error: %v", err)
			response := ErrorResponseWithTraceID("", "006", err.Error())
			return c.JSON(http.StatusUnauthorized, response)
		}

		expirationTime := time.Unix(int64(jwtClaims.Exp), 0)
		if expirationTime.Before(time.Now()) {
			response := ErrorResponseWithTraceID("", "006", "Token Expired")
			return c.JSON(http.StatusUnauthorized, response)
		}

		// TODO wrapping context
		c.Set(util.ContextTokenValueKey, validateTokenResponse)
		c.Set(util.ContextJwtClaimKey, jwtClaims)
		return next(c)
	}
}

// ErrorResponseWithTraceID generates a formatted error response
func ErrorResponseWithTraceID(traceID, code, message string) map[string]interface{} {
	return map[string]interface{}{
		"traceId": traceID,
		"status":  false,
		"code":    code,
		"message": message,
		"data":    EmptyObject{},
	}
}

// Convert map to ValidateTokenResponse struct
func convertToValidateTokenResponse(validateTokenResult map[string]interface{}) (ValidateTokenResponse, error) {
	if validateTokenResult == nil {
		return ValidateTokenResponse{}, errors.New("validateTokenResult map is nil")
	}

	userDataMap, ok := validateTokenResult["data"].(map[string]interface{})
	if !ok {
		return ValidateTokenResponse{}, errors.New("data field is not a map")
	}

	// Create a new instance of ValidateTokenResponse
	var validateTokenResponse ValidateTokenResponse

	// Populate the fields of ValidateTokenResponse
	validateTokenResponse.Status = validateTokenResult["status"].(bool)
	validateTokenResponse.Code = validateTokenResult["code"].(string)
	validateTokenResponse.Message = validateTokenResult["message"].(string)

	// Populate the "Data" field of ValidateTokenResponse
	if userDataMap != nil {
		validateTokenResponse.Data.ID = userDataMap["id"].(string)
		validateTokenResponse.Data.Username = userDataMap["username"].(string)
		validateTokenResponse.Data.FullName = userDataMap["full_name"].(string)
		validateTokenResponse.Data.MerchantID = userDataMap["merchant_id"].(string)
		validateTokenResponse.Data.MerchantName = userDataMap["merchant_name"].(string)
		validateTokenResponse.Data.RegisterDate = userDataMap["register_date"].(string)
		validateTokenResponse.Data.FcmToken = userDataMap["fcm_token"].(string)
		validateTokenResponse.Data.SecretKey = userDataMap["secret_key"]
		validateTokenResponse.Data.IsFirstLogin = int(userDataMap["is_first_login"].(float64))
		validateTokenResponse.Data.IsLogin = int(userDataMap["is_login"].(float64))
		validateTokenResponse.Data.IsBod = userDataMap["is_bod"].(string)
		validateTokenResponse.Data.StatusID = userDataMap["status_id"].(string)
		validateTokenResponse.Data.Note = userDataMap["note"].(string)
		validateTokenResponse.Data.EsignID = userDataMap["esign_id"]
	}

	return validateTokenResponse, nil
}
