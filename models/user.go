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
	RoleStaff     Role = "staff" // New role for branch office staff
)

type User struct {
	ID            string    `bson:"_id,omitempty" json:"id"`
	Username      string    `bson:"username" json:"username"`
	Password      string    `bson:"password" json:"-"`
	Roles         []Role    `bson:"roles" json:"roles"`
	FullName      string    `bson:"full_name" json:"full_name"`
	Bio           string    `bson:"bio" json:"bio"`
	Avatar        string    `bson:"avatar" json:"avatar"`
	CompanyID     string    `bson:"company_id" json:"company_id"`         // Associated company
	BranchOffices []string  `bson:"branch_offices" json:"branch_offices"` // Associated branch offices
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Username      string   `json:"username" binding:"required,min=3,max=32"`
	Password      string   `json:"password" binding:"required,min=6"`
	FullName      string   `json:"full_name" binding:"required"`
	Roles         []Role   `json:"roles"`
	CompanyID     string   `json:"company_id"`
	BranchOffices []string `json:"branch_offices"`
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

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required,min=6"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

// Staff specific types
type StaffAssignment struct {
	UserID        string    `bson:"user_id" json:"user_id"`
	CompanyID     string    `bson:"company_id" json:"company_id"`
	BranchOffices []string  `bson:"branch_offices" json:"branch_offices"`
	AssignedAt    time.Time `bson:"assigned_at" json:"assigned_at"`
	AssignedBy    string    `bson:"assigned_by" json:"assigned_by"`
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

func (u *User) SetDefaultRoles() {
	if len(u.Roles) == 0 {
		u.Roles = []Role{RoleUser}
	}
}

// IsStaff checks if the user has staff role
func (u *User) IsStaff() bool {
	for _, role := range u.Roles {
		if role == RoleStaff {
			return true
		}
	}
	return false
}

// HasAccessToBranch checks if user has access to a specific branch office
func (u *User) HasAccessToBranch(branchID string) bool {
	if !u.IsStaff() {
		return false
	}
	for _, branch := range u.BranchOffices {
		if branch == branchID {
			return true
		}
	}
	return false
}
