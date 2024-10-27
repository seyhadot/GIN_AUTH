package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleUser      Role = "user"
	RoleSuperUser Role = "super"
)

type User struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Username  string    `bson:"username" json:"username"`
	Password  string    `bson:"password" json:"-"`
	Roles     []Role    `bson:"roles" json:"roles"`
	FullName  string    `bson:"full_name" json:"full_name"`
	Bio       string    `bson:"bio" json:"bio"`
	Avatar    string    `bson:"avatar" json:"avatar"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6"`
	Roles    []Role `json:"roles"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type UpdateProfileRequest struct {
	FullName string `json:"full_name" binding:"omitempty,max=100"`
	Bio      string `json:"bio" binding:"omitempty,max=500"`
	Avatar   string `json:"avatar" binding:"omitempty,url"`
}

// Add this new struct for password update
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required,min=6"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// SetDefaultRoles sets the default roles if none are provided
func (u *User) SetDefaultRoles() {
	if len(u.Roles) == 0 {
		u.Roles = []Role{RoleUser}
	}
}
