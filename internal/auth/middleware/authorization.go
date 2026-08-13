package middleware

import (
	"github.com/gofiber/fiber/v3"

	"github.com/panuwat39/my-api/internal/shared/fiberx"
)

func RequireRole(
	allowedRoles ...string,
) fiber.Handler {
	allowed := make(
		map[string]struct{},
		len(allowedRoles),
	)

	for _, role := range allowedRoles {
		allowed[role] = struct{}{}
	}

	return func(ctx fiber.Ctx) error {
		user, ok := CurrentUser(ctx)
		if !ok {
			return fiberx.Error(
				ctx,
				fiber.StatusUnauthorized,
				"UNAUTHORIZED",
				"authentication is required",
			)
		}

		if _, exists := allowed[user.Role]; !exists {
			return fiberx.Error(
				ctx,
				fiber.StatusForbidden,
				"FORBIDDEN",
				"you do not have permission to perform this action",
			)
		}

		return ctx.Next()
	}
}

func RequireSelfOrRole(
	paramName string,
	allowedRoles ...string,
) fiber.Handler {
	allowed := make(
		map[string]struct{},
		len(allowedRoles),
	)

	for _, role := range allowedRoles {
		allowed[role] = struct{}{}
	}

	return func(ctx fiber.Ctx) error {
		user, ok := CurrentUser(ctx)
		if !ok {
			return fiberx.Error(
				ctx,
				fiber.StatusUnauthorized,
				"UNAUTHORIZED",
				"authentication is required",
			)
		}

		targetUserID := fiber.Params[string](
			ctx,
			paramName,
		)

		if targetUserID != "" &&
			targetUserID == user.UserID {
			return ctx.Next()
		}

		if _, exists := allowed[user.Role]; exists {
			return ctx.Next()
		}

		return fiberx.Error(
			ctx,
			fiber.StatusForbidden,
			"FORBIDDEN",
			"you do not have permission to perform this action",
		)
	}
}
