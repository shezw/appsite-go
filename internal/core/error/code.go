// Copyright (c) 2026 shezw. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package error

// ErrorCode is a custom error code type
type ErrorCode int

const (
	// Success code
	Success ErrorCode = 200

	// Client Side Errors (4xx)
	InvalidParams ErrorCode = 400
	Unauthorized  ErrorCode = 401
	Forbidden     ErrorCode = 403
	NotFound      ErrorCode = 404
	MethodNotAllowed ErrorCode = 405
	Conflict		 ErrorCode = 409

	// Server Side Errors (5xx)
	ServerError        ErrorCode = 500
	ServiceUnavailable ErrorCode = 503

	// Business Logic Errors (10000+)
	// Access: 100xx
	TokenInvalid ErrorCode = 10001
	TokenExpired ErrorCode = 10002
	
	// User: 200xx
	UserNotFound  ErrorCode = 20001
	UserExists    ErrorCode = 20002
	InvalidPass   ErrorCode = 20003
)

// Int returns the integer value of the error code
func (e ErrorCode) Int() int {
	return int(e)
}

// String returns the generic message for the error code
func (e ErrorCode) String() string {
	switch e {
	case Success:
		return "Success"
	case InvalidParams:
		return "Invalid Parameters"
	case Unauthorized:
		return "Unauthorized"
	case Forbidden:
		return "Forbidden"
	case NotFound:
		return "Not Found"
	case ServerError:
		return "Internal Server Error"
	case TokenInvalid:
		return "Invalid Token"
	case TokenExpired:
		return "Token Expired"
	case UserNotFound:
		return "User Not Found"
	case UserExists:
		return "User Already Exists"
	default:
		return "Unknown Error"
	}
}
