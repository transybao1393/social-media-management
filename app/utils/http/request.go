package http

import (
	"log"
	"net/http"
)

type HTTPMethod interface {
	GET()
	POST()
	PUT()
	DELETE()
	PATCH()
}

func GET() {
	//Build The URL string
	URL := "https://jsonplaceholder.typicode.com/todos/1"
	//We make HTTP request using the Get function
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal("ooopsss an error occurred, please try again")
	}
	defer resp.Body.Close()
	//- response body handling
	// var cResp model.Cryptoresponse
	//Decode the data
	// if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
	// 	log.Fatal("ooopsss! an error occurred, please try again")
	// }
	//Invoke the text output function & return it with nil as the error value
	// return cResp.TextOutput(), nil
}

func PUT()    {}
func POST()   {}
func PATCH()  {}
func DELETE() {}
