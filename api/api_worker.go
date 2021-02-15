package api

import(
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"html/template"
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
    // Replace the reference to models.AdvertModel with an interface
    // describing its methods instead.
	adverts interface {
		GetPage(prevAdv models.Advert, orderType string, typeSorting string) ([]models.Advert, error)
		AddItem(adv models.Advert) (int, int)
		GetAdv(id int) (models.Advert, error)
	}
}

func ApiWorker() {
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
	http.HandleFunc("/addget", env.advertGet)
	http.ListenAndServe(":3000", nil)
}

type AdvertSaver struct{
	Advs []models.Advert
	firstAdv models.Advert
	lastAdv models.Advert
}

func (env *Env) advertsIndex(w http.ResponseWriter, r *http.Request) {
    // Execute the SQL query by calling the All() method.
    	fmt.Println(r)
    	var adv models.Advert
	adv.ID = 1100
	advs, err := env.adverts.GetPage(adv, "DESC", "price")
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	var output AdvertSaver
	output.Advs = advs
	output.
	t := template.Must(template.ParseFiles("./templates/page.html"))
	err = t.Execute(w, output)

}

func (env *Env) advertGet(w http.ResponseWriter, r *http.Request) {
	adv, err := env.adverts.GetAdv(2)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	fmt.Fprintf(w, "%d, %s, %f, %s\n", adv.ID, adv.Name, adv.Price, adv.Photo[0])
}

func (env *Env) addAdv(w http.ResponseWriter, r *http.Request){
	adv, err := models.JsonLoader(w, r)
	if err != nil{
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	id, response := env.adverts.AddItem(adv)
	if response == 400{
		log.Println("badRequest")
		http.Error(w, http.StatusText(400), 400)
		return
	}
	fmt.Println(id)
}
