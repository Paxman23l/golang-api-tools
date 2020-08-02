package enums

// RoleType enums access the possible role types for the system
type RoleType string

const (
	// SuperAdmin is a systemadmin
	SuperAdmin = "superadmin"
	// SchoolAdmin is a school admin
	SchoolAdmin = "schooladmin"
	// Teacher is a teacher
	Teacher = "teacher"
	// Parent is a parent
	Parent = "parent"
	// Student is a student
	Student = "student"
	// User should have no rights
	User = "user"
)
