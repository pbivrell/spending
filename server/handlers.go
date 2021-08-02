package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pbivrell/spending/storage"
	"github.com/pbivrell/spending/storage/sqllite"
)

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {

	var request storage.User

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid create user body", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	if request.Name == "" {
		http.Error(w, "name must be not empty", http.StatusBadRequest)
		return
	}
	if request.Password == "" {
		http.Error(w, "password must be not empty", http.StatusBadRequest)
		return
	}
	if request.Email == "" {
		http.Error(w, "email must be not empty", http.StatusBadRequest)
		return
	}

	err = sqllite.CreateUser(r.Context(), request)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		fmt.Println(err)
		return
	}
}

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {

	var request struct {
		ID int `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid create user body", http.StatusBadRequest)
		return
	}

	users, err := sqllite.ReadUser(r.Context(), "WHERE id=?", request.ID)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	user := users[0]
	user.Password = ""

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {

	var request storage.User

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid update user body", http.StatusBadRequest)
		return
	}

	users, err := sqllite.ReadUser(r.Context(), "WHERE id=?", request.ID)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	user := users[0]

	if request.Name != "" {
		user.Name = request.Name
	}

	if request.Password != "" {
		password, err := HashPassword(request.Password)
		if err != nil {
			http.Error(w, "could not hash password", http.StatusInternalServerError)
			return
		}
		user.Password = password
	}

	if request.Email != "" {
		user.Email = request.Email
	}

	err = sqllite.UpdateUser(r.Context(), user)
	if err != nil {
		http.Error(w, "failed to update", http.StatusInternalServerError)
		return
	}
}

func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {

	var request struct {
		ID int `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid delete user body", http.StatusBadRequest)
		return
	}

	err = sqllite.DeleteUser(r.Context(), request.ID)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
}

func (s *Server) CreatePage(w http.ResponseWriter, r *http.Request) {

	var request storage.Page

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid create user body", http.StatusBadRequest)
		return
	}

	if request.Name == "" {
		http.Error(w, "name must be not empty", http.StatusBadRequest)
		return
	}
	if request.UserID == 0 {
		http.Error(w, "userID must be not empty", http.StatusBadRequest)
		return
	}

	err = sqllite.CreatePage(r.Context(), request)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
}

func (s *Server) GetPage(w http.ResponseWriter, r *http.Request) {

	var request struct {
		ID int `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid create user body", http.StatusBadRequest)
		return
	}

	users, err := sqllite.ReadPage(r.Context(), "WHERE user_id=?", request.ID)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (s *Server) UpdatePage(w http.ResponseWriter, r *http.Request) {

	var request storage.Page

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid update user body", http.StatusBadRequest)
		return
	}

	users, err := sqllite.ReadPage(r.Context(), "WHERE id=?", request.ID)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	user := users[0]

	if request.Name != "" {
		user.Name = request.Name
	}

	if request.UserID != 0 {
		user.UserID = request.UserID
	}

	err = sqllite.UpdatePage(r.Context(), user)
	if err != nil {
		http.Error(w, "failed to update", http.StatusInternalServerError)
		return
	}
}

func (s *Server) DeletePage(w http.ResponseWriter, r *http.Request) {

	var request struct {
		ID int `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid delete user body", http.StatusBadRequest)
		return
	}

	err = sqllite.DeletePage(r.Context(), request.ID)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
}

func (s *Server) CreateEstimate(w http.ResponseWriter, r *http.Request) {

	var request storage.Estimate

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid create user body", http.StatusBadRequest)
		return
	}

	err = sqllite.CreateEstimate(r.Context(), request)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
}

func (s *Server) GetEstimate(w http.ResponseWriter, r *http.Request) {

	var request struct {
		ID int `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid create user body", http.StatusBadRequest)
		return
	}

	users, err := sqllite.ReadEstimate(r.Context(), "WHERE page_id=?", request.ID)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (s *Server) UpdateEstimate(w http.ResponseWriter, r *http.Request) {

	var request storage.Estimate

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid update user body", http.StatusBadRequest)
		return
	}

	users, err := sqllite.ReadEstimate(r.Context(), "WHERE id=?", request.ID)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	user := users[0]

	if request.Name != "" {
		user.Name = request.Name
	}

	if request.PageID != 0 {
		user.PageID = request.PageID
	}

	if request.Amount != 0 {
		user.Amount = request.Amount
	}

	if request.Decimal != 0 {
		user.Decimal = request.Decimal
	}

	if request.Type != "" {
		user.Type = request.Type
	}

	if request.Occurance != nil {
		user.Occurance = request.Occurance
	}

	if request.Period != nil {
		user.Period = request.Period
	}

	if request.Date != nil {
		user.Date = request.Date
	}

	err = sqllite.UpdateEstimate(r.Context(), user)
	if err != nil {
		http.Error(w, "failed to update", http.StatusInternalServerError)
		return
	}
}

func (s *Server) DeleteEstimate(w http.ResponseWriter, r *http.Request) {

	var request struct {
		ID int `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid delete user body", http.StatusBadRequest)
		return
	}

	err = sqllite.DeleteEstimate(r.Context(), request.ID)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
}

func (s *Server) CreateTransaction(w http.ResponseWriter, r *http.Request) {

	var request storage.Transaction

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid create user body", http.StatusBadRequest)
		return
	}

	err = sqllite.CreateTransaction(r.Context(), request)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
}

func (s *Server) GetTransaction(w http.ResponseWriter, r *http.Request) {

	var request struct {
		ID int `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid create user body", http.StatusBadRequest)
		return
	}

	users, err := sqllite.ReadTransaction(r.Context(), "WHERE page_id=?", request.ID)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (s *Server) UpdateTransaction(w http.ResponseWriter, r *http.Request) {

	var request storage.Transaction

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid update user body", http.StatusBadRequest)
		return
	}

	users, err := sqllite.ReadTransaction(r.Context(), "WHERE id=?", request.ID)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	user := users[0]

	if request.Name != "" {
		user.Name = request.Name
	}

	if request.PageID != 0 {
		user.PageID = request.PageID
	}

	if request.Amount != 0 {
		user.Amount = request.Amount
	}

	if request.Decimal != 0 {
		user.Decimal = request.Decimal
	}

	if request.Type != "" {
		user.Type = request.Type
	}

	var t time.Time
	if request.Date != t {
		user.Date = request.Date
	}

	if request.Description != "" {
		user.Description = request.Description
	}

	err = sqllite.UpdateTransaction(r.Context(), user)
	if err != nil {
		http.Error(w, "failed to update", http.StatusInternalServerError)
		return
	}
}

func (s *Server) DeleteTransaction(w http.ResponseWriter, r *http.Request) {

	var request struct {
		ID int `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid delete user body", http.StatusBadRequest)
		return
	}

	err = sqllite.DeleteTransaction(r.Context(), request.ID)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
}

func (s *Server) CreateRole(w http.ResponseWriter, r *http.Request) {

	var request storage.Role

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid create user body", http.StatusBadRequest)
		return
	}

	err = sqllite.CreateRole(r.Context(), request)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
}

func (s *Server) GetRole(w http.ResponseWriter, r *http.Request) {

	var request struct {
		ID int `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid create user body", http.StatusBadRequest)
		return
	}

	users, err := sqllite.ReadRole(r.Context(), "WHERE user_id=?", request.ID)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (s *Server) UpdateRole(w http.ResponseWriter, r *http.Request) {

	var request storage.Role

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid update user body", http.StatusBadRequest)
		return
	}

	users, err := sqllite.ReadRole(r.Context(), "WHERE id=?", request.ID)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	user := users[0]

	if request.PageID != 0 {
		user.PageID = request.PageID
	}

	if request.UserID != 0 {
		user.UserID = request.UserID
	}

	if request.Permission != 0 {
		user.Permission = request.Permission
	}

	err = sqllite.UpdateRole(r.Context(), user)
	if err != nil {
		http.Error(w, "failed to update", http.StatusInternalServerError)
		return
	}
}

func (s *Server) DeleteRole(w http.ResponseWriter, r *http.Request) {

	var request struct {
		ID int `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid delete user body", http.StatusBadRequest)
		return
	}

	err = sqllite.DeleteRole(r.Context(), request.ID)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
}
