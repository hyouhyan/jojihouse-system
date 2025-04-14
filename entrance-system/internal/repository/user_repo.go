package repository

import (
	"database/sql"
	"fmt"

	"jojihouse-entrance-system/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAllUsers() ([]model.User, error) {
	rows, err := r.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Description,
			&user.Barcode,
			&user.Contact,
			&user.Remaining_entries,
			&user.Registered_at,
			&user.Total_entries,
			&user.Allergy,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) GetUserByID(id int) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(
		&user.ID,
		&user.Name,
		&user.Description,
		&user.Barcode,
		&user.Contact,
		&user.Remaining_entries,
		&user.Registered_at,
		&user.Total_entries,
		&user.Allergy,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByBarcode(barcode string) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow("SELECT * FROM users WHERE barcode = $1", barcode).Scan(
		&user.ID,
		&user.Name,
		&user.Description,
		&user.Barcode,
		&user.Contact,
		&user.Remaining_entries,
		&user.Registered_at,
		&user.Total_entries,
		&user.Allergy,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) CreateUser(user *model.User) (*model.User, error) {
	_, err := r.db.Exec(
		"INSERT INTO users (name, description, barcode, contact, remaining_entries, allergy) VALUES ($1, $2, $3, $4, $5, $6)",
		user.Name,
		user.Description,
		user.Barcode,
		user.Contact,
		user.Remaining_entries,
		user.Allergy,
	)
	if err != nil {
		return nil, err
	}

	// ユーザー情報を取得
	user, err = r.GetUserByBarcode(user.Barcode)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) UpdateUser(user *model.User) error {
	_, err := r.db.Exec(
		"UPDATE users SET name = $1, description = $2, barcode = $3, contact = $4, remaining_entries = $5, total_entries = $6, allergy = $7 WHERE id = $8",
		user.Name,
		user.Description,
		user.Barcode,
		user.Contact,
		user.Remaining_entries,
		user.Total_entries,
		user.Allergy,
		user.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) DeleteUser(id int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

// 入場可能回数を減らす
func (r *UserRepository) DecreaseRemainingEntries(id int, count int) (int, int, error) {
	var before int
	var after int

	// remaining_entries を更新しつつ、更新前後の値を取得
	err := r.db.QueryRow(`
		UPDATE users 
		SET remaining_entries = remaining_entries - $1 
		WHERE id = $2
		RETURNING remaining_entries + $1, remaining_entries
	`, count, id).Scan(&before, &after)

	if err != nil {
		return 0, 0, err
	}
	return before, after, nil
}

// 入場可能回数を増やす
func (r *UserRepository) IncreaseRemainingEntries(id int, count int) (int, int, error) {
	var before int
	var after int

	// remaining_entries を更新しつつ、更新前後の値を取得
	err := r.db.QueryRow(`
		UPDATE users
		SET remaining_entries = remaining_entries + $1
		WHERE id = $2
		RETURNING remaining_entries - $1, remaining_entries
	`, count, id).Scan(&before, &after)

	if err != nil {
		return 0, 0, err
	}

	return before, after, nil
}

// 総入場回数を増やす
func (r *UserRepository) IncreaseTotalEntries(id int) error {
	_, err := r.db.Exec("UPDATE users SET total_entries = total_entries + 1 WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

// 複数の UserID から User を取得
func (r *UserRepository) GetUsersByIDs(userIDs []int) ([]model.User, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE id IN (%s)", intArrayToString(userIDs))
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Description,
			&user.Barcode,
			&user.Contact,
			&user.Remaining_entries,
			&user.Registered_at,
			&user.Total_entries,
			&user.Allergy,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// int のスライスをカンマ区切りの文字列に変換
func intArrayToString(arr []int) string {
	result := ""
	for i, val := range arr {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("%d", val)
	}
	return result
}
