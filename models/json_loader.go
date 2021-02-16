package models

import(
    "net/http"
    "net/url"
    "time"
    "fmt"
    "errors"
    "unicode/utf8"
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

func CheckLen(str string, length int) bool{
	if utf8.RuneCountInString(str) > length{
		return false
	}
	return true
}

func CheckNameDescription(name, description string) error{
	if CheckLen(name, 200) == false{
		return errors.New("Too long name")
	}
	if utf8.RuneCountInString(name) == 0{
		return errors.New("No name")
	}
	if CheckLen(description, 1000) == false{
		return errors.New("Too long description")
	}
	if utf8.RuneCountInString(description) == 0{
		return errors.New("No description")
	}
	return nil
}


func JsonLoader(rw http.ResponseWriter, request *http.Request) (Advert, error){
	var adv Advert

	err := DecodeJSONBody(rw, request, &adv)
	if err != nil{
		return adv, err
	}
	err = CheckNameDescription(adv.Name, adv.Description)
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

type Parameters struct{
	SortingOrder string
	SortingParameter string
}

func JsonLoadParams(rw http.ResponseWriter, request *http.Request) (Parameters, error){
	var par Parameters
	err := DecodeJSONBody(rw, request, &par)

	if err != nil{
		return par, err
	}
	return par, nil
}
