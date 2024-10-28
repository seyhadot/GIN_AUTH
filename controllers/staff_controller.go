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

type StaffController struct {
	config *config.Config
}

func NewStaffController(config *config.Config) *StaffController {
	return &StaffController{config: config}
}

// AssignStaffToBranch assigns a staff member to branch offices
func (sc *StaffController) AssignStaffToBranch(c *gin.Context) {
	var assignment models.StaffAssignment
	if err := c.ShouldBindJSON(&assignment); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	// Verify user exists and is a staff member
	var user models.User
	err := sc.config.MongoDB.Collection("users").FindOne(c, bson.M{"_id": assignment.UserID}).Decode(&user)
	if err != nil {
		utils.BadRequest(c, "User not found")
		return
	}

	if !user.IsStaff() {
		utils.BadRequest(c, "User is not a staff member")
		return
	}

	// Verify company exists
	var company models.Company
	err = sc.config.MongoDB.Collection("companies").FindOne(c, bson.M{"_id": assignment.CompanyID}).Decode(&company)
	if err != nil {
		utils.BadRequest(c, "Company not found")
		return
	}

	// Verify all branch offices exist and belong to the company
	for _, branchID := range assignment.BranchOffices {
		var branch models.BranchOffice
		err = sc.config.MongoDB.Collection("branch_offices").FindOne(c,
			bson.M{
				"_id":        branchID,
				"company_id": assignment.CompanyID,
			}).Decode(&branch)
		if err != nil {
			utils.BadRequest(c, "Invalid branch office ID: "+branchID)
			return
		}
	}

	// Update user's company and branch office assignments
	update := bson.M{
		"$set": bson.M{
			"company_id":     assignment.CompanyID,
			"branch_offices": assignment.BranchOffices,
			"updated_at":     time.Now(),
		},
	}

	result := sc.config.MongoDB.Collection("users").FindOneAndUpdate(
		c,
		bson.M{"_id": assignment.UserID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var updatedUser models.User
	if err := result.Decode(&updatedUser); err != nil {
		utils.InternalError(c, "Error updating user")
		return
	}

	// Record the assignment
	assignment.AssignedAt = time.Now()
	assignment.AssignedBy = c.GetString("user_id")

	_, err = sc.config.MongoDB.Collection("staff_assignments").InsertOne(c, assignment)
	if err != nil {
		utils.InternalError(c, "Error recording assignment")
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// ListStaffByBranch lists all staff members assigned to a branch office
func (sc *StaffController) ListStaffByBranch(c *gin.Context) {
	branchID := c.Param("branch_id")
	companyID := c.Param("id")

	page, limit := utils.GetPaginationParams(c)
	skip := (page - 1) * limit

	// Verify branch exists
	var branch models.BranchOffice
	err := sc.config.MongoDB.Collection("branch_offices").FindOne(c,
		bson.M{
			"_id":        branchID,
			"company_id": companyID,
		}).Decode(&branch)
	if err != nil {
		utils.BadRequest(c, "Branch office not found")
		return
	}

	// Find users assigned to this branch
	filter := bson.M{
		"roles":          models.RoleStaff,
		"company_id":     companyID,
		"branch_offices": branchID,
	}

	// Get total count
	total, err := sc.config.MongoDB.Collection("users").CountDocuments(c, filter)
	if err != nil {
		utils.InternalError(c, "Error counting staff members")
		return
	}

	// Set options for pagination and sorting
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := sc.config.MongoDB.Collection("users").Find(c, filter, opts)
	if err != nil {
		utils.InternalError(c, "Error fetching staff members")
		return
	}
	defer cursor.Close(c)

	var staff []models.User
	if err = cursor.All(c, &staff); err != nil {
		utils.InternalError(c, "Error parsing staff members")
		return
	}

	utils.SendPaginatedResponse(c, staff, total, page, limit)
}

// RemoveStaffFromBranch removes a staff member from a branch office
func (sc *StaffController) RemoveStaffFromBranch(c *gin.Context) {
	userID := c.Param("user_id")
	branchID := c.Param("branch_id")
	companyID := c.Param("id")

	// Verify user exists and is assigned to the branch
	var user models.User
	err := sc.config.MongoDB.Collection("users").FindOne(c,
		bson.M{
			"_id":            userID,
			"company_id":     companyID,
			"branch_offices": branchID,
		}).Decode(&user)
	if err != nil {
		utils.BadRequest(c, "Staff member not found")
		return
	}

	// Remove branch from user's assignments
	newBranches := make([]string, 0)
	for _, branch := range user.BranchOffices {
		if branch != branchID {
			newBranches = append(newBranches, branch)
		}
	}

	update := bson.M{
		"$set": bson.M{
			"branch_offices": newBranches,
			"updated_at":     time.Now(),
		},
	}

	result := sc.config.MongoDB.Collection("users").FindOneAndUpdate(
		c,
		bson.M{"_id": userID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var updatedUser models.User
	if err := result.Decode(&updatedUser); err != nil {
		utils.InternalError(c, "Error updating user")
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// RegisterStaff registers a new staff member
func (sc *StaffController) RegisterStaff(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	// Ensure staff role is assigned
	hasStaffRole := false
	for _, role := range req.Roles {
		if role == models.RoleStaff {
			hasStaffRole = true
			break
		}
	}
	if !hasStaffRole {
		req.Roles = append(req.Roles, models.RoleStaff)
	}

	// Verify company exists if provided
	if req.CompanyID != "" {
		var company models.Company
		err := sc.config.MongoDB.Collection("companies").FindOne(c, bson.M{"_id": req.CompanyID}).Decode(&company)
		if err != nil {
			utils.BadRequest(c, "Company not found")
			return
		}

		// Verify branch offices exist and belong to the company if provided
		if len(req.BranchOffices) > 0 {
			for _, branchID := range req.BranchOffices {
				var branch models.BranchOffice
				err = sc.config.MongoDB.Collection("branch_offices").FindOne(c,
					bson.M{
						"_id":        branchID,
						"company_id": req.CompanyID,
					}).Decode(&branch)
				if err != nil {
					utils.BadRequest(c, "Invalid branch office ID: "+branchID)
					return
				}
			}
		}
	}

	// Check if username already exists
	var existingUser models.User
	err := sc.config.MongoDB.Collection("users").FindOne(c, bson.M{"username": req.Username}).Decode(&existingUser)
	if err == nil {
		utils.BadRequest(c, "Username already taken")
		return
	}

	// Create new user
	user := models.User{
		ID:            primitive.NewObjectID().Hex(),
		Username:      req.Username,
		Password:      req.Password,
		FullName:      req.FullName,
		Roles:         req.Roles,
		CompanyID:     req.CompanyID,
		BranchOffices: req.BranchOffices,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		utils.InternalError(c, "Error creating user")
		return
	}

	// Insert user
	_, err = sc.config.MongoDB.Collection("users").InsertOne(c, user)
	if err != nil {
		utils.InternalError(c, "Error creating user")
		return
	}

	user.Password = "" // Don't send password back

	// Generate token
	token, err := config.GenerateToken(user.ID, user.Roles)
	if err != nil {
		utils.InternalError(c, "Error generating token")
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		Token: token,
		User:  user,
	})
}
