package storage

import (
	"errors"
	"fmt"
	"sync"

	"dz_go/models"
)

type InMemoryStorage struct {
	users      map[int]*models.User
	mutex      sync.RWMutex
	nextUserID int
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		users:      make(map[int]*models.User),
		nextUserID: 1,
	}
}

func (s *InMemoryStorage) CreateUser(user *models.User) (int, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user.ID = s.nextUserID
	s.users[user.ID] = user
	s.nextUserID++

	return user.ID, nil
}

func (s *InMemoryStorage) GetUserByID(id int) (*models.User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *InMemoryStorage) UpdateUser(updatedUser *models.User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.users[updatedUser.ID]; !exists {
		return errors.New("user not found")
	}

	s.users[updatedUser.ID] = updatedUser

	return nil
}

func (s *InMemoryStorage) DeleteUser(id int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.users[id]; !exists {
		return errors.New("user not found")
	}

	// Remove user from friends' lists
	for _, user := range s.users {
		for i, friendID := range user.Friends {
			if friendID == id {
				user.Friends = append(user.Friends[:i], user.Friends[i+1:]...)
				break
			}
		}
	}

	delete(s.users, id)

	return nil
}

func (s *InMemoryStorage) MakeFriends(sourceID, targetID int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	sourceUser, sourceExists := s.users[sourceID]
	targetUser, targetExists := s.users[targetID]

	if !sourceExists || !targetExists {
		return errors.New("one or both users not found")
	}

	// Check if they are already friends
	for _, friendID := range sourceUser.Friends {
		if friendID == targetID {
			return fmt.Errorf("%s and %s are already friends", sourceUser.Name, targetUser.Name)
		}
	}

	sourceUser.Friends = append(sourceUser.Friends, targetID)
	targetUser.Friends = append(targetUser.Friends, sourceID)

	return nil
}

func (s *InMemoryStorage) ListFriends(userID int) ([]int, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user.Friends, nil
}
