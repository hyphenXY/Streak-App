package dataprovider

import (
	// "gorm.io/gorm"
	"errors"
	"time"

	"github.com/hyphenXY/Streak-App/internal/models"
	"gorm.io/gorm"
)

func IfClassExists(classID uint) (bool, error) {
	var count int64
	err := DB.Model(&models.Classes{}).Where("id = ?", classID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func CreateClass(class *models.Classes) error {
	return DB.Create(class).Error
}

func MarkAttendanceByUser(classID uint, userID uint, status string) error {
	// check in attendances table if record exists
	var attendance models.Attendance
	err := DB.Model(&models.Attendance{}).
		Where("class_id = ? AND marked_by_id = ? AND marked_by_role = ? AND DATE(created_at) = CURRENT_DATE", classID, userID, "user").
		First(&attendance).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		attendance = models.Attendance{
			ClassID:      classID,
			MarkedById:   userID,
			MarkedByRole: "user",
			Status:       status,
		}
		return DB.Create(&attendance).Error
	}
	if err != nil {
		return err
	}
	return errors.New("already marked")
}

func MarkAttendanceByAdmin(classID uint, userID uint) error {
	// check in attendances table if record exists
	var attendance models.Attendance
	err := DB.Model(&models.Attendance{}).
		Where("class_id = ? AND marked_by_id = ? AND marked_by_role = ? AND DATE(created_at) = CURRENT_DATE", classID, userID, "admin").
		First(&attendance).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		attendance = models.Attendance{
			ClassID:      classID,
			MarkedById:   userID,
			MarkedByRole: "admin",
			Status:       "present",
		}
		return DB.Create(&attendance).Error
	}
	if err != nil {
		return err
	}
	return errors.New("already marked")
}

func IsUserAdmin(userID uint, classID uint) (bool, error) {
	var count int64
	err := DB.Model(&models.Classes{}).Where("created_by_admin_id = ? AND id = ?", userID, classID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetStudentsByClassID(classID uint) ([]models.User, error) {
	var students []models.User
	err := DB.Joins("JOIN enrollments ON enrollments.user_id = users.id").
		Where("enrollments.class_id = ?", classID).
		Find(&students).Error
	if err != nil {
		return nil, err
	}
	return students, nil
}

func GetClassIDByCode(classCode string) (uint, error) {
	var class models.Classes
	err := DB.Where("class_code = ?", classCode).First(&class).Error
	if err != nil {
		return 0, err
	}
	return class.ID, nil
}

func GetClassByID(classID uint) (*models.Classes, error) {
	var class models.Classes
	err := DB.Where("id = ?", classID).First(&class).Error
	if err != nil {
		return nil, err
	}
	return &class, nil
}

func GetUserCalendar(userID uint, classID uint, role string) ([]models.Attendance, error) {
	var attendanceRecords []models.Attendance
	err := DB.
		Where("marked_by_id = ? AND marked_by_role = ? AND class_id = ?", userID, role, classID).
		Order("created_at ASC").
		Find(&attendanceRecords).Error
	if err != nil {
		return nil, err
	}
	return attendanceRecords, err
}

func GetUserStreak(userID uint, classID uint, role string) (int, int, error) {
	var attendances []models.Attendance
	err := DB.
		Where("marked_by_id = ?  AND marked_by_role = ? AND class_id = ?", userID, role, classID).
		Order("created_at ASC").
		Find(&attendances).Error
	if err != nil {
		return 0, 0, err
	}

	// Calculate the best streak
	bestStreak := 0
	currentStreak := 0

	// Group records by date (keep the last status for a given date)
	type dayRec struct {
		date   time.Time
		status string
	}
	var days []dayRec
	for _, a := range attendances {
		d := time.Date(a.CreatedAt.Year(), a.CreatedAt.Month(), a.CreatedAt.Day(), 0, 0, 0, 0, a.CreatedAt.Location())
		if len(days) == 0 || !days[len(days)-1].date.Equal(d) {
			days = append(days, dayRec{date: d, status: a.Status})
		} else {
			// overwrite with the later status on the same day
			days[len(days)-1].status = a.Status
		}
	}

	var prevDate time.Time
	var prevSet bool
	for _, day := range days {
		if day.status == "present" {
			if prevSet && day.date.Equal(prevDate.AddDate(0, 0, 1)) {
				// consecutive day
				currentStreak++
			} else {
				// start new streak
				currentStreak = 1
			}
			prevDate = day.date
			prevSet = true
		} else {
			if currentStreak > bestStreak {
				bestStreak = currentStreak
			}
			currentStreak = 0
			prevSet = false
		}
	}
	if currentStreak > bestStreak {
		bestStreak = currentStreak
	}
	return currentStreak, bestStreak, nil
}

func GetUserQuickSummary(userID uint, classID uint, role string) (map[string]interface{}, error) {
	var todayAttendance models.Attendance
	err := DB.Model(&models.Attendance{}).
		Where("marked_by_id = ? AND marked_by_role = ? AND class_id = ? AND DATE(created_at) = CURRENT_DATE", userID, role, classID).
		First(&todayAttendance).Error

	todayStatus := "unmarked"
	if err == nil {
		switch todayAttendance.Status {
		case "present":
			todayStatus = "present"
		case "absent":
			todayStatus = "absent"
		default:
			todayStatus = todayAttendance.Status
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var currentWeekPresent int64
	if err := DB.Model(&models.Attendance{}).
		Where("marked_by_id = ? AND marked_by_role = ? AND class_id = ? AND status = ? AND YEARWEEK(created_at) = YEARWEEK(CURRENT_DATE)", userID, role, classID, "present").
		Count(&currentWeekPresent).Error; err != nil {
		return nil, err
	}

	var currentWeekAbsent int64
	if err := DB.Model(&models.Attendance{}).
		Where("marked_by_id = ? AND marked_by_role = ? AND class_id = ? AND status = ? AND YEARWEEK(created_at) = YEARWEEK(CURRENT_DATE)", userID, role, classID, "absent").
		Count(&currentWeekAbsent).Error; err != nil {
		return nil, err
	}

	var currentWeekNotMarked int64
	if err := DB.Model(&models.Attendance{}).
		Distinct("DATE(created_at)").Where("marked_by_id = ? AND marked_by_role = ? AND class_id = ? AND status = ?", userID, role, classID, "not_marked").
		Count(&currentWeekNotMarked).Error; err != nil {
		return nil, err
	}

	var totalPresent int64
	if err := DB.Model(&models.Attendance{}).
		Where("marked_by_id = ? AND marked_by_role = ? AND class_id = ? AND status = ?", userID, role, classID, "present").
		Count(&totalPresent).Error; err != nil {
		return nil, err
	}

	var totalAbsent int64
	if err := DB.Model(&models.Attendance{}).
		Where("marked_by_id = ? AND marked_by_role = ? AND class_id = ? AND status = ?", userID, role, classID, "absent").
		Count(&totalAbsent).Error; err != nil {
		return nil, err
	}

	var totalSessions int64
	if err := DB.Model(&models.Attendance{}).
		Distinct("created_at").
		Where("class_id = ? AND marked_by_role = ?", classID, role).
		Count(&totalSessions).Error; err != nil {
		return nil, err
	}

	totalNotMarked := max(totalSessions-(totalPresent+totalAbsent), 0)

	// quick summary map (kept here for future use; function returns total_not_marked)
	summary := map[string]interface{}{
		"today_status":            todayStatus,
		"current_week_present":    currentWeekPresent,
		"current_week_absent":     currentWeekAbsent,
		"current_week_not_marked": currentWeekNotMarked,
		"total_present":           totalPresent,
		"total_absent":            totalAbsent,
		"total_not_marked":        totalNotMarked,
	}

	return summary, nil
}

func IfAlreadyEnrolled(userID uint, classID uint, enrollment *models.User_Classes) (bool, error) {
	err := DB.Where("user_id = ? AND class_id = ?", userID, classID).First(enrollment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil // Not enrolled
		}
		return false, err // Other error
	}
	return true, nil // Already enrolled
}

func EnrollUser(userID uint, classID uint) error {
	enrollment := models.User_Classes{
		UserID:  userID,
		ClassID: classID,
	}
	result := DB.Create(&enrollment)
	return result.Error
}

func GetClassSummary(classID uint) (map[string]interface{}, error) {
	summary := make(map[string]interface{})

	var totalStudents int64
	if err := DB.Model(&models.User_Classes{}).
		Where("class_id = ?", classID).
		Count(&totalStudents).Error; err != nil {
		return nil, err
	}
	summary["total_students"] = totalStudents

	var totalPresent int64
	if err := DB.Model(&models.Attendance{}).
		Where("class_id = ? AND status = ?", classID, "present").
		Count(&totalPresent).Error; err != nil {
		return nil, err
	}
	summary["total_present"] = totalPresent

	var totalAbsent int64
	if err := DB.Model(&models.Attendance{}).
		Where("class_id = ? AND status = ?", classID, "absent").
		Count(&totalAbsent).Error; err != nil {
		return nil, err
	}
	summary["total_absent"] = totalAbsent

	// Current week present/absent (uses YEARWEEK to match earlier queries)
	var currentWeekPresent int64
	if err := DB.Model(&models.Attendance{}).
		Where("class_id = ? AND status = ? AND YEARWEEK(created_at) = YEARWEEK(CURRENT_DATE)", classID, "present").
		Count(&currentWeekPresent).Error; err != nil {
		return nil, err
	}
	summary["current_week_present"] = currentWeekPresent

	var currentWeekAbsent int64
	if err := DB.Model(&models.Attendance{}).
		Where("class_id = ? AND status = ? AND YEARWEEK(created_at) = YEARWEEK(CURRENT_DATE)", classID, "absent").
		Count(&currentWeekAbsent).Error; err != nil {
		return nil, err
	}
	summary["current_week_absent"] = currentWeekAbsent

	// // Current month present/absent
	// var currentMonthPresent int64
	// if err := DB.Model(&models.Attendance{}).
	// 	Where("class_id = ? AND status = ? AND MONTH(created_at) = MONTH(CURRENT_DATE) AND YEAR(created_at) = YEAR(CURRENT_DATE)", classID, "present").
	// 	Count(&currentMonthPresent).Error; err != nil {
	// 	return nil, err
	// }
	// summary["current_month_present"] = currentMonthPresent

	// var currentMonthAbsent int64
	// if err := DB.Model(&models.Attendance{}).
	// 	Where("class_id = ? AND status = ? AND MONTH(created_at) = MONTH(CURRENT_DATE) AND YEAR(created_at) = YEAR(CURRENT_DATE)", classID, "absent").
	// 	Count(&currentMonthAbsent).Error; err != nil {
	// 	return nil, err
	// }
	// summary["current_month_absent"] = currentMonthAbsent

	return summary, nil
}

func GetTodaySummary(classID uint) (map[string]interface{}, error) {
	summary := make(map[string]interface{})

	var totalPresent int64
	if err := DB.Model(&models.Attendance{}).
		Where("class_id = ? AND status = ? AND DATE(created_at) = CURDATE()", classID, "present").
		Count(&totalPresent).Error; err != nil {
		return nil, err
	}
	summary["total_present"] = totalPresent

	var totalAbsent int64
	if err := DB.Model(&models.Attendance{}).
		Where("class_id = ? AND status = ? AND DATE(created_at) = CURDATE()", classID, "absent").
		Count(&totalAbsent).Error; err != nil {
		return nil, err
	}
	summary["total_absent"] = totalAbsent

	var totalStudents int64
	if err := DB.Model(&models.User_Classes{}).
		Where("class_id = ?", classID).
		Count(&totalStudents).Error; err != nil {
		return nil, err
	}
	summary["total_students"] = totalStudents

	return summary, nil
}

func GetUserReport(id uint, classID uint, role string) (map[string]interface{}, error) {
	var presentMonth, absentMonth, notMarkedMonth int64
	if err := DB.Model(&models.Attendance{}).
		Where("marked_by_id = ? AND marked_by_role = ? AND class_id = ? AND status = ? AND MONTH(created_at) = MONTH(CURRENT_DATE) AND YEAR(created_at) = YEAR(CURRENT_DATE)", id, role, classID, "present").
		Count(&presentMonth).Error; err != nil {
		return nil, err
	}
	if err := DB.Model(&models.Attendance{}).
		Where("marked_by_id = ? AND marked_by_role = ? AND class_id = ? AND status = ? AND MONTH(created_at) = MONTH(CURRENT_DATE) AND YEAR(created_at) = YEAR(CURRENT_DATE)", id, role, classID, "absent").
		Count(&absentMonth).Error; err != nil {
		return nil, err
	}
	if err := DB.Model(&models.Attendance{}).
		Where("marked_by_id = ? AND marked_by_role = ? AND class_id = ? AND status = ? AND MONTH(created_at) = MONTH(CURRENT_DATE) AND YEAR(created_at) = YEAR(CURRENT_DATE)", id, role, classID, "not_marked").
		Count(&notMarkedMonth).Error; err != nil {
		return nil, err
	}

	var presentYear, absentYear, notMarkedYear int64
	if err := DB.Model(&models.Attendance{}).
		Where("marked_by_id = ? AND marked_by_role = ? AND class_id = ? AND status = ? AND YEAR(created_at) = YEAR(CURRENT_DATE)", id, role, classID, "present").
		Count(&presentYear).Error; err != nil {
		return nil, err
	}
	if err := DB.Model(&models.Attendance{}).
		Where("marked_by_id = ? AND marked_by_role = ? AND class_id = ? AND status = ? AND YEAR(created_at) = YEAR(CURRENT_DATE)", id, role, classID, "absent").
		Count(&absentYear).Error; err != nil {
		return nil, err
	}
	if err := DB.Model(&models.Attendance{}).
		Where("marked_by_id = ? AND marked_by_role = ? AND class_id = ? AND status = ? AND YEAR(created_at) = YEAR(CURRENT_DATE)", id, role, classID, "not_marked").
		Count(&notMarkedYear).Error; err != nil {
		return nil, err
	}

	report := make(map[string]interface{})
	report["current_month"] = map[string]int64{
		"present":    presentMonth,
		"absent":     absentMonth,
		"not_marked": notMarkedMonth,
	}
	report["current_year"] = map[string]int64{
		"present":    presentYear,
		"absent":     absentYear,
		"not_marked": notMarkedYear,
	}
	return report, nil
}

func GetClassReport(classID uint) (map[string]interface{}, error) {
	var presentMonth, absentMonth, notMarkedMonth int64
	if err := DB.Model(&models.Attendance{}).
		Where("class_id = ? AND status = ? AND MONTH(created_at) = MONTH(CURRENT_DATE) AND YEAR(created_at) = YEAR(CURRENT_DATE)", classID, "present").
		Count(&presentMonth).Error; err != nil {
		return nil, err
	}
	if err := DB.Model(&models.Attendance{}).
		Where("class_id = ? AND status = ? AND MONTH(created_at) = MONTH(CURRENT_DATE) AND YEAR(created_at) = YEAR(CURRENT_DATE)", classID, "absent").
		Count(&absentMonth).Error; err != nil {
		return nil, err
	}
	if err := DB.Model(&models.Attendance{}).
		Where("class_id = ? AND status = ? AND MONTH(created_at) = MONTH(CURRENT_DATE) AND YEAR(created_at) = YEAR(CURRENT_DATE)", classID, "not_marked").
		Count(&notMarkedMonth).Error; err != nil {
		return nil, err
	}

	var presentYear, absentYear, notMarkedYear int64
	if err := DB.Model(&models.Attendance{}).
		Where("class_id = ? AND status = ? AND YEAR(created_at) = YEAR(CURRENT_DATE)", classID, "present").
		Count(&presentYear).Error; err != nil {
		return nil, err
	}
	if err := DB.Model(&models.Attendance{}).
		Where("class_id = ? AND status = ? AND YEAR(created_at) = YEAR(CURRENT_DATE)", classID, "absent").
		Count(&absentYear).Error; err != nil {
		return nil, err
	}
	if err := DB.Model(&models.Attendance{}).
		Where("class_id = ? AND status = ? AND YEAR(created_at) = YEAR(CURRENT_DATE)", classID, "not_marked").
		Count(&notMarkedYear).Error; err != nil {
		return nil, err
	}

	report := make(map[string]interface{})
	report["current_month"] = map[string]int64{
		"present":    presentMonth,
		"absent":     absentMonth,
		"not_marked": notMarkedMonth,
	}
	report["current_year"] = map[string]int64{
		"present":    presentYear,
		"absent":     absentYear,
		"not_marked": notMarkedYear,
	}
	return report, nil
}
