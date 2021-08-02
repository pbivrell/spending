package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pbivrell/spending/storage"
	"github.com/pbivrell/spending/storage/sqllite"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key")
var expiration = 15 * time.Minute

func AuthMiddlewear(minRole int, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		_, err := validate(r, minRole)
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		handler(w, r)
	}
}

var invalidToken = errors.New("invalid token")
var insufficientPerms = errors.New("insufficient permisions")

func validate(r *http.Request, role int) (Claims, error) {
	claims := Claims{}
	c, err := r.Cookie("token")
	if err != nil {
		return claims, err
	}
	tknStr := c.Value
	tkn, err := jwt.ParseWithClaims(tknStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return claims, err
	}
	if !tkn.Valid {
		return claims, invalidToken
	}

	// if claims.Permission < role {
	//	return claims, insufficientPerms
	// }
	return claims, nil
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {

	var request storage.User

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid login body", http.StatusBadRequest)
		return
	}

	users, err := sqllite.ReadUser(r.Context(), "WHERE name=?", request.Name)
	if err != nil {
		http.Error(w, "failed to login", http.StatusInternalServerError)
		return
	}

	user := users[0]
	user.Password = ""

	if !CheckPasswordHash(request.Password, user.Password) {
		http.Error(w, "invalid user", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(expiration)
	claims := &Claims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

// Claims is the JWT data
type Claims struct {
	storage.User
	jwt.StandardClaims
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
