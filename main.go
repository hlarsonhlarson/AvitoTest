package main

import(
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"AvitoTest/models"

	_ "github.com/lib/pq"

)

const (
	host     = "localhost"
	port     = 5432
	user     = "AvitoTest"
	password = "123"
	dbname   = "AvitoTest"
)

type Env struct {
    // Replace the reference to models.BookModel with an interface
    // describing its methods instead. All the other code remains exactly
    // the same.
	books interface {
		All() ([]models.Book, error)
	}
}

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	// Initalise Env with a models.BookModel instance (which in turn wraps
    // the connection pool).
	env := &Env{
		books: models.BookModel{DB: db},
	}

	http.HandleFunc("/books", env.booksIndex)
	http.ListenAndServe(":3000", nil)
}

func (env *Env) booksIndex(w http.ResponseWriter, r *http.Request) {
    // Execute the SQL query by calling the All() method.
	bks, err := env.books.All()
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for _, bk := range bks {
		fmt.Fprintf(w, "%s, %s, %s, £%.2f\n", bk.Isbn, bk.Title, bk.Author, bk.Price)
	}
}
