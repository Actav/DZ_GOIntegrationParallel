package storage

import "dz_go/models"

// Storage определяет интерфейс для хранилища данных пользователей
type Storage interface {
	CreateUser(user *models.User) (int, error)
	GetUserByID(id int) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
	MakeFriends(sourceID, targetID int) error
	ListFriends(userID int) ([]int, error)
}
