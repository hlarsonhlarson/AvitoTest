package models

import (
	"database/sql"
	"time"
	"fmt"
	pq "github.com/lib/pq"
)

const BAD_REQUEST = 400
const GOOD_REQUEST = 200

type Advert struct {
	ID int
	Price  float32
	Name string
	Description string
	Photo []string
	Created_at time.Time
}

// Create a custom AdvertModel type which wraps the sql.DB connection pool.
type AdvertModel struct {
	DB *sql.DB
}

// Use a method on the custom AdvertModel type to run the SQL query.
func (m AdvertModel) All() ([]Advert, error) {
	rows, err := m.DB.Query("SELECT * FROM adverts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var advs []Advert

	for rows.Next() {
		var adv Advert

		err := rows.Scan(&adv.ID, &adv.Price, &adv.Name, &adv.Description, pq.Array(&adv.Photo), &adv.Created_at)
		if err != nil {
			return nil, err
		}

		advs = append(advs, adv)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return advs, nil
}

func (m AdvertModel) AddItem(adv Advert) (int, int) {
	sqlStatement := `
	INSERT INTO adverts (price, name, description, photo, created_at)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id`
	id := 0
	err := m.DB.QueryRow(sqlStatement, adv.Price, adv.Name, adv.Description, pq.Array(adv.Photo), adv.Created_at).Scan(&id)
	if err != nil{
		fmt.Println(err)
		return 0, BAD_REQUEST
	}
	return id, GOOD_REQUEST
}
