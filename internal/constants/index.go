package constants

var UserEnrollment = struct {
	Enrolled  int8
	Kicked    int8
	Banned    int8
	Suspended int8
}{
	Enrolled:  1,
	Kicked:    2,
	Banned:    3,
	Suspended: 4,
}

var AttendanceStatus = struct {
	Present  int8
	Absent   int8
	UnMarked int8
}{
	Present:  1,
	Absent:   0,
	UnMarked: 2,
}
