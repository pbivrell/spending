package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/felixge/httpsnoop"
	"github.com/gorilla/mux"
	"github.com/pbivrell/gatekeeper/log"
	"golang.org/x/crypto/bcrypt"

	ua "github.com/mileusna/useragent"
)

var JwtKey string

type Server struct {
	u UserdataStorage
}

func New(u UserdataStorage) *Server {
	return &Server{
		u: u,
	}
}

func LogMiddlewear(logger log.Logger, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := "n/a"
		c, err := r.Cookie("token")
		if err == nil {
			claims := Claims{}
			_, err = jwt.ParseWithClaims(c.Value, &claims, func(token *jwt.Token) (interface{}, error) {
				return JwtKey, nil
			})
			if err == nil {
				user = claims.Username
			}
		}

		t := time.Now()

		// Preform the request and capture metrics about the result
		m := httpsnoop.CaptureMetrics(handler, w, r)

		// Parse useragent
		ua := ua.Parse(r.Header.Get("User-Agent"))

		deviceType := "n/a"
		if ua.Bot {
			deviceType = "bot"
		} else if ua.Mobile {
			deviceType = "mobile"
		} else if ua.Tablet {
			deviceType = "tablet"
		} else if ua.Desktop {
			deviceType = "desktop"
		}

		logger.Write(log.Data{
			Method:           r.Method,
			Endpoint:         r.URL.String(),
			Referer:          r.Header.Get("Referer"),
			Code:             m.Code,
			Duration:         m.Duration.Milliseconds(),
			IP:               r.RemoteAddr[0:strings.LastIndex(r.RemoteAddr, ":")],
			UserAgent:        ua.Name,
			UserAgentVersion: ua.Version,
			OS:               ua.OS,
			OSVersion:        ua.OSVersion,
			Device:           ua.Device,
			DeviceType:       deviceType,
			Time:             t,
			User:             user,
		})
	}
}

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
		return JwtKey, nil
	})
	if err != nil {
		return claims, err
	}
	if !tkn.Valid {
		return claims, invalidToken
	}

	if claims.Role < role {
		return claims, insufficientPerms
	}
	return claims, nil
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {

	var userdata Userdata

	err := json.NewDecoder(r.Body).Decode(&userdata)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	c, _ := r.Cookie("token")
	claims := Claims{}
	_, _ = jwt.ParseWithClaims(c.Value, &claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if userdata.Role >= claims.Role {
		http.Error(w, "can not create user with higher permission then yourself", http.StatusUnauthorized)
		return
	}

	err = s.u.Insert(userdata)
	if err != nil {
		http.Error(w, "could not create user", http.StatusInternalServerError)
		return
	}

}

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("token")
	claims := Claims{}
	_, _ = jwt.ParseWithClaims(c.Value, &claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	user, err := s.u.Search(claims.Username, true)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	user.Password = ""

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("token")
	claims := Claims{}
	_, _ = jwt.ParseWithClaims(c.Value, &claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	user := mux.Vars(r)["user"]

	// Only an admin can delete a user other than themself
	if claims.Username != user && claims.Role < OwnerRole {
		http.Error(w, "", http.StatusForbidden)
		return
	}

	// Lookup user to ensure the user we are deleting does not have the
	// same permissions as us
	foundUser, err := s.u.Search(claims.Username, true)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if foundUser.Role >= claims.Role {
		http.Error(w, "you don't have permission to delete this user", http.StatusForbidden)
		return

	}

	err = s.u.Remove(claims.Username)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

}

func (s *Server) User(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("token")
	claims := Claims{}
	_, _ = jwt.ParseWithClaims(c.Value, &claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	user := mux.Vars(r)["user"]

	if claims.Username != user && claims.Role < AdminRole {
		http.Error(w, "", http.StatusForbidden)
		return
	}
	result, err := s.u.Search(user, true)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

}

func (s *Server) Users(w http.ResponseWriter, r *http.Request) {
	sstart := r.URL.Query().Get("start")
	scount := r.URL.Query().Get("count")

	start := 0
	count := 50

	if s, err := strconv.Atoi(sstart); err == nil {
		start = s
	}

	if c, err := strconv.Atoi(scount); err == nil {
		count = c
	}

	users, err := s.u.List(start, count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {

	request := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid login body", http.StatusBadRequest)
		return
	}

	user, err := s.u.Search(request.Username, false)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to lookup user: %v", err), http.StatusInternalServerError)
		return
	}

	if !CheckPasswordHash(request.Password, user.Password) {
		http.Error(w, "invalid user", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
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

func Refresh(w http.ResponseWriter, r *http.Request) {

	c, _ := r.Cookie("token")
	claims := Claims{}
	_, _ = jwt.ParseWithClaims(c.Value, &claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the new token as the users `token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

// Claims is the JWT data
type Claims struct {
	Username string `json:"username"`
	Role     int    `json:"role"`
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

type Proxy struct {
	Role int    `json:"role"`
	Url  string `json:"url"`
}

func LoadProxies(path string) (map[string]Proxy, error) {
	var proxys map[string]Proxy
	f, err := os.Open(path)
	if err != nil {
		return proxys, err
	}
	defer f.Close()

	return proxys, json.NewDecoder(f).Decode(&proxys)
}

func (s *Server) Proxy(lock *sync.Mutex, proxys map[string]Proxy) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get proxyName from url path
		vars := mux.Vars(r)
		proxyName, ok := vars["proxy"]
		if !ok {
			http.Error(w, "no proxy provided", http.StatusBadRequest)
			return
		}

		// Get proxy from storage protected by mutex
		lock.Lock()
		proxy, ok := proxys[proxyName]
		lock.Unlock()
		if !ok {
			http.Error(w, "invalid proxy", http.StatusBadRequest)
			return
		}

		// Get request path from url
		path, ok := vars["path"]
		if !ok {
			path = "/"
		}

		fmt.Println("Hey role", proxy.Role)

		_, err := validate(r, 0)
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		req, err := http.NewRequest(r.Method, proxy.Url+path, r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to make request: %+v", err), http.StatusInternalServerError)
			return
		}
		// Copy headers
		for v, h := range r.Header {
			for _, z := range h {
				req.Header.Add(v, z)
			}
		}
		req.Header.Add("X-Forwarded-Host", req.Host)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to do: %+v", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		for v, h := range resp.Header {
			for _, z := range h {
				w.Header().Add(v, z)
			}
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)

	}
}
