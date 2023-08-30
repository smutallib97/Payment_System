package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type UserModel struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

var Users = []UserModel{}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "Application/json")
		err := json.NewEncoder(w).Encode(Users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	} else {
		w.WriteHeader(http.StatusBadGateway)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var newUser UserModel
		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newUser.Id = len(Users) + 1
		newUser.CreatedAt = time.Now().Format("2006-01-02T15:04:05")
		newUser.UpdatedAt = time.Now().Format("2006-01-02T15:04:05")
		Users = append(Users, newUser)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "User created successfully!")
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		var updateUser UserModel
		err := json.NewDecoder(r.Body).Decode(&updateUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for i, user := range Users {
			if user.Id == updateUser.Id {
				updateUser.CreatedAt = Users[i].CreatedAt
				updateUser.UpdatedAt = time.Now().Format("2006-01-02T15:04:05")
				Users[i] = updateUser
				fmt.Fprintln(w, updateUser)
				w.WriteHeader(http.StatusOK)
			}
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		userIdStr := r.URL.Query().Get("id")
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		for i, user := range Users {
			if user.Id == userId {
				Users = append(Users[:i], Users[i+1:]...)
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, "User deleted successfully!")
			}
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
func getUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		userIDStr := r.URL.Query().Get("id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user Id", http.StatusNotFound)
			return
		}
		var foundUser *UserModel
		for _, user := range Users {
			if user.Id == userID {
				foundUser = &user
				w.Header().Set("Content-Type", "application/json")
				err := json.NewEncoder(w).Encode(user)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				break
			}
		}
		if foundUser != nil {
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "User not found", http.StatusNotFound)
		}
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {

	http.HandleFunc("/getUser", index)
	http.HandleFunc("/addUser", createUser)
	http.HandleFunc("/updateUser", updateUser)
	http.HandleFunc("/deleteUser", deleteUser)
	http.HandleFunc("/getUserById", getUserByID)
	log.Println("Listening...!!!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
