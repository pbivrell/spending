package sqllite

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/pbivrell/spending/storage"
)

var db *sql.DB

func init() {

	var err error

	db, err = sql.Open("sqlite3", DatabaseFile)
	if err != nil {
		panic(fmt.Sprintf("sqllite connection could not be opened: %v", err))
	}

	const query = `CREATE TABLE estimate (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		page_id INTEGER NOT NULL,
		amount INTEGER NOT NULL,
		decimal INTEGER NOT NULL,
		type TEXT NOT NULL,
		name TEXT NOT NULL,
		date DATETIME,
		period INTEGER,
		occurance INTEGER,
		FOREIGN KEY(page_id) REFERENCES pages(id)
	)`

	stmt, err := db.Prepare(query)
	if err != nil {
		fmt.Printf("sqllite table 'estimate' error: %v\n", err)
		return
	}

	_, err = stmt.Exec()
	if err != nil {
		fmt.Printf("sqllite table 'estimate' error: %v\n", err)
	}

}

func CreateEstimate(ctx context.Context, estimate storage.Estimate) error {
	const query = `INSERT into estimate (
		page_id,		
		amount,
		decimal,
		type,
		name,
		date,
		period,
		occurance
	) values (?, ?, ?, ?, ?, ?, ?, ?)`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("sqllite table 'estimate' error: %v", err)
	}

	_, err = stmt.Exec(estimate.PageID, estimate.Amount, estimate.Decimal, estimate.Type, estimate.Name, estimate.Date, estimate.Period, estimate.Occurance)
	if err != nil {
		return fmt.Errorf("sqllite table 'estimate' error: %v", err)
	}

	return nil
}

func ReadEstimate(ctx context.Context, where string, data ...interface{}) ([]storage.Estimate, error) {

	query := fmt.Sprintf(`SELECT id, page_id, amount, decimal, tye, name, date, occurance, period FROM estimate %s`, where)

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("sqllite table 'estimate' error: %v", err)
	}

	rows, err := stmt.Query(data...)
	if err != nil {
		return nil, fmt.Errorf("sqllite table 'estimate' error: %v", err)
	}

	defer rows.Close()

	results := make([]storage.Estimate, 0)

	for rows.Next() {
		var estimate storage.Estimate

		err = rows.Scan(&estimate.ID, &estimate.PageID, &estimate.Amount, &estimate.Decimal, &estimate.Type, &estimate.Name, estimate.Date, estimate.Occurance, estimate.Period)
		if err != nil {
			return nil, fmt.Errorf("sqllite table 'estimate' error: %v", err)
		}
		results = append(results, estimate)
	}

	return results, nil

}

func UpdateEstimate(ctx context.Context, estimate storage.Estimate) error {
	const query = `UPDATE estimate set page_id=?, amount=?, decimal=?, type=?, name=?, date=?, occurance=?, period=? where id=?`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("sqllite table 'estimate' error: %v", err)
	}

	_, err = stmt.Exec(estimate.PageID, estimate.Amount, estimate.Decimal, estimate.Type, estimate.Name, estimate.Date, estimate.Occurance, estimate.Period)
	if err != nil {
		return fmt.Errorf("sqllite table 'estimate' error: %v", err)
	}

	return nil
}

func DeleteEstimate(ctx context.Context, id int) error {
	const query = `DELETE FROM estimate WHERE id=?`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("sqllite table 'estimate' error: %v", err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("sqllite table 'estimate' error: %v", err)
	}

	return nil
}
