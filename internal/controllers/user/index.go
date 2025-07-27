package user_controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	dataprovider "github.com/hyphenXY/Streak-App/internal/dataproviders"
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
	classID := c.Param("id")

	classIDUint, err := strconv.ParseUint(classID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID"})
		return
	}

	ifClassExists, err := dataprovider.IfClassExists(uint(classIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check class existence"})
		return
	}
	if !ifClassExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
		return
	}

	// check if user is enrolled in the class
	userIDVal, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if err := dataprovider.IfAlreadyEnrolled(uint(userIDVal.(float64)), uint(classIDUint), &models.User_Classes{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check enrollment status"})
		return
	}

	err = dataprovider.MarkAttendanceByUser(uint(classIDUint), uint(userIDVal.(float64)))
	if err != nil {
		// Check if attendance is already marked
		if err.Error() == "already marked" {
			c.JSON(http.StatusConflict, gin.H{"error": "Attendance already marked"})
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Attendance record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark attendance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance marked", "class_id": classID})
}

// GET /user/profile/:id
func Profile(c *gin.Context) {
	userID := c.Param("id")

	userIdVal, err := strconv.ParseFloat(userID, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// userIdval, err := c.Get("userId")
	// if !err {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
	// 	return
	// }

	// if userIdval.(float64) != strconv.ParseFloat(userID, 64) {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own profile"})
	// 	return
	// }

	var user models.User
	if err := dataprovider.DB.Where("id = ?", userIdVal).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.FirstName + " " + user.LastName,
		"email": user.Email,
		"phone": user.Phone,
		"dob":   user.DOB,
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

	otp, err := utils.GenerateOTP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
		return
	}

	payload := map[string]string{
		"phone":       "+91" + req.Phone,
		"otp":         otp,
		"gateway_key": os.Getenv("FAZPASS_GATEWAY_KEY"),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal payload"})
		return
	}

	url := "https://api.fazpass.com/v1/otp/send"

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create HTTP request"})
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+os.Getenv("FAZPASS_MERCHANT_KEY"))

	// Print all content of httpReq for debugging
	println("HTTP Request Method:", httpReq.Method)
	println("HTTP Request URL:", httpReq.URL.String())
	for k, v := range httpReq.Header {
		println("Header:", k, "=", v[0])
	}
	if httpReq.Body != nil {
		bodyBytes, _ := payloadBytes, error(nil)
		if bodyBytes == nil {
			bodyBytes, _ = io.ReadAll(httpReq.Body)
		}
		println("Body:", string(bodyBytes))
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		// Print the error for debugging
		println("Error sending OTP:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		errorMsg := buf.String()
		println("Failed to send OTP, status code:", resp.Status, "response:", errorMsg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":        "Failed to send OTP",
			"status_code":  resp.StatusCode,
			"status":       resp.Status,
			"response_msg": errorMsg,
		})
		return
	}

	phoneUint, err := strconv.ParseUint(req.Phone, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number"})
		return
	}
	err = dataprovider.StoreOTP(uint(phoneUint), otp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent",
		"phone":   req.Phone,
	})
}

func VerifyOTP(c *gin.Context) {
	type VerifyRequest struct {
		Phone string `json:"phone"`
		OTP   string `json:"otp"`
	}

	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	phoneUint, err := strconv.ParseUint(req.Phone, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number"})
		return
	}

	isValid, err := dataprovider.VerifyOTP(uint(phoneUint), req.OTP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify OTP"})
		return
	}
	switch isValid {
	case "Verified!":
		c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
		return
	case "Expired!":
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP expired"})
		return
	case "Wrong!":
		c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP not found"})
		return
	case "Failed!":
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OTP verification failed"})
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": isValid})
		return
	}
}

func Enroll(c *gin.Context) {
	classID := c.Param("id")

	classIDUint, err := strconv.ParseUint(classID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID"})
		return
	}
	ifClassExists, err := dataprovider.IfClassExists(uint(classIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check class existence"})
		return
	}
	if !ifClassExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
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

	// check if user is already enrolled
	var existingEnrollment models.User_Classes
	if err := dataprovider.IfAlreadyEnrolled(uint(userID), uint(classIDUint), &existingEnrollment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check enrollment status"})
		return
	}
	if existingEnrollment.ID != 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "User already enrolled"})
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

func RefreshTokenUser(c *gin.Context) {
	type RefreshRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// 🔍 Look up user with this refresh token
	var user models.User
	if err := dataprovider.DB.
		Where("refresh_token = ?", req.RefreshToken).
		First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	// ⏳ Check expiry
	if time.Now().After(*user.RefreshTokenExpiry) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token expired"})
		return
	}

	// ✅ Generate new access token
	accessToken, err := utils.GenerateJWT(map[string]any{
		"user_id": user.ID,
		"role":    "user",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

func ClassDetails(c *gin.Context) {
	type ClassDetailsRequest struct {
		ClassID uint `json:"class_id" binding:"required"`
	}

	var req ClassDetailsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var class models.Classes
	if err := dataprovider.DB.Where("id = ?", req.ClassID).First(&class).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch class details"})
		return
	}

	// take out admin name form id
	var admin models.Admin
	if err := dataprovider.AdminNameById(class.CreatedByAdminId, &admin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch admin details"})
		return
	}
	adminName := admin.FirstName + " " + admin.LastName

	c.JSON(http.StatusOK, gin.H{
		"class_id":    class.ID,
		"class_name":  class.Name,
		"description": class.Name,
		"class_email": class.Email,
		"class_phone": class.Phone,
		"class_owner": adminName,
	})
}
