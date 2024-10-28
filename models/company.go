package models

import (
	"time"
)

type CompanyType string

const (
	CompanyTypeCorporation CompanyType = "corporation"
	CompanyTypeLLC         CompanyType = "llc"
	CompanyTypePartnership CompanyType = "partnership"
)

type Company struct {
	ID           string    `bson:"_id,omitempty" json:"id"`
	Name         string    `bson:"name" json:"name" binding:"required"`
	Address      string    `bson:"address" json:"address" binding:"required"`
	Phone        string    `bson:"phone" json:"phone" binding:"required"`
	Email        string    `bson:"email" json:"email" binding:"required,email"`
	Website      string    `bson:"website" json:"website"`
	TaxID        string    `bson:"tax_id" json:"tax_id" binding:"required"`
	BusinessType string    `bson:"business_type" json:"business_type" binding:"required"`
	CreatedBy    string    `bson:"created_by" json:"created_by"`
	UpdatedBy    string    `bson:"updated_by" json:"updated_by"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at" json:"updated_at"`
}

type BranchOffice struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	CompanyID string    `bson:"company_id" json:"company_id"`
	Name      string    `bson:"name" json:"name" binding:"required"`
	Address   string    `bson:"address" json:"address" binding:"required"`
	Phone     string    `bson:"phone" json:"phone" binding:"required"`
	Email     string    `bson:"email" json:"email" binding:"required,email"`
	CreatedBy string    `bson:"created_by" json:"created_by"`
	UpdatedBy string    `bson:"updated_by" json:"updated_by"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type CompanyInfo struct {
	ID                 string         `bson:"_id" json:"id"`
	Name               string         `bson:"name" json:"name"`
	Type               CompanyType    `bson:"type" json:"type"`
	RegistrationNumber string         `bson:"registration_number" json:"registration_number"`
	TaxID              string         `bson:"tax_id" json:"tax_id"`
	HeadOfficeAddress  string         `bson:"head_office_address" json:"head_office_address"`
	BranchOffices      []BranchOffice `bson:"branch_offices" json:"branch_offices"`
	Phone              string         `bson:"phone" json:"phone"`
	Email              string         `bson:"email" json:"email"`
	Website            string         `bson:"website" json:"website"`
	LogoURL            string         `bson:"logo_url" json:"logo_url"`
	Documents          []Document     `bson:"documents" json:"documents"`
	CreatedAt          time.Time      `bson:"created_at" json:"created_at"`
	UpdatedAt          time.Time      `bson:"updated_at" json:"updated_at"`
}

type Document struct {
	Type        string    `bson:"type" json:"type"`
	Number      string    `bson:"number" json:"number"`
	URL         string    `bson:"url" json:"url"`
	ExpiryDate  time.Time `bson:"expiry_date" json:"expiry_date"`
	IssuedDate  time.Time `bson:"issued_date" json:"issued_date"`
	Description string    `bson:"description" json:"description"`
}

type CompanyUpdate struct {
	Name               *string      `json:"name,omitempty"`
	TradingName        *string      `json:"trading_name,omitempty"`
	Type               *CompanyType `json:"type,omitempty"`
	RegistrationNumber *string      `json:"registration_number,omitempty"`
	TaxID              *string      `json:"tax_id,omitempty"`
	HeadOfficeAddress  *string      `json:"head_office_address,omitempty"`
	Phone              *string      `json:"phone,omitempty"`
	Email              *string      `json:"email,omitempty"`
	Website            *string      `json:"website,omitempty"`
	LogoURL            *string      `json:"logo_url,omitempty"`
}

type CreateCompanyRequest struct {
	Name         string `json:"name" binding:"required"`
	Address      string `json:"address" binding:"required"`
	Phone        string `json:"phone" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Website      string `json:"website"`
	TaxID        string `json:"tax_id" binding:"required"`
	BusinessType string `json:"business_type" binding:"required"`
}

type UpdateCompanyRequest struct {
	Name         string `json:"name"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
	Email        string `json:"email" binding:"omitempty,email"`
	Website      string `json:"website"`
	TaxID        string `json:"tax_id"`
	BusinessType string `json:"business_type"`
}
