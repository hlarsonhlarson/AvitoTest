package main

import(
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"AvitoTest/models"
	"time"

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
    // Replace the reference to models.AdvertModel with an interface
    // describing its methods instead. All the other code remains exactly
    // the same.
	adverts interface {
		All() ([]models.Advert, error)
		AddItem(adv models.Advert) (int, int)
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
	// Initalise Env with a models.AdvertModel instance (which in turn wraps
    // the connection pool).
	env := &Env{
		adverts: models.AdvertModel{DB: db},
	}

	http.HandleFunc("/adverts", env.advertsIndex)
	http.HandleFunc("/addadv", env.addAdv)
	http.ListenAndServe(":3000", nil)
}

func (env *Env) advertsIndex(w http.ResponseWriter, r *http.Request) {
    // Execute the SQL query by calling the All() method.
	advs, err := env.adverts.All()
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for _, adv := range advs {
		fmt.Fprintf(w, "%f, %s, %s\n", adv.Price, adv.Name, adv.Description)
		for _, elem := range adv.Photo{
			fmt.Println(elem)
		}
	}
}

func (env *Env) addAdv(w http.ResponseWriter, r *http.Request){
	var adv models.Advert = models.JsonLoader()
	adv.Created_at = time.Now()
	id, response := env.adverts.AddItem(adv)
	if response == 400{
		log.Println("badRequest")
		http.Error(w, http.StatusText(400), 400)
		return
	}
	fmt.Println(id)
}
