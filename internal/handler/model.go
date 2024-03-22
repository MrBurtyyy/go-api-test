package handler

import "net/http"

var (
	httpServerError            = ErrorResponse{StatusCode: http.StatusInternalServerError, ErrorCode: "SERVER_ERROR", ErrorMessage: "An unexpected jsonError has occurred"}
	httpPayloadValidationError = ErrorResponse{StatusCode: http.StatusBadRequest, ErrorCode: "PAYLOAD_VALIDATION", ErrorMessage: "There are errors with the provided payload"}
	httpMalformedRequest       = ErrorResponse{StatusCode: http.StatusBadRequest, ErrorCode: "MALFORMED_REQUEST", ErrorMessage: "There were errors parsing the request"}
)

type PropertyValidationErrorResponse struct {
	Property string `json:"property"`
	Message  string `json:"message"`
}

type ErrorResponse struct {
	StatusCode       int                               `json:"-"`
	ErrorCode        string                            `json:"code"`
	ErrorMessage     string                            `json:"message"`
	ValidationErrors []PropertyValidationErrorResponse `json:"errors,omitempty"`
}

type ProductRequest struct {
	Name  string  `json:"name" validate:"required"`
	Price float64 `json:"price" validate:"required,dp=2"`
}

type ProductResponse struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
