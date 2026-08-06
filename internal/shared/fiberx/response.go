package fiberx

import "github.com/gofiber/fiber/v3"

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

func Success(
	c fiber.Ctx,
	statusCode int,
	data any,
) error {
	return c.Status(statusCode).JSON(
		SuccessResponse{
			Data: data,
		},
	)
}

func Error(
	c fiber.Ctx,
	statusCode int,
	code string,
	message string,
) error {
	return c.Status(statusCode).JSON(
		ErrorResponse{
			Error: ErrorDetail{
				Code:    code,
				Message: message,
			},
		},
	)
}
