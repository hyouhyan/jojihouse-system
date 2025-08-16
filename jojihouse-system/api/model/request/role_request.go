package request

type AddRole struct {
	RoleID int `json:"role_id" binding:"required"`
}
