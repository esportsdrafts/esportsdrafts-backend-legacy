package db

import (
	"time"

	"database/sql/driver"

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
	Username        string     `gorm:"varchar(128);not null;unique_index" json:"username"`
	Email           string     `gorm:"varchar(256);not null;unique_index" json:"email"`
	Password        string     `gorm:"column:password_hash;varchar(256);not null" json:"-"`
	AcceptedTermsAt *time.Time `json:"accepted_terms_at"`
	MFA             *MFAMethod `json:"mfa_method"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
}

// EmailVerificationCode is used to verify a users email
// Uses the built-in ID of the object as the verify code
type EmailVerificationCode struct {
	Base
	User      Account   `gorm:"foreignkey:UserID"`
	UserID    uuid.UUID `gorm:"varchar(36);not null;index;" json:"user_id"`
	ExpiresAt time.Time `gorm:"not null;" json:"expires_at"`
}

// PasswordResetToken is used to reset a user's password
// Uses the built-in ID of the object as the verify code
type PasswordResetToken struct {
	Base
	User      Account   `gorm:"foreignkey:UserID"`
	UserID    uuid.UUID `gorm:"varchar(36);not null;index;" json:"user_id"`
	ExpiresAt time.Time `gorm:"not null;" json:"expires_at"`
}

// MFACode is very similar to email verification codes. But we have explicit
// code property since it needs to be a bit more human-readable compared to
// UUID:s.
type MFACode struct {
	Base
	User      Account   `gorm:"foreignkey:UserID"`
	UserID    uuid.UUID `gorm:"varchar(36);not null;index;" json:"user_id"`
	Code      string    `gorm:"varchar(10);not null;unique_index" json:"code"`
	ExpiresAt time.Time `gorm:"not null;" json:"expires_at"`
}

// MFAMethodOption denotes a type of MFA
type MFAMethodOption string

// NullAccount denotes a dummy user used to do perform password compare
var /* const */ NullAccount = Account{
	Password: "$argon2id$v=19$m=65536,t=3,p=2$DkjlxfkTZlrgGY8XyfFLjw$w+sqKlCc8ju5ZacLtmyAngicTDMZUAcAGIHzqQa9Wi0",
}

const (
	emailMFA MFAMethodOption = "email"
)

// Scan turns a value into a MFAMethodOption
func (p *MFAMethodOption) Scan(value interface{}) error {
	*p = MFAMethodOption(value.([]byte))
	return nil
}

// Value gives the value as string of a MFAMethodOption
func (p MFAMethodOption) Value() (driver.Value, error) {
	return string(p), nil
}

// MFAMethod denotes a MFA device type.
type MFAMethod struct {
	Base
	Type string `gorm:"type:ENUM('email');not null;" json:"type"`
}

// SetMFAMethod sets MFA method on an account
func (a *Account) SetMFAMethod(db *gorm.DB, method MFAMethodOption) error {
	return nil
}

// VerifyEmail marks and account as verified
func (a *Account) VerifyEmail(db *gorm.DB) error {
	timeNow := time.Now()
	a.EmailVerifiedAt = &timeNow
	err := db.Save(a).Error
	if err != nil {
		return err
	}
	// Ignore all errors here since deleting is not really important
	db.Where("user_id = ?", a.ID).Delete(EmailVerificationCode{})
	return nil
}

// IsEmailVerified returns true if an email has been verified for the account
func (a *Account) IsEmailVerified() bool {
	return a.EmailVerifiedAt != nil
}
