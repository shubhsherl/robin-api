package models

import (
	"time"

	"github.com/RobinHoodArmyHQ/robin-api/pkg/nanoid"
)

type UserInfo struct {
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	PasswordHash string `json:"password_hash,omitempty"`
}

type UserVerification struct {
	ID             uint64        `json:"-" gorm:"primaryKey"`
	UserID         nanoid.NanoID `json:"user_id,omitempty"`
	EmailId        string        `json:"email_id"`
	Otp            uint64        `json:"otp,omitempty"`
	OtpGeneratedAt time.Time     `json:"otp_generated_at,omitempty"`
	OtpExpiresAt   time.Time     `json:"otp_expires_at,omitempty"`
	OtpRetryCount  uint64        `json:"otp_retry_count,omitempty"`
	IsVerified     int8          `json:"is_verified,omitempty"`
	UserInfo       UserInfo      `json:"user_info,omitempty" gorm:"serializer:json"`
}
