package util

import (
	jsoniter "github.com/json-iterator/go"
)

var Json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	ContextTokenValueKey = "token-value"
	ContextJwtClaimKey   = "jwt-claim"
	ContextRouterKey     = "router-property"
	ApiKey               = "x-api-token"

	TagRouteDefault = "default"

	SettingValueTrue = "1"

	TypeSocialMedia int32 = 1
	TypeOnlineShop  int32 = 2

	ShowList                  int32 = 99999999
	PageSizeMicrositeProducts int32 = 15
	DefaultPage               int32 = 1
	DefaultCount              int32 = 15
)

type EmptyObject struct{}

type Response struct {
	Status     interface{}            `json:"status"`
	Code       interface{}            `json:"code"`
	HTTPStatus int                    `json:"-"`
	Message    string                 `json:"message"`
	Data       interface{}            `json:"data"`
	Errors     []string               `json:"errors,omitempty"`
	Header     map[string]interface{} `json:"-"`
}
