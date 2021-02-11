package models

import(
    "fmt"
    "encoding/json"
    "net/http"
)

func JsonLoader(rw http.ResponseWriter, request *http.Request) (Advert, error){
	decoder := json.NewDecoder(request.Body)

	var adv Advert

	err := decoder.Decode(&adv)
	if err != nil{
		return adv, err
	}
	fmt.Println("HI")
	return adv, nil
}
