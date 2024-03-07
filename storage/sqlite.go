package storage

import (
	"database/sql"
	"dz_go/models"
	"errors"
	_ "modernc.org/sqlite"
)

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(dsn string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	s := &SQLiteStorage{db: db}

	err = s.createTables()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *SQLiteStorage) createTables() error {
	q := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			age INTEGER NOT NULL
		);
		CREATE TABLE IF NOT EXISTS friendships (
			user_id INTEGER,
			friend_id INTEGER,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (friend_id) REFERENCES users(id)
		);`
	if _, err := s.db.Exec(q); err != nil {
		return err
	}

	return nil
}

func (s *SQLiteStorage) CreateUser(user *models.User) (int, error) {
	result, err := s.db.Exec(`INSERT INTO users (name, age) VALUES (?, ?)`, user.Name, user.Age)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *SQLiteStorage) GetUserByID(id int) (*models.User, error) {
	var user models.User

	row := s.db.QueryRow("SELECT id, name, age FROM users WHERE id = ?", id)
	err := row.Scan(&user.ID, &user.Name, &user.Age)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	// Заполнение списка друзей пользователя
	rows, err := s.db.Query("SELECT friend_id FROM friendships WHERE user_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var friendID int
		if err := rows.Scan(&friendID); err != nil {
			return nil, err
		}
		user.Friends = append(user.Friends, friendID)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *SQLiteStorage) UpdateUser(user *models.User) error {
	result, err := s.db.Exec("UPDATE users SET age = ? WHERE id = ?", user.Age, user.ID)
	if err != nil {
		return err
	}

	// Проверка, что был обновлен хотя бы один ряд
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows affected, user not found or data is the same")
	}

	return nil
}

func (s *SQLiteStorage) DeleteUser(id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// Удаление пользователя из таблицы `users`
	_, err = tx.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		tx.Rollback() // Откат в случае ошибки

		return err
	}

	// Удаление всех записей о дружбе связанных с пользователем
	_, err = tx.Exec("DELETE FROM friendships WHERE user_id = ? OR friend_id = ?", id, id)
	if err != nil {
		tx.Rollback() // Откат в случае ошибки

		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *SQLiteStorage) MakeFriends(sourceID, targetID int) error {
	// Проверка, существуют ли оба пользователя
	if !s.userExists(sourceID) || !s.userExists(targetID) {
		return errors.New("one or both users not found")
	}

	// Начало транзакции
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// Проверка существования дружбы
	if s.friendshipExists(tx, sourceID, targetID) {
		tx.Rollback()

		return errors.New("the users are already friends")
	}

	// Создание записи о дружбе
	q := `INSERT INTO friendships (user_id, friend_id) VALUES (?, ?), (?, ?)`
	_, err = tx.Exec(q, sourceID, targetID, targetID, sourceID)
	if err != nil {
		tx.Rollback() // Откат в случае ошибки

		return err
	}

	// Подтверждение транзакции
	return tx.Commit()
}

// userExists проверяет существование пользователя в базе данных
func (s *SQLiteStorage) userExists(userID int) bool {
	var exists bool
	err := s.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`, userID).Scan(&exists)

	return err == nil && exists
}

// friendshipExists проверяет существование записи о дружбе между двумя пользователями
func (s *SQLiteStorage) friendshipExists(tx *sql.Tx, sourceID, targetID int) bool {
	var exists bool

	q := `SELECT EXISTS(SELECT 1 FROM friendships WHERE user_id = ? AND friend_id = ?)`
	err := tx.QueryRow(q, sourceID, targetID).Scan(&exists)

	return err == nil && exists
}

func (s *SQLiteStorage) ListFriends(userID int) ([]int, error) {
	var friendIDs []int

	rows, err := s.db.Query(`SELECT friend_id FROM friendships WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Итерация по результатам запроса и добавление ID каждого друга в список
	for rows.Next() {
		var friendID int
		if err := rows.Scan(&friendID); err != nil {
			return nil, err
		}
		friendIDs = append(friendIDs, friendID)
	}

	// Проверка на ошибки при итерации
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return friendIDs, nil
}
