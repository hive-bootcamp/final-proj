package api

import (
    "net/http"
    "os"

    "github.com/golang-jwt/jwt/v4"
)

func validateToken(tokenStr string) bool {
    serverPass := os.Getenv("TODO_PASSWORD")
    if serverPass == "" {
        return true // пароль не задан — аутентификация отключена
    }

    token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil || !token.Valid {
        return false
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return false
    }

    // токен должен содержать hash == текущий пароль
    if claims["hash"] != serverPass {
        return false
    }

    return true
}

func auth(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        pass := os.Getenv("TODO_PASSWORD")

        if pass == "" {
            next(w, r)
            return
        }

        var token string
        cookie, err := r.Cookie("token")
        if err == nil {
            token = cookie.Value
        }

        if !validateToken(token) {
            http.Error(w, "Authentification required", http.StatusUnauthorized)
            return
        }

        next(w, r)
    })
}
