package handlers

import (
	"fmt"
	"net/http"
)

func StudentsHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
		case http.MethodGet:
			fmt.Println("Placeholder")
		case http.MethodPost:
			fmt.Println("Placeholder")
		case http.MethodDelete:
			fmt.Println("Placeholder")
		case http.MethodPut:
			fmt.Println("Placeholder")
		case http.MethodPatch:
			fmt.Println("Placeholder")
	}
}