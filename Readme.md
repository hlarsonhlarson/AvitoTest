This is test task for Avito internship

SETUP
Hope that this will be started by docker-compose up command, because on my machine it didn't
as people say that's because of docker containers that already exist on my machine (https://github.com/docker-library/postgres/issues/203#issuecomment-255200501).
If it don't wnat to start you need to do things that are presented at file needed_commands.txt.
and init user with parameters presented below
After this you should initialize global variables 
            POSTGRES_USER: "postgres"
            POSTGRES_PASSWORD: "123"
            POSTGRES_DB: "postgres"
            POSTGRES_HOST: "localhost"
            POSTGRES_PORT: "5432"

After this go to src directory and run
go mod download
go build -o avito main.go
and run application by ./main

RUNNING
To pass adverts to database you should pass same command as in file json_pass.txt
You can type another name, description, price in dollars (cause course is unstable) and link to photos, there might be another number of them.
To see and paginate through database you can use your browser by typing localhost:3000/adverts
To pass sorting and order params you should use the command same as in file json_pass1.txt
In browser you can get adverts by clicking on their photo. Parameters that returned already set, but you can change them in the code

METHODS

In the file models.go there are methods which were necessary

func (m AdvertModel) GetPage(offset int, orderType string, typeSorting string) ([]Advert, error)

After connecting to db, it takes ofset as in, orderType and sorting param and return necessary fields. Also it returns id of advert, because it is much simplier to get adverts by id on the next step. The formed query using offset and limit 10. This is prety simple and slow but it works.
Also it return slice of structures which help to get rid of useless massive of variables. In the next method use this structure for the same reason.

 
func (m AdvertModel) GetAdv(id int, fields []string) (Advert, error) 

I choose to get advert by id, cause it's the only unique value, also we can choose timestamp, but form me this is better. It takes id of advert as parameter and fields.
Fields are slice of strings. If we pass string "photos" we can return several photos. If we pass string "description" we cav retrun description in addition to name, price and main photo.


func (m AdvertModel) AddItem(adv Advert) (int, int) 

This function take our structure as a parameter. Cause we fill it with method JsonLoader.
The full application of this method is in function addAdv, but the core is simple
We just pass our structure and form a query row

