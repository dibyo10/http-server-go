package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var (
	userCache  = make(map[int]User)
	cacheMutex sync.RWMutex
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleRoot)

	mux.HandleFunc("POST /signup", signup)
	mux.HandleFunc("POST /signin", signin)
	mux.HandleFunc("PUT /users/{id}", updateUser)
	mux.HandleFunc("DELETE /users/{id}", deleteUser)

	http.ListenAndServe(":8080", mux)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server is running"))
}

func signup(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.Password == "" || user.Name == "" {
		http.Error(w, "name, email and password are required", http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	cacheMutex.Lock()
	newID := len(userCache) + 1
	user.ID = newID
	userCache[newID] = user
	cacheMutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("User " + user.Name + " signed up successfully")
}

func signin(w http.ResponseWriter, r *http.Request) {
	var input User

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if input.Email == "" || input.Password == "" {
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}

	cacheMutex.RLock()

	for _, user := range userCache {
		if user.Email != input.Email {
			continue
		}
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

		if err != nil {

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Invalid email or password")
			return

		}
		cacheMutex.RUnlock()
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("User " + user.Name + " signed in successfully")
		return

	}

	cacheMutex.RUnlock()

	http.Error(w, "invalid email or password", http.StatusUnauthorized)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cacheMutex.RLock()
	user, exists := userCache[id]
	cacheMutex.RUnlock()
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var input User
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		user.Password = string(hashed)
	}
	cacheMutex.Lock()
	userCache[id] = user
	cacheMutex.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("User " + user.Name + " updated successfully")

}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cacheMutex.Lock()
	user, exists := userCache[id]
	if !exists {
		cacheMutex.Unlock()
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	delete(userCache, id)
	cacheMutex.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("User " + user.Name + " deleted successfully")
}
