package repository

import (
	"jojihouse-management-system/internal/model"

	"github.com/jmoiron/sqlx"
)

type RoleRepository struct {
	db *sqlx.DB
}

func NewRoleRepository(db *sqlx.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) GetAllRoles() ([]model.Role, error) {
	roles := []model.Role{}
	err := r.db.Select(&roles, "SELECT * FROM roles")
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleRepository) GetRoleByID(id int) (*model.Role, error) {
	role := &model.Role{}
	err := r.db.Get(role, "SELECT * FROM roles WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) GetRoleByName(name string) (*model.Role, error) {
	role := &model.Role{}
	err := r.db.Get(role, "SELECT * FROM roles WHERE name = $1", name)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) GetRolesByUserID(userID int) ([]model.Role, error) {
	roles := []model.Role{}
	err := r.db.Select(&roles, "SELECT r.* FROM roles r JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleRepository) IsMember(userID int) (bool, error) {
	roles, err := r.GetRolesByUserID(userID)
	if err != nil {
		return false, err
	}
	for _, role := range roles {
		if role.Name == model.RoleMember {
			return true, nil
		}
	}
	return false, nil
}

func (r *RoleRepository) IsStudent(userID int) (bool, error) {
	roles, err := r.GetRolesByUserID(userID)
	if err != nil {
		return false, err
	}
	for _, role := range roles {
		if role.Name == model.RoleStudent {
			return true, nil
		}
	}
	return false, nil
}

func (r *RoleRepository) IsHouseAdmin(userID int) (bool, error) {
	roles, err := r.GetRolesByUserID(userID)
	if err != nil {
		return false, err
	}
	for _, role := range roles {
		if role.Name == model.RoleHouseAdmin {
			return true, nil
		}
	}
	return false, nil
}

func (r *RoleRepository) IsSystemAdmin(userID int) (bool, error) {
	roles, err := r.GetRolesByUserID(userID)
	if err != nil {
		return false, err
	}
	for _, role := range roles {
		if role.Name == model.RoleSystemAdmin {
			return true, nil
		}
	}
	return false, nil
}

func (r *RoleRepository) IsGuest(userID int) (bool, error) {
	roles, err := r.GetRolesByUserID(userID)
	if err != nil {
		return false, err
	}
	for _, role := range roles {
		if role.Name == model.RoleGuest {
			return true, nil
		}
	}
	return false, nil
}

func (r *RoleRepository) AddRoleToUser(userID, roleID int) error {
	_, err := r.db.Exec("INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)", userID, roleID)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoleRepository) RemoveRoleFromUser(userID, roleID int) error {
	_, err := r.db.Exec("DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2", userID, roleID)
	if err != nil {
		return err
	}
	return nil
}
