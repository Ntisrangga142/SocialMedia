package models

import "time"

type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type BlacklistToken struct {
	Token     string        `json:"token"`
	ExpiresAt time.Duration `json:"expires_in"`
}
