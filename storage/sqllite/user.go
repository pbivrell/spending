package sqllite

import (
	"context"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/pbivrell/spending/storage"
)

const (
	DatabaseFile = "spending.db"
)

func init() {

	const query = `CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		password TEXT NOT NULL,
		email TEXT NOT NULL
	)`

	stmt, err := db.Prepare(query)
	if err != nil {
		fmt.Printf("sqllite table 'users' error: %v\n", err)
		return
	}

	_, err = stmt.Exec()
	if err != nil {
		fmt.Printf("sqllite table 'users' error: %v\n", err)
	}

}

func CreateUser(ctx context.Context, user storage.User) error {
	const query = `INSERT into users (
		name,
		password ,
		email
	) values (?, ?, ?)`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("sqllite table 'users' error: %v", err)
	}

	_, err = stmt.Exec(user.Name, user.Password, user.Email)
	if err != nil {
		return fmt.Errorf("sqllite table 'users' error: %v", err)
	}

	return nil
}

func ReadUser(ctx context.Context, where string, data ...interface{}) ([]storage.User, error) {

	query := fmt.Sprintf(`SELECT id, name, password, email FROM users %s`, where)

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("sqllite table 'users' error: %v", err)
	}

	rows, err := stmt.Query(data)
	if err != nil {
		return nil, fmt.Errorf("sqllite table 'users' error: %v", err)
	}

	defer rows.Close()

	results := make([]storage.User, 0)

	for rows.Next() {
		var user storage.User

		err = rows.Scan(&user.ID, &user.Name, &user.Password, &user.Email)
		if err != nil {
			return nil, fmt.Errorf("sqllite table 'users' error: %v", err)
		}

		results = append(results, user)
	}

	return results, nil

}

func UpdateUser(ctx context.Context, user storage.User) error {
	const query = `UPDATE users set name=?, password=?, email=? where id=?`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("sqllite table 'users' error: %v", err)
	}

	_, err = stmt.Exec(user.Name, user.Password, user.Email, user.ID)
	if err != nil {
		return fmt.Errorf("sqllite table 'users' error: %v", err)
	}

	return nil
}

func DeleteUser(ctx context.Context, id int) error {
	const query = `DELETE FROM users WHERE id=?`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("sqllite table 'users' error: %v", err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("sqllite table 'users' error: %v", err)
	}

	return nil
}
