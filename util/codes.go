package util

import "google.golang.org/grpc/codes"

const (
	// Success status
	Success codes.Code = 200
	//SuccessCreated status
	SuccessCreated codes.Code = 201
	// SuccessNoContent status
	SuccessNoContent codes.Code = 204
	// InvalidArgument status
	InvalidArgument codes.Code = 400
	// Unauthorized status
	Unauthorized codes.Code = 401
	// Forbidden status
	Forbidden codes.Code = 403
	// NotFound status
	NotFound codes.Code = 404
	// Cancelled status
	Cancelled codes.Code = 405
	// RequestTimeout status
	RequestTimeout codes.Code = 408

	// InactiveAccount status
	InactiveAccount codes.Code = 410
	// InvalidToken status
	InvalidToken codes.Code = 411
	// InvalidAPIKey status
	InvalidAPIKey codes.Code = 412
	// InvalidSession status
	InvalidSession codes.Code = 413
	// ResourceExhausted status
	ResourceExhausted codes.Code = 414

	// InvalidSubdomain status
	InvalidSubdomain codes.Code = 420
	// InactiveSubdomain status
	InactiveSubdomain codes.Code = 421
	// SuspendedSubdomain status
	SuspendedSubdomain codes.Code = 422

	// InvalidTransaction status
	InvalidTransaction codes.Code = 430
	// DuplicateTransaction status
	DuplicateTransaction codes.Code = 431

	// InternalError status
	InternalError codes.Code = 500
	// ProcessingError status
	ProcessingError codes.Code = 502
)

// StatusMessage represent string message for code
var StatusMessage = map[codes.Code]string{
	Success:              "Success",
	SuccessCreated:       "Success, Data Saved",
	SuccessNoContent:     "Success, Data Not Found",
	InvalidArgument:      "Invalid Parameter",
	Unauthorized:         "Invalid Username or Password",
	Forbidden:            "Access Denied or Unauthorized",
	NotFound:             "Data Not Found",
	Cancelled:            "Request Cancelled",
	RequestTimeout:       "Request Timeout",
	InactiveAccount:      "Inactive Account",
	InvalidToken:         "Invalid Access Token or Expired Login",
	InvalidAPIKey:        "Invalid API Key",
	InvalidSession:       "Invalid or Expired Session",
	ResourceExhausted:    "Resource Limit Exceeded",
	InvalidSubdomain:     "Invalid Store Name",
	InactiveSubdomain:    "Store Not Activated",
	SuspendedSubdomain:   "Blocked, please contact pan via email pancanugrahawicaksana@gmail.com",
	InvalidTransaction:   "Invalid or Mismatched Transaction Data",
	DuplicateTransaction: "Duplicate Transaction Data",
	InternalError:        "Server Error",
}
