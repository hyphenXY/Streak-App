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
- Responses: 200 OK

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
- Responses: 201 Created

### GET /root/health-check
- Handler: root_controller.HealthCheck
- Responses: 200 OK

Protected (AuthRootMiddleware)

### GET /root/homepage/:id
- Handler: root_controller.Homepage
- Path Params: id (string)
- Responses: 200 OK

### GET /root/profile/:id
- Handler: root_controller.Profile
- Path Params: id (string)
- Responses: 200 OK

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
- Responses: 200 OK

### DELETE /root/admin/:id
- Handler: root_controller.DeleteAdmin
- Path Params: id (string)
- Responses: 200 OK

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
- Responses: 200 OK | 401 Unauthorized | 400 Bad Request

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
- Responses: 201 Created | 400 Bad Request | 401 Unauthorized | 409 Conflict

### POST /admin/sendOTP
- Handler: admin_controller.SendOTP
- Body
```json
{
  "phone": "string (digits)"
}
```
- Responses: 200 OK | 400 Bad Request | 500 Internal Server Error

### POST /admin/verifyOTP
- Handler: admin_controller.VerifyOTP
- Body
```json
{
  "phone": "string (digits)",
  "otp": "string"
}
```
- Responses: 200 OK | 400 Bad Request | 401 Unauthorized | 500 Internal Server Error

### POST /admin/refreshToken
- Handler: admin_controller.RefreshTokenUser
- Cookie: refresh_token
- Responses: 200 OK | 401 Unauthorized | 500 Internal Server Error

Protected (AuthAdminMiddleware)

### GET /admin/classList
- Handler: admin_controller.ClassList
- Responses: 200 OK

### GET /admin/profile
- Handler: admin_controller.Profile
- Responses: 200 OK | 401 Unauthorized

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
- Responses: 200 OK | 400 Bad Request | 500 Internal Server Error

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
- Responses: 201 Created | 400 Bad Request | 500 Internal Server Error

### POST /admin/logOutAdmin
- Handler: admin_controller.LogOutAdmin
- Cookie: refresh_token
- Responses: 200 OK | 400 Bad Request | 500 Internal Server Error

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
- Responses: 200 OK | 400 Bad Request | 401 Unauthorized | 500 Internal Server Error

Protected (AuthAdminMiddleware + IsAdminClass)

### GET /admin/quickSummary/:classId
- Handler: admin_controller.QuickSummary
- Responses: 200 OK | 400 Bad Request | 500 Internal Server Error

### GET /admin/todaySummary/:classId
- Handler: admin_controller.TodaySummary
- Responses: 200 OK | 400 Bad Request | 500 Internal Server Error

### GET /admin/calendar/:classId
- Handler: admin_controller.Calendar
- Responses: 200 OK | 400 Bad Request | 500 Internal Server Error

### POST /admin/markAttendance/:classId
- Handler: admin_controller.MarkAttendance
- Responses: 200 OK | 400 Bad Request | 401 Unauthorized | 403 Forbidden | 404 Not Found | 409 Conflict | 500 Internal Server Error

### GET /admin/studentsList/:classId
- Handler: admin_controller.StudentsList
- Responses: 200 OK | 400 Bad Request | 401 Unauthorized | 403 Forbidden | 500 Internal Server Error

### GET /admin/streak/:classId
- Handler: admin_controller.Streak
- Responses: 200 OK | 400 Bad Request | 401 Unauthorized | 500 Internal Server Error

### GET /admin/personalSummary/:classId
- Handler: admin_controller.PersonalSummary
- Responses: 202 Accepted | 400 Bad Request | 401 Unauthorized | 500 Internal Server Error

### GET /admin/report/:classId
- Handler: admin_controller.Report
- Responses: 200 OK | 400 Bad Request | 500 Internal Server Error

### GET /admin/personalReport/:classId
- Handler: admin_controller.PersonalReport
- Responses: 200 OK | 400 Bad Request | 401 Unauthorized | 500 Internal Server Error

### POST /admin/kickStudent/:classId
- Middlewares: AuthAdminMiddleware, IsAdminClass, IsUserEnrolledInClass
- Handler: admin_controller.KickStudent
- Body
```json
{
  "studentId": 0
}
```
- Responses: 200 OK | 400 Bad Request | 500 Internal Server Error

### POST /admin/banStudent/:classId
- Middlewares: AuthAdminMiddleware, IsAdminClass, IsUserEnrolledInClass
- Handler: admin_controller.BanStudent
- Body
```json
{
  "studentId": 0
}
```
- Responses: 200 OK | 400 Bad Request | 500 Internal Server Error

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
- Responses: 200 OK | 401 Unauthorized | 400 Bad Request

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
- Responses: 201 Created | 400 Bad Request | 401 Unauthorized | 409 Conflict

### POST /user/sendOTP
- Handler: user_controller.SendOTP
- Body
```json
{
  "phone": "string (digits)"
}
```
- Responses: 200 OK | 400 Bad Request | 500 Internal Server Error

### POST /user/verifyOTP
- Handler: user_controller.VerifyOTP
- Body
```json
{
  "phone": "string (digits)",
  "otp": "string"
}
```
- Responses: 200 OK | 400 Bad Request | 401 Unauthorized | 500 Internal Server Error

### POST /user/refreshToken
- Handler: user_controller.RefreshTokenUser
- Cookie: refresh_token
- Responses: 200 OK | 401 Unauthorized | 500 Internal Server Error

Protected (AuthUserMiddleware + IsUserClass)

### POST /user/markAttendance/:classID
- Handler: user_controller.MarkAttendance
- Body
```json
{
  "status": "string"
}
```
- Responses: 200 OK | 400 Bad Request | 401 Unauthorized | 409 Conflict | 500 Internal Server Error

### GET /user/classDetails/:classID
- Handler: user_controller.ClassDetails
- Responses: 200 OK | 400 Bad Request | 404 Not Found | 500 Internal Server Error

### GET /user/calendar/:classID
- Handler: user_controller.Calendar
- Responses: 200 OK | 400 Bad Request | 500 Internal Server Error

### GET /user/streak/:classID
- Handler: user_controller.Streak
- Responses: 200 OK | 400 Bad Request | 401 Unauthorized | 500 Internal Server Error

### GET /user/quickSummary/:classID
- Handler: user_controller.QuickSummary
- Responses: 200 OK | 400 Bad Request | 401 Unauthorized | 500 Internal Server Error

### GET /user/report/:classID
- Handler: user_controller.Report
- Responses: 200 OK | 400 Bad Request | 401 Unauthorized | 500 Internal Server Error

Protected (AuthUserMiddleware)

### POST /user/enroll/:classCode
- Middlewares: AuthUserMiddleware, IsAllowedToEnroll
- Handler: user_controller.Enroll
- Path Param: classCode (string) [mapped to context classId via middleware]
- Responses: 200 OK | 400 Bad Request | 401 Unauthorized | 500 Internal Server Error

### GET /user/classList
- Handler: user_controller.ClassList
- Responses: 200 OK | 401 Unauthorized | 500 Internal Server Error

### POST /user/logOutUser
- Handler: user_controller.LogOutUser
- Cookie: refresh_token
- Responses: 200 OK | 400 Bad Request | 500 Internal Server Error

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
- Responses: 200 OK | 400 Bad Request | 500 Internal Server Error

### GET /user/profile
- Handler: user_controller.Profile
- Responses: 200 OK | 401 Unauthorized | 404 Not Found | 500 Internal Server Error

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
- Responses: 200 OK | 400 Bad Request | 401 Unauthorized | 500 Internal Server Error
