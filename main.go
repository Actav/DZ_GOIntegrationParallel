package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
)

type User struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Friends []int  `json:"friends"` // Список идентификаторов друзей
}

var (
	users      = make(map[int]*User) // Использование int в качестве типа ключа
	usersMutex = sync.RWMutex{}
	nextUserID = 1 // Следующий ID для нового пользователя
)

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usersMutex.Lock()
	userID := nextUserID
	users[userID] = &user
	nextUserID++
	usersMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": userID})
}

func makeFriends(w http.ResponseWriter, r *http.Request) {
	var ids struct {
		SourceID int `json:"source_id"`
		TargetID int `json:"target_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usersMutex.Lock()
	defer usersMutex.Unlock()

	sourceUser, sourceExists := users[ids.SourceID]
	targetUser, targetExists := users[ids.TargetID]

	if !sourceExists || !targetExists {
		http.Error(w, "One or both users not found", http.StatusBadRequest)
		return
	}

	sourceUser.Friends = append(sourceUser.Friends, ids.TargetID)
	targetUser.Friends = append(targetUser.Friends, ids.SourceID)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s и %s теперь друзья", sourceUser.Name, targetUser.Name)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	var id struct {
		TargetID int `json:"target_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usersMutex.Lock()
	defer usersMutex.Unlock()

	user, exists := users[id.TargetID]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	delete(users, id.TargetID)

	for _, friendID := range user.Friends {
		friend := users[friendID]
		for i, fid := range friend.Friends {
			if fid == id.TargetID {
				friend.Friends = append(friend.Friends[:i], friend.Friends[i+1:]...)
				break
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s удален", user.Name)
}

func listFriends(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	usersMutex.RLock()
	defer usersMutex.RUnlock()

	user, exists := users[userID]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user.Friends)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var data struct {
		NewAge int `json:"new_age"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usersMutex.Lock()
	defer usersMutex.Unlock()

	user, exists := users[userID]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user.Age = data.NewAge
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "возраст пользователя успешно обновлён")
}

func main() {
	r := chi.NewRouter()

	r.Post("/create", createUser)
	r.Post("/make_friends", makeFriends)
	r.Delete("/user", deleteUser)
	r.Get("/friends/{user_id}", listFriends)
	r.Put("/{user_id}", updateUser)

	log.Println("Start server")
	http.ListenAndServe(":8080", r)
}
