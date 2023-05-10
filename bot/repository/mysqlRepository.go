package repository

import (
	"VKTestBot/configs"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

var (
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

var db *sql.DB

type MySQLRepository struct {
	db *sql.DB
}

func NewSQLiteRepository() *MySQLRepository {
	return &MySQLRepository{
		db: db,
	}
}

// Инициализация MySQL базы
func Init() {
	addr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		configs.Getenv("SQL_USER", "root"),
		configs.Getenv("SQL_PASSWORD", "password"),
		configs.Getenv("SQL_HOST", "127.0.0.1"),
		configs.Getenv("SQL_PORT", "3306"),
		configs.Getenv("SQL_DATABASE", "bot"),
	)

	var err error
	db, err = sql.Open("mysql", addr)

	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	log.Printf("MySQL Connected!")
}

// Добавление записи в таблицу
func (r *MySQLRepository) Create(password Password) error {
	resp, err := r.GetByUserAndResource(password.UserID, password.Resource)
	if err != nil {
		_, err = r.db.Exec("INSERT INTO password(user_id, resource, login, password) VALUES(?,?,?,?)",
			password.UserID, password.Resource, password.Login, password.Password)

		if err != nil {
			return err
		}
	} else {
		if resp.Login != password.Login || resp.Password != password.Password {
			_, err = r.Update(password)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Получение всех паролей пользователя
func (r *MySQLRepository) AllUserPassword(userID string) ([]Password, error) {
	rows, err := r.db.Query("SELECT * FROM password WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Password
	for rows.Next() {
		var password Password
		if err := rows.Scan(&password.ID, &password.UserID, &password.Resource, &password.Login, &password.Password); err != nil {
			return nil, err
		}
		all = append(all, password)
	}
	return all, nil
}

// Получение пароля по пользователю и ресурсу
func (r *MySQLRepository) GetByUserAndResource(userID, resource string) (*Password, error) {
	row := r.db.QueryRow("SELECT * FROM password WHERE user_id = ? AND resource = ?", userID, resource)

	var password Password
	if err := row.Scan(&password.ID, &password.UserID, &password.Resource, &password.Login, &password.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &password, nil
}

// Обновление пароля
func (r *MySQLRepository) Update(updated Password) (*Password, error) {
	res, reqErr := r.db.Exec("UPDATE password SET login = ?, password = ? WHERE user_id = ? AND resource = ?",
		updated.Login, updated.Password, updated.UserID, updated.Resource)

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}

	if reqErr != nil {
		return nil, err
	}

	return &updated, nil
}

// Удаление пароля
func (r *MySQLRepository) Delete(userID, resource string) error {
	res, err := r.db.Exec("DELETE FROM password WHERE user_id = ? AND resource = ?", userID, resource)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return err
}
