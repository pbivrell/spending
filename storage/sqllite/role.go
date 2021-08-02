package sqllite

import (
	"context"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/pbivrell/spending/storage"
)

func init() {

	const query = `CREATE TABLE roles(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		permission INTEGER NOT NULL,
		page_id INTEGER NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id),
		FOREIGN KEY(page_id) REFERENCES pages(id)
	)`

	stmt, err := db.Prepare(query)
	if err != nil {
		fmt.Printf("sqllite table 'roles' error: %v\n", err)
		return
	}

	_, err = stmt.Exec()
	if err != nil {
		fmt.Printf("sqllite table 'roles' error: %v\n", err)
	}

}

func CreateRole(ctx context.Context, role storage.Role) error {
	const query = `INSERT into roles (
		user_id,
		permission,
		page_id,
	) values (?, ?, ?)`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("sqllite table 'roles' error: %v", err)
	}

	_, err = stmt.Exec(role.UserID, role.Permission, role.PageID)
	if err != nil {
		return fmt.Errorf("sqllite table 'roles' error: %v", err)
	}

	return nil
}

func ReadRole(ctx context.Context, where string, data ...interface{}) ([]storage.Role, error) {

	query := fmt.Sprintf(`SELECT id, user_id, permission, page_id FROM roles %s`, where)

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("sqllite table 'roles' error: %v", err)
	}

	rows, err := stmt.Query(data)
	if err != nil {
		return nil, fmt.Errorf("sqllite table 'roles' error: %v", err)
	}

	defer rows.Close()

	results := make([]storage.Role, 0)

	for rows.Next() {
		var role storage.Role

		err = rows.Scan(&role.ID, &role.UserID, &role.Permission, &role.PageID)
		if err != nil {
			return nil, fmt.Errorf("sqllite table 'roles' error: %v", err)
		}

		results = append(results, role)
	}

	return results, nil

}

func UpdateRole(ctx context.Context, user storage.Role) error {
	const query = `UPDATE roles set user_id=?, permission=?, page_id=? where id=?`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("sqllite table 'roles' error: %v", err)
	}

	_, err = stmt.Exec(user.UserID, user.Permission, user.PageID, user.ID)
	if err != nil {
		return fmt.Errorf("sqllite table 'roles' error: %v", err)
	}

	return nil
}

func DeleteRole(ctx context.Context, id int) error {
	const query = `DELETE FROM roles WHERE id=?`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("sqllite table 'roles' error: %v", err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("sqllite table 'roles' error: %v", err)
	}

	return nil
}
