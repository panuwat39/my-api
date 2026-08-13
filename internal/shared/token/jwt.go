package token

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	authport "github.com/panuwat39/my-api/internal/auth/port"
)

var (
	ErrJWTSecretRequired  = errors.New("JWT secret is required")
	ErrInvalidAccessToken = errors.New("invalid access token")
)

type JWTIssuer struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`

	jwt.RegisteredClaims
}

func NewJWTIssuer(
	secret string,
	issuer string,
	ttl time.Duration,
) (*JWTIssuer, error) {
	if strings.TrimSpace(secret) == "" {
		return nil, ErrJWTSecretRequired
	}

	if len([]byte(secret)) < 32 {
		return nil, errors.New(
			"JWT secret must be at least 32 bytes",
		)
	}

	if strings.TrimSpace(issuer) == "" {
		return nil, errors.New(
			"JWT issuer is required",
		)
	}

	if ttl <= 0 {
		return nil, errors.New(
			"JWT access TTL must be greater than zero",
		)
	}

	return &JWTIssuer{
		secret: []byte(secret),
		issuer: issuer,
		ttl:    ttl,
	}, nil
}

func (i *JWTIssuer) IssueAccessToken(
	ctx context.Context,
	userID string,
	email string,
	role string,
) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	now := time.Now().UTC()

	claims := Claims{
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    i.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(
				now.Add(i.ttl),
			),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, err := token.SignedString(
		i.secret,
	)
	if err != nil {
		return "", fmt.Errorf(
			"sign access token: %w",
			err,
		)
	}

	return signedToken, nil
}

func (i *JWTIssuer) VerifyAccessToken(
	ctx context.Context,
	accessToken string,
) (authport.AccessTokenClaims, error) {
	if err := ctx.Err(); err != nil {
		return authport.AccessTokenClaims{}, err
	}

	accessToken = strings.TrimSpace(accessToken)
	if accessToken == "" {
		return authport.AccessTokenClaims{},
			ErrInvalidAccessToken
	}

	claims := &Claims{}

	parsedToken, err := jwt.ParseWithClaims(
		accessToken,
		claims,
		func(token *jwt.Token) (any, error) {
			return i.secret, nil
		},
		jwt.WithValidMethods(
			[]string{
				jwt.SigningMethodHS256.Alg(),
			},
		),
		jwt.WithIssuer(i.issuer),
		jwt.WithExpirationRequired(),
		jwt.WithNotBeforeRequired(),
	)
	if err != nil {
		return authport.AccessTokenClaims{},
			ErrInvalidAccessToken
	}

	if !parsedToken.Valid {
		return authport.AccessTokenClaims{},
			ErrInvalidAccessToken
	}

	if strings.TrimSpace(claims.Subject) == "" {
		return authport.AccessTokenClaims{},
			ErrInvalidAccessToken
	}

	if strings.TrimSpace(claims.Email) == "" {
		return authport.AccessTokenClaims{},
			ErrInvalidAccessToken
	}

	if strings.TrimSpace(claims.Role) == "" {
		return authport.AccessTokenClaims{},
			ErrInvalidAccessToken
	}

	return authport.AccessTokenClaims{
		UserID: claims.Subject,
		Email:  claims.Email,
		Role:   claims.Role,
	}, nil
}

var _ authport.TokenIssuer = (*JWTIssuer)(nil)
var _ authport.TokenVerifier = (*JWTIssuer)(nil)
