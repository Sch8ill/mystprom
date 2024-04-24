package auth

import "time"

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshResponse struct {
	AccessToken    string        `json:"accessToken"`
	AccessTokenTTL time.Duration `json:"accessTokenTTLMs"`
}
