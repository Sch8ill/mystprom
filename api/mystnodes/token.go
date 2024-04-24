package mystnodes

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Token struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

func NewToken(token string, ttl time.Duration) *Token {
	return &Token{
		Token:   token,
		Expires: time.Now().Add(ttl),
	}
}

func NewTokenFromFile(path string) (*Token, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	token := new(Token)
	if err := json.NewDecoder(file).Decode(token); err != nil {
		return nil, fmt.Errorf("failed to parse token file: %w", err)
	}

	return token, nil
}

func (t *Token) Value() string {
	return t.Token
}

func (t *Token) Expired() bool {
	return time.Now().After(t.Expires)
}

func (t *Token) Save(path string) error {
	b, err := json.Marshal(t)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, b, 0644); err != nil {
		return err
	}
	return nil
}
