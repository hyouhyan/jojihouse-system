package repository

import (
	"database/sql"
	"jojihouse-management-system/internal/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type CurrentUsersRepository struct {
	db *sqlx.DB
}

func NewCurrentUsersRepository(db *sqlx.DB) *CurrentUsersRepository {
	return &CurrentUsersRepository{db: db}
}

// 在室ユーザーに追加
func (r *CurrentUsersRepository) AddUserToCurrentUsers(userID int) error {
	_, err := r.db.Exec(`
		INSERT INTO current_users (user_id) VALUES ($1)
		ON CONFLICT (user_id) DO NOTHING;
	`, userID)
	if err != nil {
		return err
	}

	return nil
}

// 在室ユーザーから削除
func (r *CurrentUsersRepository) DeleteUserToCurrentUsers(userID int) error {
	_, err := r.db.Exec("DELETE FROM current_users WHERE user_id = $1;", userID)
	if err != nil {
		return err
	}

	return nil
}

// 在室ユーザー一覧を取得
func (r *CurrentUsersRepository) GetCurrentUsers() ([]model.CurrentUser, error) {
	var users []model.CurrentUser

	query := `
		SELECT u.id AS user_id, u.name, c.entered_at, u.allergy
		FROM current_users c
		JOIN users u ON c.user_id = u.id
	`

	err := r.db.Select(&users, query)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// 在室中のユーザーの入室時間を取得
func (r *CurrentUsersRepository) GetEnteredTime(userID int) (time.Time, error) {
	var enteredAt time.Time

	query := `
		SELECT entered_at
		FROM current_users
		WHERE user_id = $1
	`

	err := r.db.QueryRow(query, userID).Scan(&enteredAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return time.Time{}, nil
		}
		return time.Time{}, err
	}

	return enteredAt, nil
}
