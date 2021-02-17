package api

import(
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"html/template"
	"AvitoTest/models"
	"os"

	_ "github.com/lib/pq"

)

type Env struct {
	adverts interface {
		GetPage(offset int, orderType string, typeSorting string) ([]models.Advert, error)
		AddItem(adv models.Advert) (int, int)
		GetAdv(id int, fields []string) (models.Advert, error)
	}
	Params models.Parameters
}

type AdvertWraper struct{
	Adverts []models.Advert
}


func ApiWorker() {
	host := os.Getenv("POSTGRES_HOST")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_USER")
	port := os.Getenv("POSTGRES_PORT")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
	"password=%s dbname=%s sslmode=disabled",
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
	_, err = db.Query(`CREATE TABLE IF NOT EXISTS adverts(
		id SERIAL PRIMARY KEY,
		price FLOAT,
		name TEXT,
		description TEXT,
		photo TEXT[],
		created_at TIMESTAMP);
		`)
	if err != nil {
		panic(err)
	}
	env := &Env{
		adverts: models.AdvertModel{DB: db},
	}
	env.Params.SortingOrder = "ASC"
	env.Params.SortingParameter = "price"
	rows, err := db.Query(`SELECT count(*) from adverts;`)
	if err != nil {
		panic(err)
	}
	rows.Next()
	err = rows.Scan(&env.Params.Elem_nums)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/setparams", env.advertsPassParams)
	http.HandleFunc("/adverts", env.advertsIndex)
	http.HandleFunc("/addadv", env.addAdv)
	http.HandleFunc("/addget", env.advertGet)
	http.ListenAndServe(":3000", nil)
}

func (env *Env) advertsPassParams(w http.ResponseWriter, r *http.Request){
	params_tmp,err := models.JsonLoadParams(w, r)

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	env.Params.SortingOrder = params_tmp.SortingOrder
	env.Params.SortingParameter = params_tmp.SortingParameter
}


func (env *Env) advertsIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.RawQuery == "prev_page"{
		env.Params.Offset -= 10
		if env.Params.Offset < 0{
			env.Params.Offset = 0
		}
	}else if r.URL.RawQuery == "next_page"{
		env.Params.Offset += 10
		if env.Params.Offset > env.Params.Elem_nums{
			env.Params.Offset -= 10
		}
	}else{
		env.Params.Offset = 0
	}
	advs, err := env.adverts.GetPage(env.Params.Offset, env.Params.SortingOrder, env.Params.SortingParameter)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	var wraper AdvertWraper
	wraper.Adverts = advs
	t := template.Must(template.ParseFiles("./templates/page.html"))
	err = t.Execute(w, wraper)

}

func (env *Env) advertGet(w http.ResponseWriter, r *http.Request) {
	var ID int
	var err error
	if r.URL.RawQuery != ""{
		ID, err = strconv.Atoi(r.URL.RawQuery)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
	}
	fields := make([]string, 0)
	fields = append(fields, "description")
	adv, err := env.adverts.GetAdv(ID, fields)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	t := template.Must(template.ParseFiles("./templates/item.html"))
	err = t.Execute(w, adv)
}

func (env *Env) addAdv(w http.ResponseWriter, r *http.Request){
	adv, err := models.JsonLoader(w, r)
	if err != nil{
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	_, response := env.adverts.AddItem(adv)
	if response == 400{
		log.Println("badRequest")
		http.Error(w, http.StatusText(400), 400)
		return
	}
	env.Params.Elem_nums += 1
}
