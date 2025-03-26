package model

type Role struct {
	ID   int
	Role string
}

const (
	RoleMember      = "member"
	RoleStudent     = "student"
	RoleSystemAdmin = "system-admin"
	RoleHouseAdmin  = "house-admin"
	RoleGuest       = "guest"
)
