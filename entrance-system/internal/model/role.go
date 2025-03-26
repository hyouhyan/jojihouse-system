package model

type Role struct {
	ID   int
	Name string
}

type UserRoles struct {
	UserID int
	RoleID int
}

const (
	RoleMember      = "member"
	RoleStudent     = "student"
	RoleSystemAdmin = "system-admin"
	RoleHouseAdmin  = "house-admin"
	RoleGuest       = "guest"
)
