package port

import "context"

type AccessTokenClaims struct {
	UserID string
	Email  string
	Role   string
}

type TokenIssuer interface {
	IssueAccessToken(
		ctx context.Context,
		userID string,
		email string,
		role string,
	) (string, error)
}

type TokenVerifier interface {
	VerifyAccessToken(
		ctx context.Context,
		accessToken string,
	) (AccessTokenClaims, error)
}
