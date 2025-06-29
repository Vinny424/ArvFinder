package models

import (
	"time"
)

type User struct {
	ID           string    `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	TenantID     string    `json:"tenant_id" db:"tenant_id"`
	Role         string    `json:"role" db:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Tenant struct {
	ID             string    `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	SubscriptionTier string  `json:"subscription_tier" db:"subscription_tier"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}