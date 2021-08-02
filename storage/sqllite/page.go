package sqllite

import (
	"context"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/pbivrell/spending/storage"
)

func init() {

	const query = `CREATE TABLE pages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
	)`

	stmt, err := db.Prepare(query)
	if err != nil {
		fmt.Printf("sqllite table 'pages' error: %v\n", err)
		return
	}

	_, err = stmt.Exec()
	if err != nil {
		fmt.Printf("sqllite table 'pages' error: %v\n", err)
	}

}

func CreatePage(ctx context.Context, page storage.Page) error {
	const query = `INSERT into pages (
		name,
		user_id,
	) values (?, ?, ?)`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("sqllite table 'pages' error: %v", err)
	}

	_, err = stmt.Exec(page.Name, page.UserID)
	if err != nil {
		return fmt.Errorf("sqllite table 'pages' error: %v", err)
	}

	return nil
}

func ReadPage(ctx context.Context, where string, data ...interface{}) ([]storage.Page, error) {

	query := fmt.Sprintf(`SELECT id, name, user_id FROM pages %s`, where)

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("sqllite table 'pages' error: %v", err)
	}

	rows, err := stmt.Query(data)
	if err != nil {
		return nil, fmt.Errorf("sqllite table 'pages' error: %v", err)
	}

	defer rows.Close()

	results := make([]storage.Page, 0)

	for rows.Next() {
		var page storage.Page

		err = rows.Scan(&page.ID, &page.Name, &page.UserID)
		if err != nil {
			return nil, fmt.Errorf("sqllite table 'pages' error: %v", err)
		}
		results = append(results, page)
	}

	return results, nil

}

func UpdatePage(ctx context.Context, page storage.Page) error {
	const query = `UPDATE users set name=?, user_id=? where id=?`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("sqllite table 'pages' error: %v", err)
	}

	_, err = stmt.Exec(page.Name, page.UserID, page.ID)
	if err != nil {
		return fmt.Errorf("sqllite table 'pages' error: %v", err)
	}

	return nil
}

func DeletePage(ctx context.Context, id int) error {
	const query = `DELETE FROM pages WHERE id=?`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("sqllite table 'pages' error: %v", err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("sqllite table 'pages' error: %v", err)
	}

	return nil
}
