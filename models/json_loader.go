package models

import(
    "net/http"
    "time"
)

func JsonLoader(rw http.ResponseWriter, request *http.Request) (Advert, error){
	var adv Advert

	err := DecodeJSONBody(rw, request, &adv)
	if err != nil{
		return adv, err
	}
	adv.Created_at = time.Now()
	return adv, nil
}
