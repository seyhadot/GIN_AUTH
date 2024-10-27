package controllers

import (
	"loan/config"
	"loan/models"
	"loan/utils"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	config *config.Config
}

func NewAuthController(config *config.Config) *AuthController {
	return &AuthController{config: config}
}

func (ac *AuthController) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	// Check if username already exists
	var existingUser models.User
	err := ac.config.MongoDB.Collection("users").FindOne(c, bson.M{"username": req.Username}).Decode(&existingUser)
	if err == nil {
		utils.BadRequest(c, "Username already taken")
		return
	}

	// Generate new ID as string
	id := primitive.NewObjectID().Hex()

	// Create new user with string ID
	user := models.User{
		ID:        id,
		Username:  req.Username,
		Password:  req.Password,
		Roles:     req.Roles,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set default roles if none provided
	user.SetDefaultRoles()

	// Hash password
	if err := user.HashPassword(); err != nil {
		utils.InternalError(c, "Error creating user")
		return
	}

	// Insert user
	_, err = ac.config.MongoDB.Collection("users").InsertOne(c, user)
	if err != nil {
		utils.InternalError(c, "Error creating user")
		return
	}

	user.Password = "" // Don't send password back

	// Generate token with string ID
	token, err := config.GenerateToken(id, user.Roles)
	if err != nil {
		utils.InternalError(c, "Error generating token")
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

func (ac *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	// Find user by username
	var user models.User
	err := ac.config.MongoDB.Collection("users").FindOne(c, bson.M{"username": req.Username}).Decode(&user)
	if err != nil {
		utils.BadRequest(c, "Invalid credentials")
		return
	}

	// Compare password
	if err := user.ComparePassword(req.Password); err != nil {
		utils.BadRequest(c, "Invalid credentials")
		return
	}

	// Generate token with roles
	token, err := config.GenerateToken(user.ID, user.Roles)
	if err != nil {
		utils.InternalError(c, "Error generating token")
		return
	}

	user.Password = "" // Don't send password back
	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

// GetCurrentProfile gets the current user's profile
func (ac *AuthController) GetCurrentProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.BadRequest(c, "User not found")
		return
	}

	var user models.User
	err := ac.config.MongoDB.Collection("users").FindOne(c, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		utils.BadRequest(c, "User not found")
		return
	}

	user.Password = "" // Don't send password
	c.JSON(http.StatusOK, user)
}

// UpdateProfile updates the current user's profile
func (ac *AuthController) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.BadRequest(c, "User not found")
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	// Prepare update data
	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	// Only update fields that are provided
	if req.FullName != "" {
		update["$set"].(bson.M)["full_name"] = req.FullName
	}
	if req.Bio != "" {
		update["$set"].(bson.M)["bio"] = req.Bio
	}
	if req.Avatar != "" {
		update["$set"].(bson.M)["avatar"] = req.Avatar
	}

	// Update user using string ID
	result := ac.config.MongoDB.Collection("users").FindOneAndUpdate(
		c,
		bson.M{"_id": userID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var updatedUser models.User
	if err := result.Decode(&updatedUser); err != nil {
		utils.InternalError(c, "Error updating profile")
		return
	}

	updatedUser.Password = "" // Don't send password
	c.JSON(http.StatusOK, updatedUser)
}

// UpdatePassword updates the current user's password
func (ac *AuthController) UpdatePassword(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.BadRequest(c, "User not found")
		return
	}

	var req models.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body")
		return
	}

	// Get current user using string ID
	var user models.User
	err := ac.config.MongoDB.Collection("users").FindOne(c, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		utils.BadRequest(c, "User not found")
		return
	}

	// Verify current password
	if err := user.ComparePassword(req.CurrentPassword); err != nil {
		utils.BadRequest(c, "Current password is incorrect")
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalError(c, "Error updating password")
		return
	}

	// Update password in database using string ID
	update := bson.M{
		"$set": bson.M{
			"password":   string(hashedPassword),
			"updated_at": time.Now(),
		},
	}

	result := ac.config.MongoDB.Collection("users").FindOneAndUpdate(
		c,
		bson.M{"_id": userID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var updatedUser models.User
	if err := result.Decode(&updatedUser); err != nil {
		utils.InternalError(c, "Error updating password")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}

// GetAllUsers gets all users with pagination
func (ac *AuthController) GetAllUsers(c *gin.Context) {
	// Get pagination parameters from query
	page := 1
	limit := 10
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Calculate skip value
	skip := (page - 1) * limit

	// Get total count
	total, err := ac.config.MongoDB.Collection("users").CountDocuments(c, bson.M{})
	if err != nil {
		utils.InternalError(c, "Error counting users")
		return
	}

	// Set options for pagination and sorting
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	// Fetch users from database
	cursor, err := ac.config.MongoDB.Collection("users").Find(c, bson.M{}, opts)
	if err != nil {
		utils.InternalError(c, "Error fetching users")
		return
	}
	defer cursor.Close(c)

	var users []models.User
	if err = cursor.All(c, &users); err != nil {
		utils.InternalError(c, "Error parsing users")
		return
	}

	// Remove sensitive information
	for i := range users {
		users[i].Password = ""
	}

	// Calculate pagination metadata
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	hasNext := page < totalPages
	hasPrev := page > 1

	// Return response with pagination metadata
	response := gin.H{
		"users": users,
		"pagination": gin.H{
			"total":       total,
			"total_pages": totalPages,
			"page":        page,
			"limit":       limit,
			"has_next":    hasNext,
			"has_prev":    hasPrev,
		},
	}

	c.JSON(http.StatusOK, response)
}
