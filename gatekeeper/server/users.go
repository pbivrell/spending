package server

import (
	"errors"
	"fmt"
	"sync"

	"database/sql"
)

const (
	OwnerRole  = 0b10000000
	AdminRole  = 0b01000000
	FriendRole = 0b00001000
	FamilyRole = 0b00000100
	ReadRole   = 0b00000001
)

// Userdata encode a users authenication userdata
type Userdata struct {
	// Username is the user's name
	Username string `json:"username"`
	// Password is the user's password
	Password string `json:"password,omitempty"`
	// Email is the users email
	Email string `json:"email"`
	// Role is the role
	Role int `json:"role"`
}

// MissingUserdata indicates that the userdata could not be found
var MissingUserdata = errors.New("Could not find userdata")

// DuplicateUserdata indicates that the user already exists
var DuplicateUserdata = errors.New("Userdata already exist")

// UserdataStorage is an interface for aquiring userdata
type UserdataStorage interface {
	Search(string, bool) (Userdata, error)
	Insert(Userdata) error
	Remove(string) error
	List(int, int) ([]Userdata, error)
}

// MemoryUserdataStorage is an in memory method for storing Userdata
type MemoryUserdataStorage struct {
	l       *sync.Mutex
	storage map[string]*Userdata
}

// Returns a new memory userdata
func NewMemoryUserdataStorage() *MemoryUserdataStorage {
	return &MemoryUserdataStorage{
		l:       &sync.Mutex{},
		storage: make(map[string]*Userdata),
	}
}

// Search retrieves the userdata from the map
func (m *MemoryUserdataStorage) Search(username string, hidePassword bool) (Userdata, error) {
	m.l.Lock()
	defer m.l.Unlock()

	data, ok := m.storage[username]

	if !ok {
		return Userdata{}, MissingUserdata
	}

	result := *data

	if hidePassword {
		result.Password = ""
	}

	return result, nil
}

// Insert adds the userdata to the map
func (m *MemoryUserdataStorage) Insert(data Userdata) error {
	m.l.Lock()
	defer m.l.Unlock()

	_, ok := m.storage[data.Username]
	if ok {
		return DuplicateUserdata
	}

	m.storage[data.Username] = &data

	return nil
}

// Delete removes the userdata from the map
func (m *MemoryUserdataStorage) Delete(username string) error {
	m.l.Lock()
	defer m.l.Unlock()

	delete(m.storage, username)

	return nil
}

func (m *MemoryUserdataStorage) List(start, count int) ([]Userdata, error) {
	m.l.Lock()
	defer m.l.Unlock()

	users := make([]Userdata, 0)

	i := 0
	for _, value := range m.storage {
		if i > count {
			break
		}
		users = append(users, Userdata{
			Username: value.Username,
			Role:     value.Role,
			Email:    value.Email,
		})
	}

	return users, nil
}

// SQLUserdataStorage stores userdata in the sql Database provided
type SQLUserdataStorage struct {
	db *sql.DB
}

const createUsersTableQuery = `CREATE TABLE IF NOT EXISTS users (
	uid INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	password TEXT NOT NULL,
	username TEXT NOT NULL UNIQUE,
	email TEXT NOT NULL,
	role INTEGER NOT NULL
);`

const searchUsersTableQuery = `SELECT username,password,email,role FROM users WHERE username=? LIMIT 1`

const removeUsersTableQuery = `DELETE FROM users WHERE username=?`

const insertUsersTableQuery = `INSERT INTO users(username, password, email, role) VALUES (?, ?, ?, ?)`

const listUsersTableQuery = `SELECT username,email,role FROM users LIMIT ? OFFSET ?`

// FailedUserdataSQL is the error associated with a sql query
var FailedUserdataSQL = errors.New("userdata SQL error")

// NewSQLUserdataStorage returns a userdata storage interface using a SQL connection as
// backing store
func NewSQLUserdataStorage(db *sql.DB) (*SQLUserdataStorage, error) {

	statement, err := db.Prepare(createUsersTableQuery) // Prepare SQL Statement
	if err != nil {
		return nil, fmt.Errorf("%w: %v", FailedUserdataSQL, err)
	}
	_, err = statement.Exec()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", FailedUserdataSQL, err)
	}
	return &SQLUserdataStorage{
		db: db,
	}, nil
}

// Search finds the provided username returning an error if one could not be found
func (s *SQLUserdataStorage) Search(username string, hidePassword bool) (Userdata, error) {
	statement, err := s.db.Prepare(searchUsersTableQuery) // Prepare SQL Statement
	if err != nil {
		return Userdata{}, fmt.Errorf("%w: %v", FailedUserdataSQL, err)
	}
	result, err := statement.Query(username)
	if err != nil {
		return Userdata{}, fmt.Errorf("%w: %v", FailedUserdataSQL, err)
	}

	var password, email string
	var role int

	result.Next()
	err = result.Scan(&username, &password, &email, &role)
	if err != nil {
		return Userdata{}, fmt.Errorf("%w: %v", FailedUserdataSQL, err)

	}

	if hidePassword {
		password = ""
	}

	result.Close()

	return Userdata{
		Username: username,
		Password: password,
		Email:    email,
		Role:     role,
	}, nil
}

// Insert user data into database
func (s *SQLUserdataStorage) Insert(data Userdata) error {
	statement, err := s.db.Prepare(insertUsersTableQuery) // Prepare SQL Statement
	if err != nil {
		return fmt.Errorf("%w: %v", FailedUserdataSQL, err)
	}
	_, err = statement.Exec(data.Username, data.Password, data.Email, data.Role)
	if err != nil {
		return fmt.Errorf("%w: %v", FailedUserdataSQL, err)
	}
	return nil
}

// Remove user from database
func (s *SQLUserdataStorage) Remove(username string) error {
	statement, err := s.db.Prepare(removeUsersTableQuery) // Prepare SQL Statement
	if err != nil {
		return fmt.Errorf("%w: %v", FailedUserdataSQL, err)
	}
	_, err = statement.Exec(username)
	if err != nil {
		return fmt.Errorf("%w: %v", FailedUserdataSQL, err)

	}
	return nil
}

func (s *SQLUserdataStorage) List(start, count int) ([]Userdata, error) {
	statement, err := s.db.Prepare(listUsersTableQuery) // Prepare SQL Statement
	if err != nil {
		return nil, fmt.Errorf("%w: %v", FailedUserdataSQL, err)
	}
	result, err := statement.Query(count, start)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", FailedUserdataSQL, err)
	}

	users := make([]Userdata, 0)

	var username, email string
	var role int

	for result.Next() {
		err = result.Scan(&username, &email, &role)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", FailedUserdataSQL, err)
		}
		users = append(users, Userdata{
			Username: username,
			Email:    email,
			Role:     role,
		})
	}

	result.Close()

	return users, nil
}
