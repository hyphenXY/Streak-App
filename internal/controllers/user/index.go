package user_controller

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hyphenXY/Streak-App/internal/dataproviders"
	"github.com/hyphenXY/Streak-App/internal/models"
	"github.com/hyphenXY/Streak-App/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// POST /user/signIn
func SignIn(c *gin.Context) {
	// 1️⃣ Parse JSON body
	type SignInRequest struct {
		UserName string `json:"userName" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var req SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	var user models.User
	if err := dataprovider.DB.Where("user_name = ?", req.UserName).First(&user).Error; err != nil {
		// User not found
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username! try to remember it."})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		// Wrong password
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password! work your brain or reset it."})
		return
	}

	accessToken, err := utils.GenerateJWT(map[string]any{
		"userId": user.ID,
		"role":   "user",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create refresh token"})
		return
	}

	refreshTokenExpiry := time.Now().Add(30 * 24 * time.Hour) // 30 days
	if err := dataprovider.UpdateUserRefreshToken(user.ID, refreshToken, refreshTokenExpiry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Sign in successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":        user.ID,
			"username":  user.UserName,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"phone":     user.Phone,
		},
	})
}

func SignUp(c *gin.Context) {
	// 1️⃣ Parse JSON body
	type SignUpRequest struct {
		UserName  string `json:"userName" binding:"required"`
		Password  string `json:"password" binding:"required"`
		Email     string `json:"email" binding:"required"`
		FirstName string `json:"firstName" binding:"required"`
		LastName  string `json:"lastName" binding:"required"`
		Phone     string `json:"phone" binding:"required"`
		DoB       string `json:"dob" binding:"required"`
	}
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	var existingUser models.User
	if err := dataprovider.DB.
		Where("user_name = ? OR email = ?", req.UserName, req.Email).
		First(&existingUser).Error; err == nil {
		// Found a record
		c.JSON(http.StatusConflict, gin.H{"error": "username, or email already exists."})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Some other DB error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	err = dataprovider.CreateUser(&models.User{
		UserName:  req.UserName,
		Password:  string(hashedPassword),
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		DOB:       req.DoB,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})

}

// GET /user/homepage/:id
func Homepage(c *gin.Context) {
	userID := c.Param("id")
	// TODO: fetch homepage data for user
	c.JSON(http.StatusOK, gin.H{
		"message": "Homepage data",
		"user_id": userID,
	})
}

// POST /user/markAttendance/:id
func MarkAttendance(c *gin.Context) {
	userID := c.Param("id")
	// TODO: mark attendance for userID
	c.JSON(http.StatusOK, gin.H{
		"message": "Attendance marked",
		"user_id": userID,
	})
}

// GET /user/profile/:id
func Profile(c *gin.Context) {
	userID := c.Param("id")
	// TODO: fetch user profile from DB
	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"name":    "John Doe",
		"email":   "john@example.com",
	})
}

// PATCH /user/profile/:id
func UpdateProfile(c *gin.Context) {
	userID := c.Param("id")

	type UpdateProfileRequest struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: update profile in DB
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated",
		"user_id": userID,
		"name":    req.Name,
		"email":   req.Email,
	})
}

// POST /user/sendOTP
func SendOTP(c *gin.Context) {
	type OTPRequest struct {
		Phone string `json:"phone"`
	}

	var req OTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: send OTP to req.Phone
	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent",
		"phone":   req.Phone,
	})
}

func Enroll(c *gin.Context) {
	classID := c.Param("id")

	classIDUint, err := strconv.ParseUint(classID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID"})
		return
	}
	// Log the enrollment attempt

	userIDVal, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}

	userID, ok := userIDVal.(float64)
	if !ok {
		// if your middleware stored it as float64 or string, convert accordingly
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id type in context"})
		return
	}

	err = dataprovider.EnrollUser(uint(userID), uint(classIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enroll user"})
		return
	}

	// TODO: enroll req.UserID in classID
	c.JSON(http.StatusOK, gin.H{
		"message":  "User enrolled",
		"user_id":  userID,
		"class_id": classID,
	})
}
