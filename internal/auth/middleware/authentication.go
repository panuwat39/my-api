package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	authport "github.com/panuwat39/my-api/internal/auth/port"
	"github.com/panuwat39/my-api/internal/shared/fiberx"
)

type authenticatedUserKeyType struct{}

var authenticatedUserKey authenticatedUserKeyType

type AuthenticatedUser struct {
	UserID string
	Email  string
	Role   string
}

func Authenticate(
	tokenVerifier authport.TokenVerifier,
) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		authorization := strings.TrimSpace(
			ctx.Get("Authorization"),
		)

		parts := strings.Fields(authorization)

		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") {
			return fiberx.Error(
				ctx,
				fiber.StatusUnauthorized,
				"UNAUTHORIZED",
				"authentication is required",
			)
		}

		claims, err := tokenVerifier.VerifyAccessToken(
			ctx.Context(),
			parts[1],
		)
		if err != nil {
			return fiberx.Error(
				ctx,
				fiber.StatusUnauthorized,
				"INVALID_ACCESS_TOKEN",
				"access token is invalid or expired",
			)
		}

		ctx.Locals(
			authenticatedUserKey,
			AuthenticatedUser{
				UserID: claims.UserID,
				Email:  claims.Email,
				Role:   claims.Role,
			},
		)

		return ctx.Next()
	}
}

func CurrentUser(
	ctx fiber.Ctx,
) (AuthenticatedUser, bool) {
	user, ok := ctx.Locals(
		authenticatedUserKey,
	).(AuthenticatedUser)

	return user, ok
}
