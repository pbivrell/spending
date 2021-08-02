package sqllite

import (
	"context"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/pbivrell/spending/storage"
)

func init() {

	const query = `CREATE TABLE transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT NOT NULL,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		amount INTEGER NOT NULL,
		decimal INTEGER NOT NULL,
		page_id INTEGER NOT NULL,
		date DATETIME NOT NULL,
		FOREIGN KEY(page_id) REFERENCES pages(id)
	)`

	stmt, err := db.Prepare(query)
	if err != nil {
		fmt.Printf("sqllite table 'transactions' error: %v\n", err)
		return
	}

	_, err = stmt.Exec()
	if err != nil {
		fmt.Printf("sqllite table 'transactions' error: %v\n", err)
	}

}

func CreateTransaction(ctx context.Context, transaction storage.Transaction) error {
	const query = `INSERT into transactions (
		page_id,		
		amount,
		decimal,
		type,
		name,
		description,
		date,
	) values (?, ?, ?, ?, ?, ?, ?)`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("sqllite table 'transactions' error: %v", err)
	}

	_, err = stmt.Exec(transaction.PageID, transaction.Amount, transaction.Decimal, transaction.Type, transaction.Name, transaction.Description, transaction.Date)
	if err != nil {
		return fmt.Errorf("sqllite table 'transactions' error: %v", err)
	}

	return nil
}

func ReadTransaction(ctx context.Context, where string, data ...interface{}) ([]storage.Transaction, error) {

	query := fmt.Sprintf(`SELECT id, page_id, amount, decimal, tye, name, description, date FROM transactions %s`, where)

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("sqllite table 'transactions' error: %v", err)
	}

	rows, err := stmt.Query(data)
	if err != nil {
		return nil, fmt.Errorf("sqllite table 'transactions' error: %v", err)
	}

	defer rows.Close()

	results := make([]storage.Transaction, 0)

	for rows.Next() {
		var transaction storage.Transaction

		err = rows.Scan(&transaction.ID, &transaction.PageID, &transaction.Amount, &transaction.Decimal, &transaction.Type, &transaction.Name, &transaction.Description, &transaction.Date)
		if err != nil {
			return nil, fmt.Errorf("sqllite table 'transactions' error: %v", err)
		}
		results = append(results, transaction)
	}

	return results, nil

}

func UpdateTransaction(ctx context.Context, transaction storage.Transaction) error {
	const query = `UPDATE transaction set page_id=?, amount=?, decimal=?, type=?, name=?, description=?, date=? where id=?`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("sqllite table 'transactions' error: %v", err)
	}

	_, err = stmt.Exec(transaction.PageID, transaction.Amount, transaction.Decimal, transaction.Type, transaction.Name, transaction.Description, transaction.Date)
	if err != nil {
		return fmt.Errorf("sqllite table 'transactions' error: %v", err)
	}

	return nil
}

func DeleteTransaction(ctx context.Context, id int) error {
	const query = `DELETE FROM transactions WHERE id=?`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("sqllite table 'transactions' error: %v", err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("sqllite table 'transactions' error: %v", err)
	}

	return nil
}
