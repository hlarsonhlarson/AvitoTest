package models

import (
	"database/sql"
	"time"
	"fmt"
	"errors"
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

type AdvertModel struct {
	DB *sql.DB
}

func (m AdvertModel) GetPage(offset int, orderType string, typeSorting string) ([]Advert, error) {
	sqlStatement := `SELECT name, price, photo, id
	FROM adverts `
	switch {
	case orderType == "ASC" && typeSorting == "price":
		sqlStatement += `ORDER BY price,id `
	case orderType == "ASC" && typeSorting == "time":
		sqlStatement += `ORDER BY time,id `
	case orderType == "DESC" && typeSorting == "price":
		sqlStatement += `ORDER BY price DESC, id `
	case orderType == "DESC" && typeSorting == "time":
		sqlStatement += `ORDER BY time DESC,id `
	default:
		return nil, errors.New("Incorect page statement")
	}
	sqlStatement += `LIMIT 10
	OFFSET $1;`
	rows, err := m.DB.Query(sqlStatement, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var advs []Advert

	for rows.Next() {
		var adv Advert

		err := rows.Scan(&adv.Name, &adv.Price, pq.Array(&adv.Photo), &adv.ID)
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


func (m AdvertModel) GetAdv(id int, fields []string) (Advert, error) {
	var adv Advert
	photos := 0
	descr := 0
	sqlStatement := `SELECT name, price, photo `
	for _, elem := range fields{
		if elem == "photos"{
			photos = 1
		}else if elem == "description"{
			sqlStatement += `, description `
			descr += 1
		}
	}
	sqlStatement += `FROM adverts WHERE id=$1;`
	row := m.DB.QueryRow(sqlStatement, id)
	var err error
	if descr == 1{
		err = row.Scan(&adv.Name, &adv.Price, pq.Array(&adv.Photo), &adv.Description)
	}else{
		err = row.Scan(&adv.Name, &adv.Price, pq.Array(&adv.Photo))
	}
	switch err {
	case sql.ErrNoRows:
		return adv, errors.New("No rows were returned!")
	default:
		if photos == 0{
				tmp_photo := make([]string, 0)
				tmp_photo = append(tmp_photo, adv.Photo[0])
				adv.Photo = tmp_photo
			}
		return adv, err
	}
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
