package response

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Data any `json:"data"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

func JSON(w http.ResponseWriter, statusCode int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(payload)
}

func Success(w http.ResponseWriter, statusCode int, data any) error {
	return JSON(
		w,
		statusCode,
		SuccessResponse{
			Data: data,
		},
	)
}

func Error(
	w http.ResponseWriter,
	statusCode int,
	code string,
	message string,
) error {
	return JSON(
		w,
		statusCode,
		ErrorResponse{
			Error: ErrorDetail{
				Code:    code,
				Message: message,
			},
		},
	)
}
