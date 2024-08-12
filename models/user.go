package models

import (
	"time"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type SubscriptionLevel string

const (
	SubscriptionFree     SubscriptionLevel = "free"
	SubscriptionPremium  SubscriptionLevel = "premium"
	SubscriptionBusiness SubscriptionLevel = "business"
)

type User struct {
	ID                string            `json:"id"`
	Email             string            `json:"email"`
	Role              Role              `json:"role"`
	SubscriptionLevel SubscriptionLevel `json:"subscription_level"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}
