package handlers

import (
	"fmt"
	"net/http"
)

func ExecsHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
		case http.MethodGet:
			fmt.Println("Placeholder")
		case http.MethodPost:
			fmt.Println("Query: ", r.URL.Query())
			fmt.Println("name: ", r.URL.Query().Get("name"))

			err := r.ParseForm() 
				if err != nil {
					fmt.Println(err)
					return
			}
			fmt.Println("Form from POST: ", r.Form)

		case http.MethodDelete:
			fmt.Println("Placeholder")
		case http.MethodPut:
			fmt.Println("Placeholder")
		case http.MethodPatch:
			fmt.Println("Placeholder")
	}
}