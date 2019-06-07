package db

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Base struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.NewV4()
	return scope.SetColumn("ID", uuid)
}

type Account struct {
	Base
	Username string `gorm:"varchar(100);not null;unique_index"`
	Email    string `gorm:"varchar(100);not null;unique_index"`
	Password string `gorm:"column:password_hash;varchar(256);not null"`
	EmailVerifiedAt *time.Time
}
