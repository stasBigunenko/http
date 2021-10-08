package handlers

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"net/http"
	"src/http/pkg/model"
	"src/http/pkg/services"
	"strings"
	"time"
)

const SECRETKEY = "password"

//type Claims struct {
//	Name string `json:"name"`
//	jwt.StandardClaims
//}

func (h *PostHandler) VerifyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := "/posts/login"
		requestPath := r.URL.Path

		if notAuth == requestPath {
			next.ServeHTTP(w, r)
			return
		}

		token, err := GetJwtTokenFromRequest(r)
		if err != nil {
			msg := services.Response("Invalid token")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(msg)
			return
		}

		claims, ok := validateToken(token)
		if !ok {
			msg := services.Response("Invalid token")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(msg)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "claims", claims))
		next.ServeHTTP(w, r)
	})
}

func GetJwtTokenFromRequest(r *http.Request) (string, error) {

	tokenInHeaderVal := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	if len(tokenInHeaderVal) != 2 {
		msg := services.Response("Malformed token")
		return "", fmt.Errorf(string(msg))
	}

	return tokenInHeaderVal[1], nil
}

func validateToken(jwtToken string) (model.Claims, bool) {

	claims := &model.Claims{}

	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRETKEY), nil
	})
	if err != nil {
		return *claims, false
	}

	if !token.Valid {
		return *claims, false
	}

	return *claims, true
}

func (h *PostHandler) generateTokenStringForUser(name string) (string, error) {

	expirationTime := time.Now().Add(24 * time.Hour)
	// Create the JWT claims, which includes the username and expiry time
	claims := model.Claims{
		// In JWT, the expiry time is expressed as unix milliseconds
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Id:        uuid.New().String(),
		},
		Name: name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SECRETKEY))
	return tokenString, err
}
