# Streak-App API Documentation

Base URL
- Default: http://0.0.0.0:8080

Global Middlewares
- CORS: CORSMiddleware
- Rate Limiting: NewClientLimiter(5, 10) with LimitMiddleware
- Security Headers: SecurityHeaders

Auth Middlewares
- AuthRootMiddleware: Root-level authenticated access
- AuthAdminMiddleware: Admin authentication required
- AuthUserMiddleware: User authentication required
- IsAdminClass: Admin must be authorized for the target class
- IsUserClass: User must be enrolled/authorized for the target class
- IsAllowedToEnroll: Allow enrollment based on server policies
- IsUserEnrolledInClass: Target user must be in the class

Notes
- Path params indicated by :param
- Unless otherwise noted, payloads and responses are JSON
- Detailed logic and error cases are implemented in controller functions referenced per route

## Root Routes (/root)

### POST /root/signIn
- Handler: root_controller.SignIn
- Body
```json
{
  "email": "string",
  "password": "string"
}
```
- Response 200
```json
{
  "message": "User signed in successfully",
  "email": "john@example.com"
}
```

### POST /root/register
- Handler: root_controller.Register
- Body
```json
{
  "name": "string",
  "email": "string",
  "password": "string"
}
```
- Response 201
```json
{
  "message": "User signed up successfully",
  "name": "John Doe"
}
```

### GET /root/health-check
- Handler: root_controller.HealthCheck
- Response 200
```json
{
  "status": "success",
  "message": "API is healthy"
}
```

Protected (AuthRootMiddleware)

### GET /root/homepage/:id
- Handler: root_controller.Homepage
- Path Params: id (string)
- Response 200
```json
{
  "message": "Homepage data",
  "user_id": "123"
}
```

### GET /root/profile/:id
- Handler: root_controller.Profile
- Path Params: id (string)
- Response 200
```json
{
  "user_id": "123",
  "name": "John Doe",
  "email": "john@example.com"
}
```

### PATCH /root/profile/:id
- Handler: root_controller.UpdateProfile
- Path Params: id (string)
- Body
```json
{
  "name": "string",
  "email": "string"
}
```
- Response 200
```json
{
  "message": "Profile updated",
  "user_id": "123",
  "name": "Jane Doe",
  "email": "jane@example.com"
}
```

### DELETE /root/admin/:id
- Handler: root_controller.DeleteAdmin
- Path Params: id (string)
- Response 200
```json
{
  "message": "Admin user deleted",
  "user_id": "123"
}
```

## Admin Routes (/admin)

Public

### POST /admin/signIn
- Handler: admin_controller.SignIn
- Body
```json
{
  "userName": "string",
  "password": "string"
}
```
- Sets cookie: refresh_token (HttpOnly)
- Response 200
```json
{
  "message": "Sign in successful",
  "role": "admin",
  "access_token": "<jwt>",
  "user": {
    "id": 1,
    "username": "admin1",
    "email": "a@example.com",
    "firstName": "Alice",
    "lastName": "Admin",
    "phone": "9999999999"
  }
}
```
- Error responses: 400, 401

### POST /admin/signUp
- Handler: admin_controller.SignUp
- Body
```json
{
  "userName": "string",
  "password": "string",
  "email": "string",
  "firstName": "string",
  "lastName": "string",
  "phone": "string (digits)",
  "dob": "YYYY-MM-DD",
  "otp": "string"
}
```
- Response 201
```json
{
  "message": "User created successfully",
  "user": {
    "username": "admin1",
    "email": "a@example.com",
    "firstName": "Alice",
    "lastName": "Admin",
    "phone": "9999999999"
  }
}
```

### POST /admin/sendOTP
- Handler: admin_controller.SendOTP
- Body
```json
{
  "phone": "string (digits)"
}
```
- Response 200
```json
{
  "message": "OTP sent",
  "phone": "9999999999"
}
```
- Errors: 400, 500

### POST /admin/verifyOTP
- Handler: admin_controller.VerifyOTP
- Body
```json
{
  "phone": "string (digits)",
  "otp": "string"
}
```
- Response 200
```json
{
  "message": "OTP verified successfully"
}
```
- Errors: 400 (expired/bad), 401 (wrong), 500

### POST /admin/refreshToken
- Handler: admin_controller.RefreshTokenUser
- Cookie: refresh_token
- Response 200
```json
{
  "access_token": "<jwt>"
}
```
- Errors: 401, 500

Protected (AuthAdminMiddleware)

### GET /admin/classList
- Handler: admin_controller.ClassList
- Response 200
```json
{
  "admin_id": 1,
  "classList": [
    {
      "ID": 1,
      "Name": "Class A",
      "Email": "",
      "Phone": "",
      "CreatedByAdminId": 1,
      "ClassCode": "ABC123",
      "CreatedAt": "2024-01-01T00:00:00Z",
      "UpdatedAt": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### GET /admin/profile
- Handler: admin_controller.Profile
- Response 200
```json
{
  "user": {
    "ID": 1,
    "FirstName": "Alice",
    "LastName": "Admin",
    "Email": "a@example.com",
    "Phone": "9999999999",
    "UserName": "admin1",
    "DOB": "2024-01-01T00:00:00Z",
    "CreatedAt": "2024-01-01T00:00:00Z",
    "UpdatedAt": "2024-01-01T00:00:00Z"
  }
}
```
- Errors: 401

### PATCH /admin/profile
- Handler: admin_controller.UpdateProfile
- Body
```json
{
  "firstName": "string",
  "lastName": "string",
  "email": "string"
}
```
- Response 200
```json
{
  "message": "Profile updated",
  "user_id": 1,
  "name": "Alice Admin",
  "email": "a@example.com"
}
```
- Errors: 400, 500

### POST /admin/createClass
- Handler: admin_controller.CreateClass
- Body
```json
{
  "name": "string",
  "email": "string",
  "phone": "string"
}
```
- Response 201
```json
{
  "message": "Class created successfully",
  "class_id": 1,
  "class_code": "ABC123",
  "name": "Class A",
  "email": "",
  "phone": ""
}
```
- Errors: 400, 500

### POST /admin/logOutAdmin
- Handler: admin_controller.LogOutAdmin
- Cookie: refresh_token
- Response 200
```json
{
  "message": "Logged out successfully"
}
```
- Errors: 400, 500

### GET /admin/resetPassword
- Handler: admin_controller.ResetPassword
- Body
```json
{
  "newPassword": "string",
  "phone": 0,
  "otp": 0
}
```
- Response 200
```json
{
  "message": "Password reset successful"
}
```
- Errors: 400, 401, 500

Protected (AuthAdminMiddleware + IsAdminClass)

### GET /admin/quickSummary/:classId
- Handler: admin_controller.QuickSummary
- Response 200
```json
{
  "summary": {
    "total_students": 10,
    "total_present": 8,
    "total_absent": 2,
    "current_week_present": 5,
    "current_week_absent": 1
  }
}
```
- Errors: 400, 500

### GET /admin/todaySummary/:classId
- Handler: admin_controller.TodaySummary
- Response 200
```json
{
  "summary": {
    "total_present": 8,
    "total_absent": 2,
    "total_students": 10
  }
}
```
- Errors: 400, 500

### GET /admin/calendar/:classId
- Handler: admin_controller.Calendar
- Response 200
```json
{
  "class_id": 1,
  "user_id": 1,
  "calendar": [
    {"date": "2024-01-01", "status": "present"},
    {"date": "2024-01-02", "status": "absent"}
  ]
}
```
- Errors: 400, 500

### POST /admin/markAttendance/:classId
- Handler: admin_controller.MarkAttendance
- Response 200
```json
{ "message": "Attendance marked", "class_id": 1 }
```
- Errors: 400, 401, 403, 404, 409, 500

### GET /admin/studentsList/:classId
- Handler: admin_controller.StudentsList
- Response 200
```json
{
  "students": [
    { "ID": 10, "FirstName": "Bob", "LastName": "User", "Email": "b@example.com", "Phone": "9999999998", "UserName": "bob" }
  ]
}
```
- Errors: 400, 401, 403, 500

### GET /admin/streak/:classId
- Handler: admin_controller.Streak
- Response 200
```json
{ "currentStreak": 3, "bestStreak": 10 }
```
- Errors: 400, 401, 500

### GET /admin/personalSummary/:classId
- Handler: admin_controller.PersonalSummary
- Response 202
```json
{
  "quick_summary": {
    "today_status": "present",
    "current_week_present": 3,
    "current_week_absent": 1,
    "current_week_not_marked": 0,
    "total_present": 20,
    "total_absent": 5,
    "total_not_marked": 2
  }
}
```
- Errors: 400, 401, 500

### GET /admin/report/:classId
- Handler: admin_controller.Report
- Response 200
```json
{
  "class_report": {
    "current_month": { "present": 15, "absent": 4, "not_marked": 1 },
    "current_year": { "present": 120, "absent": 20, "not_marked": 5 }
  }
}
```
- Errors: 400, 500

### GET /admin/personalReport/:classId
- Handler: admin_controller.PersonalReport
- Response 200
```json
{
  "personal_report": {
    "current_month": { "present": 10, "absent": 2, "not_marked": 0 },
    "current_year": { "present": 80, "absent": 10, "not_marked": 3 }
  }
}
```
- Errors: 400, 401, 500

### POST /admin/kickStudent/:classId
- Middlewares: AuthAdminMiddleware, IsAdminClass, IsUserEnrolledInClass
- Handler: admin_controller.KickStudent
- Body
```json
{
  "studentId": 0
}
```
- Response 200
```json
{ "message": "Student kicked from class successfully" }
```
- Errors: 400, 500

### POST /admin/banStudent/:classId
- Middlewares: AuthAdminMiddleware, IsAdminClass, IsUserEnrolledInClass
- Handler: admin_controller.BanStudent
- Body
```json
{
  "studentId": 0
}
```
- Response 200
```json
{ "message": "Student banned from class successfully" }
```
- Errors: 400, 500

## User Routes (/user)

Public

### POST /user/signIn
- Handler: user_controller.SignIn
- Body
```json
{
  "userName": "string",
  "password": "string"
}
```
- Sets cookie: refresh_token (HttpOnly)
- Response 200
```json
{
  "message": "Sign in successful",
  "role": "user",
  "access_token": "<jwt>",
  "user": {
    "id": 1,
    "username": "user1",
    "email": "u@example.com",
    "firstName": "Uma",
    "lastName": "User",
    "phone": "9999999999"
  }
}
```
- Errors: 400, 401

### POST /user/signUp
- Handler: user_controller.SignUp
- Body
```json
{
  "userName": "string",
  "password": "string",
  "email": "string",
  "firstName": "string",
  "lastName": "string",
  "phone": "string (digits)",
  "dob": "YYYY-MM-DD",
  "otp": "string"
}
```
- Response 201
```json
{
  "message": "User created successfully",
  "user": {
    "username": "user1",
    "email": "u@example.com",
    "firstName": "Uma",
    "lastName": "User",
    "phone": "9999999999"
  }
}
```
- Errors: 400, 401, 409

### POST /user/sendOTP
- Handler: user_controller.SendOTP
- Body
```json
{
  "phone": "string (digits)"
}
```
- Response 200
```json
{
  "message": "OTP sent",
  "phone": "9999999999"
}
```
- Errors: 400, 500

### POST /user/verifyOTP
- Handler: user_controller.VerifyOTP
- Body
```json
{
  "phone": "string (digits)",
  "otp": "string"
}
```
- Response 200
```json
{
  "message": "OTP verified successfully"
}
```
- Errors: 400 (expired/bad), 401 (wrong), 500

### POST /user/refreshToken
- Handler: user_controller.RefreshTokenUser
- Cookie: refresh_token
- Response 200
```json
{
  "access_token": "<jwt>"
}
```
- Errors: 401, 500

Protected (AuthUserMiddleware + IsUserClass)

### POST /user/markAttendance/:classID
- Handler: user_controller.MarkAttendance
- Body
```json
{
  "status": "present | absent | unmarked"
}
```
- Response 200
```json
{ "message": "Attendance marked", "class_id": 1 }
```
- Errors: 400, 401, 409, 500

### GET /user/classDetails/:classID
- Handler: user_controller.ClassDetails
- Response 200
```json
{
  "class": {
    "ID": 1,
    "Name": "Class A",
    "Email": "",
    "Phone": "",
    "CreatedByAdminId": 1,
    "ClassCode": "ABC123",
    "CreatedAt": "2024-01-01T00:00:00Z",
    "UpdatedAt": "2024-01-01T00:00:00Z"
  }
}
```
- Errors: 400, 404, 500

### GET /user/calendar/:classID
- Handler: user_controller.Calendar
- Response 200
```json
{
  "class_id": 1,
  "user_id": 1,
  "calendar": [
    {"date": "2024-01-01", "status": "present"}
  ]
}
```
- Errors: 400, 500

### GET /user/streak/:classID
- Handler: user_controller.Streak
- Response 200
```json
{ "currentStreak": 3, "bestStreak": 10 }
```
- Errors: 400, 401, 500

### GET /user/quickSummary/:classID
- Handler: user_controller.QuickSummary
- Response 200
```json
{
  "quick_summary": {
    "today_status": "present",
    "current_week_present": 3,
    "current_week_absent": 1,
    "current_week_not_marked": 0,
    "total_present": 20,
    "total_absent": 5,
    "total_not_marked": 2
  }
}
```
- Errors: 400, 401, 500

### GET /user/report/:classID
- Handler: user_controller.Report
- Response 200
```json
{
  "report": {
    "current_month": { "present": 10, "absent": 2, "not_marked": 0 },
    "current_year": { "present": 80, "absent": 10, "not_marked": 3 }
  }
}
```
- Errors: 400, 401, 500

Protected (AuthUserMiddleware)

### POST /user/enroll/:classCode
- Middlewares: AuthUserMiddleware, IsAllowedToEnroll
- Handler: user_controller.Enroll
- Path Param: classCode (string) [mapped to context classId via middleware]
- Response 200
```json
{
  "message": "User enrolled",
  "user_id": 1,
  "class_id": 1
}
```
- Errors: 400, 401, 500

### GET /user/classList
- Handler: user_controller.ClassList
- Response 200
```json
{
  "classes": [
    {
      "class_id": 1,
      "class_name": "Class A",
      "class_code": "ABC123",
      "created_at": "2024-01-01T00:00:00Z",
      "joined_at": "2024-01-02T00:00:00Z",
      "email": "",
      "phone": "",
      "created_by_admin_id": 1
    }
  ]
}
```
- Errors: 401, 500

### POST /user/logOutUser
- Handler: user_controller.LogOutUser
- Cookie: refresh_token
- Response 200
```json
{
  "message": "Logged out successfully"
}
```
- Errors: 400, 500

### PATCH /user/profile/:id
- Handler: user_controller.UpdateProfile
- Body
```json
{
  "firstName": "string",
  "lastName": "string",
  "email": "string"
}
```
- Response 200
```json
{
  "message": "Profile updated",
  "user_id": 1,
  "name": "Uma User",
  "email": "u@example.com"
}
```
- Errors: 400, 500

### GET /user/profile
- Handler: user_controller.Profile
- Response 200
```json
{
  "id": 1,
  "name": "Uma User",
  "email": "u@example.com"
}
```
- Errors: 401, 404, 500

### GET /user/resetPassword
- Handler: user_controller.ResetPassword
- Body
```json
{
  "newPassword": "string",
  "otp": 0,
  "phone": 0
}
```
- Response 200
```json
{
  "message": "Password reset successful"
}
```
- Errors: 400, 401, 500
