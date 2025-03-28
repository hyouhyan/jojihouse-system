package repository

import (
	"database/sql"

	"jojihouse-entrance-system/internal/model"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) GetRoleByID(id int) (*model.Role, error) {
	role := &model.Role{}
	err := r.db.QueryRow("SELECT * FROM roles WHERE id = $1", id).Scan(
		&role.ID,
		&role.Name,
	)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) GetRoleByName(name string) (*model.Role, error) {
	role := &model.Role{}
	err := r.db.QueryRow("SELECT * FROM roles WHERE name = $1", name).Scan(
		&role.ID,
		&role.Name,
	)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) GetRolesByUserID(userID int) ([]model.Role, error) {
	roles := []model.Role{}
	rows, err := r.db.Query("SELECT r.id, r.name FROM roles r JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		role := model.Role{}
		err := rows.Scan(
			&role.ID,
			&role.Name,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
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
