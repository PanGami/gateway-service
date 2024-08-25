package middleware

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/pangami/gateway-service/util"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type JWTClaims struct {
	Aud         string   `json:"aud"`
	Exp         float64  `json:"exp"`
	Iat         float64  `json:"iat"`
	Jti         string   `json:"jti"`
	Nbf         float64  `json:"nbf"`
	Permissions []string `json:"permissions"`
	Role        string   `json:"role"`
	Scopes      []string `json:"scopes"`
	Sub         string   `json:"sub"`
	Uesign      *string  `json:"uesign"`
	Uid         string   `json:"uid"`
	Uname       string   `json:"uname"`
	Utype       *string  `json:"utype"`
}

// CheckPermission checks if the user is authenticated and decodes the JWT token
func CheckPermission(requestString string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			apiToken := c.Request().Header.Get("x-api-token")
			if apiToken != "" {
				c.Request().Header.Set("Authorization", "Bearer "+apiToken)
			}

			jwtClaims, err := DecodeJWT(apiToken)
			if err != nil {
				log.Printf("JWT decoding error: %v", err)
				response := ErrorResponseWithTraceID("", "006", err.Error())
				return c.JSON(http.StatusUnauthorized, response)
			}

			if !Contains(jwtClaims.Permissions, requestString) {
				log.Printf("JWT decoding error: %v", err)
				response := ErrorResponseWithTraceID("", "007", "Akses ditolak")
				return c.JSON(http.StatusUnauthorized, response)
			}

			c.Set(util.ContextJwtClaimKey, jwtClaims)
			return next(c)
		}
	}
}

func DecodeJWT(tokenString string) (*JWTClaims, error) {
	// Check if the token is empty
	if tokenString == "" {
		log.Println("Token is empty")
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Token is empty")
	}

	// Split the token into parts
	tokenParts := strings.Split(tokenString, ".")
	if len(tokenParts) != 3 {
		log.Println("Invalid token format")
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid token format")
	}

	// Decode the token data
	tokenData, err := base64.RawStdEncoding.DecodeString(tokenParts[1])
	if err != nil {
		log.Println("Invalid token encoding:", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid token encoding")
	}

	// Unmarshal token claims
	var claims jwt.MapClaims
	if err := json.Unmarshal(tokenData, &claims); err != nil {
		log.Println("Invalid token claims:", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid token claims")
	}

	// Extract user data from claims
	jwtClaims := &JWTClaims{
		Uid:   claims["uid"].(string),
		Uname: claims["uname"].(string),
		Exp:   claims["exp"].(float64),
	}

	// Handle the Permissions field
	permissionsInterface, ok := claims["permissions"].([]interface{})
	if ok {
		permissions := make([]string, len(permissionsInterface))
		for i, v := range permissionsInterface {
			permissions[i] = v.(string)
		}
		jwtClaims.Permissions = permissions
	} else {
		log.Println("Permissions field missing or of incorrect type")
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Permissions field missing or of incorrect type")
	}

	return jwtClaims, nil
}

func Contains[T comparable](arr []T, target T) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}
