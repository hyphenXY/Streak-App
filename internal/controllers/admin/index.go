package admin_controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var user models.Admin
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
		"role":   "admin",
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
	if err := dataprovider.UpdateAdminRefreshToken(user.ID, refreshToken, refreshTokenExpiry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update refresh token"})
		return
	}

	// Set refresh token as HttpOnly cookie
	secure := gin.Mode() == gin.ReleaseMode
	c.SetCookie(
		"refresh_token",
		refreshToken,
		int((30 * 24 * time.Hour).Seconds()), // expiry in seconds
		"/",
		"",     // domain (empty = current domain)
		secure, // secure (true = HTTPS only)
		true,   // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message":      "Sign in successful",
		"access_token": accessToken,
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

// POST /user/signUp
func SignUp(c *gin.Context) {
	// 1️⃣ Parse JSON body
	type SignUpRequest struct {
		UserName  string `json:"userName" binding:"required"`
		Password  string `json:"password" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		FirstName string `json:"firstName" binding:"required"`
		LastName  string `json:"lastName" binding:"required"`
		Phone     string `json:"phone" binding:"required"`
		DoB       string `json:"dob" binding:"required"`
		OTP       string `json:"otp" binding:"required"`
	}
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	phoneUint, err := strconv.ParseUint(req.Phone, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number"})
		return
	}
	isValid, err := dataprovider.IsPhoneVerified(uint(phoneUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify phone"})
		return
	}
	if !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Phone not verified"})
		return
	}

	req.UserName = strings.ToLower(strings.TrimSpace(req.UserName))
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	var existingUser models.Admin
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

	dob, err := time.Parse("2006-01-02", req.DoB)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	newUser := &models.Admin{
		UserName:  req.UserName,
		Password:  string(hashedPassword),
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		DOB:       dob, // parsed time
	}

	if err := dataprovider.CreateAdmin(newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully",
		"user": gin.H{
			"username":  req.UserName,
			"email":     req.Email,
			"firstName": req.FirstName,
			"lastName":  req.LastName,
			"phone":     req.Phone,
		},
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

func PersonalHomepage(c *gin.Context) {
	userID := c.Param("id")
	// TODO: fetch personal homepage data for user
	c.JSON(http.StatusOK, gin.H{
		"message": "Personal homepage data",
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
		err := dataprovider.MarkPhoneVerified(uint(phoneUint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark phone as verified"})
			return
		}
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

func RefreshTokenUser(c *gin.Context) {
	// 1️⃣ Extract refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token missing"})
		return
	}

	// 2️⃣ Find user with this refresh token
	var user models.Admin
	if err := dataprovider.DB.
		Where("refresh_token = ?", refreshToken).
		First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	// 3️⃣ Check expiry
	if time.Now().After(*user.RefreshTokenExpiry) {
		// clear cookie
		c.SetCookie("refresh_token", "", -1, "/", "", true, true)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token expired"})
		return
	}

	// 4️⃣ Generate new access token
	accessToken, err := utils.GenerateJWT(map[string]any{
		"user_id": user.ID,
		"role":    "admin",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create access token"})
		return
	}

	// 5️⃣ Rotate refresh token
	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	expiry := time.Now().Add(30 * 24 * time.Hour)
	if err := dataprovider.UpdateUserRefreshToken(user.ID, newRefreshToken, expiry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update refresh token"})
		return
	}

	c.SetCookie("refresh_token", newRefreshToken, int((30 * 24 * time.Hour).Seconds()), "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

func CreateClass(c *gin.Context) {
	adminId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	type CreateClassRequest struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	var req CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	class := models.Classes{
		Name:             req.Name,
		Email:            req.Email,
		Phone:            req.Phone,
		CreatedByAdminId: uint(adminId.(float64)),
	}

	if err := dataprovider.CreateClass(&class); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create class"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Class created successfully", "class_id": class.ID})
}

func StudentsList(c *gin.Context) {
	classId, exists := c.Get("classId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	students, err := dataprovider.GetStudentsByClassID(classId.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch students"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"students": students})
}
