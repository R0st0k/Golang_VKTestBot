package repository

// Структура таблицы 'password'
type Password struct {
	ID       int64  `db:"id"`
	UserID   string `db:"user_id"`
	Resource string `db:"resource"`
	Login    string `db:"login"`
	Password string `db:"password"`
}

// Интерфейс для реализации хранилища данных
type Repository interface {
	Create(password Password) error
	AllUserPassword(userID string) ([]Password, error)
	GetByUserAndResource(userID, resource string) (*Password, error)
	Update(updated Password) (*Password, error)
	Delete(userID, resource string) error
}
