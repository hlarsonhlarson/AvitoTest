package models

import(
    "net/http"
    "net/url"
    "time"
    "fmt"
    "errors"
)

func IsUrl(str string) bool {
	_, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}
	u, err := url.Parse(str)
	if err != nil{
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https"{
		return false
	}
	return err == nil && u.Scheme != "" && u.Host != ""
}

func CheckPhoto(photos []string) error{
	if len(photos) > 3{
		return errors.New("There are more than three photos")
	}
	if len(photos) == 0{
		return errors.New("There are no photo links")
	}
	for photo_num, photo := range photos{
		if IsUrl(photo) == false{
			s := fmt.Sprintf("Invalid url for photo %d", photo_num + 1)
			return errors.New(s)
		}
		resp, err := http.Get(photo)
		if err != nil || resp.StatusCode != http.StatusOK{
			s := fmt.Sprintf("Not reachable link number %d", photo_num + 1)
			return errors.New(s)
		}
	}
	return nil
}


func JsonLoader(rw http.ResponseWriter, request *http.Request) (Advert, error){
	var adv Advert

	err := DecodeJSONBody(rw, request, &adv)
	if err != nil{
		return adv, err
	}
	err = CheckPhoto(adv.Photo)
	if err != nil{
		return adv, err
	}
	adv.Created_at = time.Now()
	return adv, nil
}
