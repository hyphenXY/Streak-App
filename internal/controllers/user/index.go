package user_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// POST /user/signIn
func SignIn(c *gin.Context) {
	// You can bind JSON here for email/password
	type SignInRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: authenticate user
	c.JSON(http.StatusOK, gin.H{
		"message": "User signed in successfully",
		"email":   req.Email,
	})
}

// POST /user/signUp
func SignUp(c *gin.Context) {
	type SignUpRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: save user to DB
	c.JSON(http.StatusCreated, gin.H{
		"message": "User signed up successfully",
		"name":    req.Name,
	})
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
