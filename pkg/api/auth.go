package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("todo-secret-key")

type signinReq struct {
	Password string `json:"password"`
}

type signinResp struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func generateToken(pass string) (string, error) {
	claims := jwt.MapClaims{
		"hash": pass,
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// POST
func signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, signinResp{Error: "method not allowed"})
		return
	}

	var req signinReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, signinResp{Error: "invalid json"})
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		writeJSON(w, signinResp{Error: "server password not set"})
		return
	}

	if req.Password != pass {
		writeJSON(w, signinResp{Error: "Неверный пароль"})
		return
	}

	token, err := generateToken(pass)
	if err != nil {
		writeJSON(w, signinResp{Error: err.Error()})
		return
	}

	writeJSON(w, signinResp{Token: token})
}
