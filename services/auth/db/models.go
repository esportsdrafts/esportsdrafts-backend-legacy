package db

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

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

type Account struct {
	Base
	Username string `gorm:"varchar(100);not null;unique_index" json:"username"`
	Email    string `gorm:"varchar(100);not null;unique_index" json:"email"`
	Password string `gorm:"column:password_hash;varchar(256);not null" json:"password_hash"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
}
