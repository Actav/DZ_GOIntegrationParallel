package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"dz_go/models"
	"dz_go/storage"
	"github.com/go-chi/chi/v5"
)

// Handlers структура для хранения зависимостей обработчиков
type Handlers struct {
	Storage storage.Storage
}

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	userID, err := h.Storage.CreateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": userID})
}

func (h *Handlers) MakeFriends(w http.ResponseWriter, r *http.Request) {
	var ids struct {
		SourceID int `json:"source_id"`
		TargetID int `json:"target_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	sourceUser, err := h.Storage.GetUserByID(ids.SourceID)
	if err != nil {
		http.Error(w, "Source user not found", http.StatusNotFound)

		return
	}

	targetUser, err := h.Storage.GetUserByID(ids.TargetID)
	if err != nil {
		http.Error(w, "Target user not found", http.StatusNotFound)

		return
	}

	err = h.Storage.MakeFriends(ids.SourceID, ids.TargetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s и %s теперь друзья", sourceUser.Name, targetUser.Name)
}

func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var u struct {
		ID int `json:"target_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if err := h.Storage.DeleteUser(u.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User %d has been deleted", u.ID)
}

func (h *Handlers) ListFriends(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)

		return
	}

	friends, err := h.Storage.ListFriends(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(friends)
}

func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
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

	user, err := h.Storage.GetUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	user.Age = data.NewAge
	err = h.Storage.UpdateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User %d's age has been updated to %d", userID, data.NewAge)
}

func (h *Handlers) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)

		return
	}

	user, err := h.Storage.GetUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
