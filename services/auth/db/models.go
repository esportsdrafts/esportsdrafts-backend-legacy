package db

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Base fixes the default behavior in GORM to use UUID as primary key instead of
// numbered ID:s. Grab the defaults but change the primary key to UUID.
type Base struct {
	ID        uuid.UUID  `gorm:"varchar(36);primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"update_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.NewV4()
	return scope.SetColumn("ID", uuid.String())
}

// Account is a definition of the pure login info of a user. User details are
// stored in the `user` service.
type Account struct {
	Base
	Username               string     `gorm:"varchar(100);not null;unique_index" json:"username"`
	Email                  string     `gorm:"varchar(100);not null;unique_index" json:"email"`
	Password               string     `gorm:"column:password_hash;varchar(256);not null" json:"-"`
	AcceptedTermsAt        *time.Time `json:"accepted_terms_at"`
	EmailVerifiedAt        *time.Time `json:"email_verified_at"`
	EmailVerificationCodes []EmailVerificationCode
}

// EmailVerificationCode is used to verify a users email
// Just use the built-in ID of the object as the verify code
type EmailVerificationCode struct {
	Base
	UserID    string    `gorm:"varchar(36);not null;"`
	ExpiresAt time.Time `gorm:"not null;" json:"expires_at"`
}
