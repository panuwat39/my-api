package model

type LoginResponse struct {
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}
