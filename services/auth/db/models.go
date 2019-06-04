
package db

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"update_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.NewV4()
	return scope.SetColumn("ID", uuid)
}

type Account struct {
	Base
	Username string `gorm:"column:name;size:128;not null;" json:"name"`
	Email string `gorm:"column:email;size:128;not null;" json:"email"`
	Password string `gorm:"column:name;size:128;not null;"`
}
