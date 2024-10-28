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

type BranchOfficeController struct {
	config *config.Config
}

func NewBranchOfficeController(config *config.Config) *BranchOfficeController {
	return &BranchOfficeController{config: config}
}

// CreateBranchOffice creates a new branch office for a company
func (bc *BranchOfficeController) CreateBranchOffice(c *gin.Context) {
	companyID := c.Param("id")

	// Verify company exists
	var company models.Company
	err := bc.config.MongoDB.Collection("companies").FindOne(c, bson.M{"_id": companyID}).Decode(&company)
	if err != nil {
		utils.BadRequest(c, "Company not found")
		return
	}

	var req models.BranchOffice
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	branchOffice := models.BranchOffice{
		ID:        primitive.NewObjectID().Hex(),
		CompanyID: companyID,
		Name:      req.Name,
		Address:   req.Address,
		Phone:     req.Phone,
		Email:     req.Email,
		CreatedBy: c.GetString("user_id"),
		UpdatedBy: c.GetString("user_id"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = bc.config.MongoDB.Collection("branch_offices").InsertOne(c, branchOffice)
	if err != nil {
		utils.InternalError(c, "Error creating branch office")
		return
	}

	c.JSON(http.StatusCreated, branchOffice)
}

// GetBranchOffice gets a branch office by ID
func (bc *BranchOfficeController) GetBranchOffice(c *gin.Context) {
	branchID := c.Param("branch_id")
	companyID := c.Param("id")

	var branchOffice models.BranchOffice
	err := bc.config.MongoDB.Collection("branch_offices").FindOne(c,
		bson.M{
			"_id":        branchID,
			"company_id": companyID,
		}).Decode(&branchOffice)

	if err != nil {
		utils.BadRequest(c, "Branch office not found")
		return
	}

	c.JSON(http.StatusOK, branchOffice)
}

// UpdateBranchOffice updates a branch office
func (bc *BranchOfficeController) UpdateBranchOffice(c *gin.Context) {
	branchID := c.Param("branch_id")
	companyID := c.Param("id")

	var req models.BranchOffice
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

	result := bc.config.MongoDB.Collection("branch_offices").FindOneAndUpdate(
		c,
		bson.M{
			"_id":        branchID,
			"company_id": companyID,
		},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var updatedBranchOffice models.BranchOffice
	if err := result.Decode(&updatedBranchOffice); err != nil {
		utils.BadRequest(c, "Branch office not found")
		return
	}

	c.JSON(http.StatusOK, updatedBranchOffice)
}

// ListBranchOffices lists all branch offices for a company
func (bc *BranchOfficeController) ListBranchOffices(c *gin.Context) {
	companyID := c.Param("id")
	page, limit := utils.GetPaginationParams(c)
	skip := (page - 1) * limit

	// Get total count
	total, err := bc.config.MongoDB.Collection("branch_offices").CountDocuments(c,
		bson.M{"company_id": companyID})
	if err != nil {
		utils.InternalError(c, "Error counting branch offices")
		return
	}

	// Set options for pagination and sorting
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	// Fetch branch offices
	cursor, err := bc.config.MongoDB.Collection("branch_offices").Find(c,
		bson.M{"company_id": companyID}, opts)
	if err != nil {
		utils.InternalError(c, "Error fetching branch offices")
		return
	}
	defer cursor.Close(c)

	var branchOffices []models.BranchOffice
	if err = cursor.All(c, &branchOffices); err != nil {
		utils.InternalError(c, "Error parsing branch offices")
		return
	}

	utils.SendPaginatedResponse(c, branchOffices, total, page, limit)
}

// DeleteBranchOffice deletes a branch office
func (bc *BranchOfficeController) DeleteBranchOffice(c *gin.Context) {
	branchID := c.Param("branch_id")
	companyID := c.Param("id")

	result, err := bc.config.MongoDB.Collection("branch_offices").DeleteOne(c,
		bson.M{
			"_id":        branchID,
			"company_id": companyID,
		})
	if err != nil {
		utils.InternalError(c, "Error deleting branch office")
		return
	}

	if result.DeletedCount == 0 {
		utils.BadRequest(c, "Branch office not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Branch office deleted successfully"})
}
