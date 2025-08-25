package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	User    string `json:"user,omitempty"`
}

var users = map[string]string{
	"admin":    "123456",
	"dole2711": "123456",
}

var currentUser string

func login(username, password string) error {
	if pass, ok := users[username]; ok {
		if pass == password {
			currentUser = username
			return nil
		}
		return errors.New("sai mật khẩu")
	}
	return errors.New("không tìm thấy user")
}

func logout() {
	if currentUser != "" {
		fmt.Printf("User %s đã logout\n", currentUser)
		currentUser = ""
	} else {
		fmt.Println("Chưa có user nào đang login")
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = login(req.Username, req.Password)
	if err != nil {
		json.NewEncoder(w).Encode(LoginResponse{Status: "fail", Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(LoginResponse{Status: "success", User: currentUser})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	if currentUser == "" {
		json.NewEncoder(w).Encode(LoginResponse{Status: "fail", Message: "chưa có user nào login"})
		return
	}

	logout()
	json.NewEncoder(w).Encode(LoginResponse{Status: "success", Message: "đã logout"})
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	fmt.Println("Auth API running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
