package auth

import "time"

type LoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`
}

type LoginResponse struct {
	UserID          string        `json:"userId"`
	IsAdmin         bool          `json:"isAdmin"`
	AccessToken     string        `json:"accessToken"`
	RefreshToken    string        `json:"refreshToken"`
	IsFirstLogin    bool          `json:"isFirstLogin"`
	AccessTokenTTL  time.Duration `json:"accessTokenTTLMs"`
	RefreshTokenTTL time.Duration `json:"refreshTokenTTLMs"`
}
