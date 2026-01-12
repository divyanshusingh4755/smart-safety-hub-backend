package user

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          string    `db:"id"`
	FullName    string    `db:"full_name"`
	Email       string    `db:"email"`
	Password    string    `db:"password"`
	PhoneNumber string    `db:"phone_number"`
	CompanyID   *string   `db:"company_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Roles struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Description *string   `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Permissions struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Description *string   `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type RolesPermissions struct {
	Role        string `db:"role"`
	Permissions string `db:"permissions"`
}

type UserRoles struct {
	UserId string `db:"user_id"`
	RoleId string `db:"role_id"`
}

type RefreshTokens struct {
	ID        string    `db:"id"`
	UserId    string    `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
	Revoked   bool      `db:"revoked"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u *User) HashPassword() error {
	if u.Password == "" {
		return errors.New("Password is empty")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Error hashing password: %v", err)
	}

	u.Password = string(hashed)
	return nil
}

func (u *User) ComparePassword(password string) error {
	if u.Password == "" {
		return errors.New("Password is empty")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return fmt.Errorf("Incorrect password")
	}
	return nil
}
