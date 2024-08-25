package util

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
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

type AppError interface {
	Status() interface{}
	Code() interface{}
	HTTPStatus() int
	Message() string
	Data() *DataItem
	Errors() []string
	Header() map[string]interface{}
}

type GenericException struct {
	ErrorCode    string                 `json:"code"`
	ErrorMessage string                 `json:"message"`
	ErrorHTTP    int                    `json:"-"`
	ErrorData    *DataItem              `json:"data"`
	ErrorErrors  []string               `json:"errors,omitempty"`
	ErrorHeader  map[string]interface{} `json:"-"`
}

type Pagination struct {
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
	Total       int `json:"total"`
	TotalPages  int `json:"total_pages"`
}
type DataItem struct {
	Items      interface{}  `json:"items"`
	Pagination []Pagination `json:"pagination"`
}

func NewGenericException(code, message string, httpStatus int) *GenericException {
	return &GenericException{
		ErrorCode:    code,
		ErrorMessage: message,
		ErrorHTTP:    httpStatus,
	}
}

// Implement methods required by the AppError interface
func (ge *GenericException) Status() interface{} {
	return false
}

func (ge *GenericException) Code() interface{} {
	return ge.ErrorCode
}

func (ge *GenericException) HTTPStatus() int {
	return ge.ErrorHTTP
}

func (ge *GenericException) Message() string {
	return ge.ErrorMessage
}

func (ge *GenericException) Data() *DataItem {
	return ge.ErrorData
}

func (ge *GenericException) Errors() []string {
	return ge.ErrorErrors
}

func (ge *GenericException) Header() map[string]interface{} {
	return ge.ErrorHeader
}

// CustomHTTPErrorHandler handles various types of errors and renders the JSON response
func CustomHTTPErrorHandler(err error, c echo.Context) {
	var genericException AppError

	// Type switch to handle different error types
	switch e := err.(type) {
	case *echo.HTTPError:
		code := e.Code
		message := e.Message.(string)
		genericException = mapHTTPErrorToGenericException(code, message)
	case AppError:
		genericException = e
	default:
		genericException = NewGenericException("999", "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
	}

	// Convert genericException to Response struct
	response := &Response{
		Status:  genericException.Status(),
		Code:    genericException.Code(),
		Message: genericException.Message(),
		Data:    genericException.Data(),
	}

	// Marshal response to JSON and send it
	if !c.Response().Committed {
		c.JSON(response.HTTPStatus, response)
	}
}

func mapHTTPErrorToGenericException(code int, message string) AppError {
	switch code {
	case http.StatusBadRequest:
		return NewGenericException("005", message, http.StatusForbidden)
	case http.StatusUnauthorized:
		return NewGenericException("006", message, http.StatusUnauthorized)
	case http.StatusForbidden:
		return NewGenericException("007", message, http.StatusForbidden)
	case http.StatusNotFound:
		return NewGenericException("008", message, http.StatusNotFound)
	case http.StatusMethodNotAllowed:
		return NewGenericException("009", message, http.StatusMethodNotAllowed)
	case http.StatusRequestTimeout:
		return NewGenericException("011", message, http.StatusRequestTimeout)
	case http.StatusConflict:
		return NewGenericException("012", message, http.StatusConflict)
	case http.StatusRequestEntityTooLarge:
		return NewGenericException("013", message, http.StatusRequestEntityTooLarge)
	case http.StatusRequestURITooLong:
		return NewGenericException("014", message, http.StatusRequestURITooLong)
	case http.StatusUnsupportedMediaType:
		return NewGenericException("015", message, http.StatusUnsupportedMediaType)
	case http.StatusTooManyRequests:
		return NewGenericException("016", message, http.StatusTooManyRequests)
	case http.StatusRequestHeaderFieldsTooLarge:
		return NewGenericException("017", message, http.StatusRequestHeaderFieldsTooLarge)
	default:
		return NewGenericException("999", "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
	}
}

func MustAtoi64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func MustAtoi32(s string) int32 {
	i, _ := strconv.ParseInt(s, 10, 32)
	return int32(i)
}

func MustAtof64(s string) float64 {
	i, _ := strconv.ParseFloat(s, 64)
	return i
}

func IntegerToString(i int64) string {
	s := strconv.Itoa(int(i))
	return s
}

// string to array string
func Explode(s string, separator string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, separator)
}

// covert array string to array int64
func ExplodeInt64(s string, separator string) []int64 {
	var integers []int64
	for _, v := range Explode(s, separator) {
		val, _ := strconv.Atoi(v)
		integers = append(integers, int64(val))
	}
	return integers
}

// covert array productsid to array int64
// products = [1,2,3,4,5]
func ExplodeProductsArray(s string, separator string) []int64 {
	var integers []int64
	s = strings.ReplaceAll(s, "[", "")
	s = strings.ReplaceAll(s, "]", "")
	for _, v := range Explode(s, separator) {
		val, _ := strconv.Atoi(v)
		integers = append(integers, int64(val))
	}
	return integers
}

func ReplaceTimeZone(s string) string {
	s = strings.ReplaceAll(s, "T", " ")
	s = strings.ReplaceAll(s, "Z", "")
	return s
}

func StringToInteger(txt string) int {
	i, _ := strconv.Atoi(txt)
	return int(i)
}

func ArrayQueryParams(s, sp string) []string {
	var res []string
	if len(s) < 1 {
		return res
	}
	return strings.Split(s, sp)
}

func StringToBool(s string) bool {
	resp, _ := strconv.ParseBool(s)
	return resp
}

func CheckDefaultPage(s string) int32 {
	page := MustAtoi32(s)
	if page == 0 {
		return DefaultPage
	} else {
		return page
	}
}

// FindInArray is
func FindInArray(one []string, two string) bool {
	for _, val := range one {
		if two == val {
			return true
		}
	}
	return false
}

func BoolToString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

// FormatHourMinute for format hour & minute
// from : "09.00"
// to : "09:00"
func FormatHourMinute(req string) (response string) {
	var hourTmp string
	var minuteTmp string
	if req != "" {
		hourTmp = req[0:2]
		minuteTmp = req[3:5]

		response = hourTmp + ":" + minuteTmp
	}

	return
}

func Slugger(str string) (response string) {
	response = strings.ReplaceAll(strings.ToLower(str), " ", "-")

	return
}

func BindAndValidate(i interface{}, c echo.Context) error {
	if err := c.Bind(i); err != nil {
		return err
	}

	return c.Validate(i)
}

func GenerateID(prefix string) string {
	currentTime := time.Now()
	formattedTime := currentTime.Format("20060102150405.000000")
	formattedTime = formattedTime[:14] + formattedTime[15:]

	return fmt.Sprintf("%s%s", prefix, formattedTime)
}

func TimeToString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
