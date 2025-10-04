package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"jojihouse-system/internal/model"
	"strings"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAllUsers() ([]model.User, error) {
	var users []model.User
	err := r.db.Select(&users, "SELECT * FROM users ORDER BY number")

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

func (r *UserRepository) GetUserByID(id int) (*model.User, error) {
	user := &model.User{}
	err := r.db.Get(user, "SELECT * FROM users WHERE id = $1", id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetUserByBarcode(barcode string) (*model.User, error) {
	user := &model.User{}
	err := r.db.Get(user, "SELECT * FROM users WHERE barcode = $1", barcode)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetUserByNumber(number int) (*model.User, error) {
	user := &model.User{}
	err := r.db.Get(user, "SELECT * FROM users WHERE number = $1", number)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return user, nil
}

func (r *UserRepository) CreateUser(user *model.User) (*model.User, error) {
	// sql文の構築
	columns := []string{}
	values := []string{}
	args := map[string]interface{}{}

	if user.Name != nil {
		columns = append(columns, "name")
		values = append(values, ":name")
		args["name"] = *user.Name
	}
	if user.Description != nil {
		columns = append(columns, "description")
		values = append(values, ":description")
		args["description"] = *user.Description
	}
	if user.Barcode != nil {
		columns = append(columns, "barcode")
		values = append(values, ":barcode")
		args["barcode"] = *user.Barcode
	}
	if user.Contact != nil {
		columns = append(columns, "contact")
		values = append(values, ":contact")
		args["contact"] = *user.Contact
	}
	if user.Remaining_entries != nil {
		columns = append(columns, "remaining_entries")
		values = append(values, ":remaining_entries")
		args["remaining_entries"] = *user.Remaining_entries
	}
	if user.Allergy != nil {
		columns = append(columns, "allergy")
		values = append(values, ":allergy")
		args["allergy"] = *user.Allergy
	}
	if user.Number != nil {
		columns = append(columns, "number")
		values = append(values, ":number")
		args["number"] = *user.Number
	}

	query := fmt.Sprintf(`
		INSERT INTO users (%s)
		VALUES (%s)
		RETURNING *`,
		strings.Join(columns, ", "),
		strings.Join(values, ", "),
	)

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// ここがポイント！！ QueryRowxではなく NamedQuery → その結果から取り出す
	rows, err := stmt.Queryx(args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var insertedUser model.User
	if rows.Next() {
		if err := rows.StructScan(&insertedUser); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("no user inserted")
	}

	return &insertedUser, nil

}

func (r *UserRepository) UpdateUser(user *model.User) error {
	_, err := r.db.NamedExec(`
		UPDATE users SET
			name = :name,
			description = :description,
			barcode = :barcode,
			contact = :contact,
			remaining_entries = :remaining_entries,
			total_entries = :total_entries,
			allergy = :allergy,
			number = :number
		WHERE id = :id
	`, user)
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

	query, args, err := sqlx.In("SELECT * FROM users WHERE id IN (?)", userIDs)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)

	var users []model.User
	err = r.db.Select(&users, query, args...)
	if err != nil {
		return nil, err
	}
	return users, nil

}
