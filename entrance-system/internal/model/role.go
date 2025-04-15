package model

type Role struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type UserRoles struct {
	UserID int `db:"user_id"`
	RoleID int `db:"role_id"`
}

const (
	RoleMember      = "member"
	RoleStudent     = "student"
	RoleSystemAdmin = "system-admin"
	RoleHouseAdmin  = "house-admin"
	RoleGuest       = "guest"
)
