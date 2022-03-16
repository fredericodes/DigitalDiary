package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const (
	missingAuthHeaderErr = "missing Authorization header"
	verifyJwtTokenErr    = "error verifying JWT token:"
)

func Authenticate(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("{error: %v}", missingAuthHeaderErr)))
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		claims, err := verifyToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("{%v %v}", verifyJwtTokenErr, err.Error())))
			return
		}

		id := claims.(jwt.MapClaims)["id"].(string)
		iat := claims.(jwt.MapClaims)["iat"].(float64)
		exp := claims.(jwt.MapClaims)["exp"].(float64)

		r.Header.Set("id", id)
		r.Header.Set("iat", fmt.Sprintf("%v", iat))
		r.Header.Set("exp", fmt.Sprintf("%v", exp))

		handler.ServeHTTP(w, r)
	}
}

func GenerateToken(id uuid.UUID) (string, error) {
	signingKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   id.String(),
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Minute * 10).Unix(),
	})

	tokenString, err := token.SignedString(signingKey)
	return tokenString, err
}

func verifyToken(tokenString string) (jwt.Claims, error) {
	signingKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token.Claims, err
}

func ReadUserId(r *http.Request) (*uuid.UUID, error) {
	id := r.Header.Get("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return &uuid, nil
}