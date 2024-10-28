package controllers

import (
	"loan/config"
	"loan/models"
	"loan/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CompanyController struct {
	config *config.Config
}

func NewCompanyController(config *config.Config) *CompanyController {
	return &CompanyController{config: config}
}

// CreateCompany creates a new company
func (cc *CompanyController) CreateCompany(c *gin.Context) {
	var req models.CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	userID := c.GetString("user_id")
	company := models.Company{
		ID:           primitive.NewObjectID().Hex(),
		Name:         req.Name,
		Address:      req.Address,
		Phone:        req.Phone,
		Email:        req.Email,
		Website:      req.Website,
		TaxID:        req.TaxID,
		BusinessType: req.BusinessType,
		CreatedBy:    userID,
		UpdatedBy:    userID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err := cc.config.MongoDB.Collection("companies").InsertOne(c, company)
	if err != nil {
		utils.InternalError(c, "Error creating company")
		return
	}

	c.JSON(http.StatusCreated, company)
}

// GetCompany gets a company by ID
func (cc *CompanyController) GetCompany(c *gin.Context) {
	id := c.Param("id")

	var company models.Company
	err := cc.config.MongoDB.Collection("companies").FindOne(c, bson.M{"_id": id}).Decode(&company)
	if err != nil {
		utils.BadRequest(c, "Company not found")
		return
	}

	c.JSON(http.StatusOK, company)
}

// UpdateCompany updates a company
func (cc *CompanyController) UpdateCompany(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
			"updated_by": c.GetString("user_id"),
		},
	}

	// Only update fields that are provided
	if req.Name != "" {
		update["$set"].(bson.M)["name"] = req.Name
	}
	if req.Address != "" {
		update["$set"].(bson.M)["address"] = req.Address
	}
	if req.Phone != "" {
		update["$set"].(bson.M)["phone"] = req.Phone
	}
	if req.Email != "" {
		update["$set"].(bson.M)["email"] = req.Email
	}
	if req.Website != "" {
		update["$set"].(bson.M)["website"] = req.Website
	}
	if req.TaxID != "" {
		update["$set"].(bson.M)["tax_id"] = req.TaxID
	}
	if req.BusinessType != "" {
		update["$set"].(bson.M)["business_type"] = req.BusinessType
	}

	result := cc.config.MongoDB.Collection("companies").FindOneAndUpdate(
		c,
		bson.M{"_id": id},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var updatedCompany models.Company
	if err := result.Decode(&updatedCompany); err != nil {
		utils.BadRequest(c, "Company not found")
		return
	}

	c.JSON(http.StatusOK, updatedCompany)
}

// DeleteCompany deletes a company
func (cc *CompanyController) DeleteCompany(c *gin.Context) {
	id := c.Param("id")

	result, err := cc.config.MongoDB.Collection("companies").DeleteOne(c, bson.M{"_id": id})
	if err != nil {
		utils.InternalError(c, "Error deleting company")
		return
	}

	if result.DeletedCount == 0 {
		utils.BadRequest(c, "Company not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Company deleted successfully"})
}

// ListCompanies lists all companies with pagination
func (cc *CompanyController) ListCompanies(c *gin.Context) {
	page, limit := utils.GetPaginationParams(c)
	skip := (page - 1) * limit

	// Get total count
	total, err := cc.config.MongoDB.Collection("companies").CountDocuments(c, bson.M{})
	if err != nil {
		utils.InternalError(c, "Error counting companies")
		return
	}

	// Set options for pagination and sorting
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	// Fetch companies
	cursor, err := cc.config.MongoDB.Collection("companies").Find(c, bson.M{}, opts)
	if err != nil {
		utils.InternalError(c, "Error fetching companies")
		return
	}
	defer cursor.Close(c)

	var companies []models.Company
	if err = cursor.All(c, &companies); err != nil {
		utils.InternalError(c, "Error parsing companies")
		return
	}

	utils.SendPaginatedResponse(c, companies, total, page, limit)
}
