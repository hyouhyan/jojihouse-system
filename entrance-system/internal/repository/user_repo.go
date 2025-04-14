package repository

import (
	"jojihouse-entrance-system/internal/model"

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
	err := r.db.Select(&users, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetUserByID(id int) (*model.User, error) {
	user := &model.User{}
	err := r.db.Get(user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByBarcode(barcode string) (*model.User, error) {
	user := &model.User{}
	err := r.db.Get(user, "SELECT * FROM users WHERE barcode = $1", barcode)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) CreateUser(user *model.User) (*model.User, error) {
	_, err := r.db.NamedExec(`
		INSERT INTO users (name, description, barcode, contact, remaining_entries, allergy)
		VALUES (:name, :description, :barcode, :contact, :remaining_entries, :allergy)
	`, user)
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
	_, err := r.db.NamedExec(`
		UPDATE users SET
			name = :name,
			description = :description,
			barcode = :barcode,
			contact = :contact,
			remaining_entries = :remaining_entries,
			total_entries = :total_entries,
			allergy = :allergy
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
