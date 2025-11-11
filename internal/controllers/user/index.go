package user_controller

import (
	"bytes"
	"encoding/json"
	"errors"
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
		"role":         "user",
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

	// Verify OTP
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

	dob, err := time.Parse("2006-01-02", req.DoB)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	newUser := &models.User{
		UserName:  req.UserName,
		Password:  string(hashedPassword),
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		DOB:       dob, // parsed time
	}

	if err := dataprovider.CreateUser(newUser); err != nil {
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

// GET /user/classList/:id
func ClassList(c *gin.Context) {
	userID, err := c.Get("userId")
	if !err {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// TODO: fetch class list data for user
	var userClasses []models.User_Classes
	if err := dataprovider.DB.Where("user_id = ?", userID).Find(&userClasses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user classes"})
		return
	}

	// Optionally, fetch class details for each enrollment
	var classes []models.Classes
	classIDs := make([]uint, 0, len(userClasses))
	for _, uc := range userClasses {
		classIDs = append(classIDs, uc.ClassID)
	}
	if len(classIDs) > 0 {
		if err := dataprovider.DB.Where("id IN ?", classIDs).Find(&classes).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch class details"})
			return
		}
	}
	// Merge join date into classes array
	joinedClasses := make([]gin.H, 0, len(classes))
	classJoinMap := make(map[uint]time.Time)
	for _, uc := range userClasses {
		classJoinMap[uc.ClassID] = uc.CreatedAt
	}
	for _, class := range classes {
		joinedClasses = append(joinedClasses, gin.H{
			"class_id":            class.ID,
			"class_name":          class.Name,
			"class_code":          class.ClassCode,
			"created_at":          class.CreatedAt,
			"joined_at":           classJoinMap[class.ID],
			"email":               class.Email,
			"phone":               class.Phone,
			"created_by_admin_id": class.CreatedByAdminId,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"classes": joinedClasses,
	})
}

// POST /user/markAttendance/:id
func MarkAttendance(c *gin.Context) {
	// check if user is enrolled in the class
	userIDVal, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	classID, exists := c.Get("classID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "classID not provided"})
		return
	}

	type MarkAttendanceRequest struct {
		Status string `json:"status" binding:"required"`
	}
	var req MarkAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	classIDFloat, ok2 := classID.(uint)

	if !ok2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId or classID type"})
		return
	}

	println("d", uint(userIDVal.(float64)), classIDFloat)

	err := dataprovider.MarkAttendanceByUser(classIDFloat, uint(userIDVal.(float64)), req.Status)
	if err != nil {
		// Check if attendance is already marked
		if err.Error() == "already marked" {
			c.JSON(http.StatusConflict, gin.H{"error": "Attendance already marked"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark attendance", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance marked", "class_id": classID})
}

// GET /user/profile/:id
func Profile(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := dataprovider.DB.Where("id = ?", userID).First(&user).Error; err != nil {
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
	})
}

func UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("UserId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	type UpdateProfileRequest struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateData := map[string]interface{}{
		"first_name": req.FirstName,
		"last_name":  req.LastName,
		"email":      strings.ToLower(strings.TrimSpace(req.Email)),
	}

	dataprovider.UpdateProfile(updateData, uint(userID.(float64)))

	// TODO: update profile in DB using updateData
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated",
		"user_id": userID,
		"name":    req.FirstName + " " + req.LastName,
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
		err = dataprovider.MarkPhoneVerified(uint(phoneUint))
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

func Enroll(c *gin.Context) {
	classCode := c.Param("classCode")

	classID, err := dataprovider.GetClassIDByCode(classCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Class code not found"})
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
	isEnrolled, err := dataprovider.IfAlreadyEnrolled(uint(userID), uint(classID), &existingEnrollment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check enrollment status"})
		return
	}
	if isEnrolled {
		c.JSON(http.StatusConflict, gin.H{"error": "User already enrolled"})
		return
	}

	err = dataprovider.EnrollUser(uint(userID), uint(classID))
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
	// 1️⃣ Extract refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token missing"})
		return
	}

	// 2️⃣ Find user with this refresh token
	var user models.User
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
		"role":    "user",
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

func ClassDetails(c *gin.Context) {
	classID, exists := c.Get("classID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "classID not provided"})
		return
	}

	classIDUint, ok := classID.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid classID type"})
		return
	}

	// Fetch class details from the database
	class, err := dataprovider.GetClassByID(classIDUint)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch class details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"class": class})
}

func QuickSummary(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	classID, exists := c.Get("classID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "classID not provided"})
		return
	}
	classIDFloat, ok := classID.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid classID type"})
		return
	}

	quickSummary, err := dataprovider.GetUserQuickSummary(uint(userID.(float64)), classIDFloat, "user")
	if err != nil {
		println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch quick summary"})
		return
	}

	// Add summary to response
	c.JSON(http.StatusOK, gin.H{
		"quick_summary": quickSummary,
	})
}

func Calendar(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	classID, exists := c.Get("classID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "classID not provided"})
		return
	}

	classIDFloat, ok := classID.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid classID type"})
		return
	}

	// Get all attendance records for this user in this class
	attendanceRecords, err := dataprovider.GetUserCalendar(uint(userID.(float64)), classIDFloat, "user")
	if err != nil {
		println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attendance records"})
		return
	}

	// Prepare response: list of {date, status}
	calendar := make([]gin.H, 0, len(attendanceRecords))
	for _, record := range attendanceRecords {
		calendar = append(calendar, gin.H{
			"date":   record.CreatedAt.Format("2006-01-02"),
			"status": record.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"class_id": classID,
		"user_id":  userID,
		"calendar": calendar,
	})

}

func LogOutUser(c *gin.Context) {
	// 1️⃣ Extract refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing refresh token"})
		return
	}

	// 2️⃣ Validate and revoke the refresh token
	if err := dataprovider.RevokeUserRefreshToken(refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke refresh token"})
		return
	}

	// 3️⃣ Clear the refresh token cookie
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func Streak(c *gin.Context) {
	classID, exists := c.Get("classID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing classID"})
		return
	}

	userIDVal, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	classIDFloat, ok := classID.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid classID type"})
		return
	}

	curr, best, err := dataprovider.GetUserStreak(uint(userIDVal.(float64)), classIDFloat, "user")
	if err != nil {
		println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch streak data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"currentStreak": curr, "bestStreak": best})
}

func Report(c *gin.Context) {
	classID, exists := c.Get("classID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing classID"})
		return
	}

	userIDVal, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	classIDFloat, ok := classID.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid classID type"})
		return
	}

	report, err := dataprovider.GetUserReport(uint(userIDVal.(float64)), classIDFloat, "user")
	if err != nil {
		println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch report data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"report": report})
}

func ResetPassword(c *gin.Context) {
	type ResetPasswordRequest struct {
		NewPassword string `json:"newPassword" binding:"required"`
		OTP         uint   `json:"otp" binding:"required"`
		Phone       uint   `json:"phone" binding:"required"`
	}
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	userIDVal, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// check if otp is verified in last 10 mins
	isValid, err := dataprovider.IsOTPRecentlyVerified(req.OTP, req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify OTP status"})
		return
	}
	if !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP not verified recently"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	if err := dataprovider.UpdateUserPassword(uint(userIDVal.(float64)), string(hashedPassword)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}
