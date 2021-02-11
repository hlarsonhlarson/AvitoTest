package models

import (
	"database/sql"
)

const BAD_REQUEST := 400
const GOOD_REQUEST := 200

type Book struct {
	Isbn   string
	Title  string
	Author string
	Price  float32
}

// Create a custom BookModel type which wraps the sql.DB connection pool.
type BookModel struct {
	DB *sql.DB
}

// Use a method on the custom BookModel type to run the SQL query.
func (m BookModel) All() ([]Book, error) {
	rows, err := m.DB.Query("SELECT * FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bks []Book

	for rows.Next() {
		var bk Book

		err := rows.Scan(&bk.Isbn, &bk.Title, &bk.Author, &bk.Price)
		if err != nil {
			return nil, err
		}

		bks = append(bks, bk)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bks, nil
}

func (m BookModel) AddItem(price, name, description string, photo string[]) (int, int){
	sqlStatement := `
	INSERT INTO books (price, name, description, photo)
	VALUES ($1, $2, $3, $4)
	RETURNING id`
	id := 0
	err = db.QueryRow(sqlStatement, price, name, description, photo)
	if err != nil{
		fmt.Println(err)
		return nil, BAD_REQUEST
	}
	return id, GOOD_REQUEST
}
