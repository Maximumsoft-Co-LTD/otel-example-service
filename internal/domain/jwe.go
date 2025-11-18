package domain

import (
	"github.com/go-jose/go-jose/v4/jwt"
)

type TokenPayload struct {
	Roles []string `json:"roles,omitempty"`
	jwt.Claims
}
