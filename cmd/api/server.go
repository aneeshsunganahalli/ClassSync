package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	port := ":3000"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		// fmt.Fprintf(w, "Hello Root")
		w.Write([]byte("Hello Root"))
		fmt.Println("Hello Root")
	})

	http.HandleFunc("/teachers", func(w http.ResponseWriter, r *http.Request){
		if r.Method == http.MethodGet {
			w.Write([]byte("Hello GET Method on Teachers"))
			fmt.Println("Hello GET Method on Teachers")
			return 
		}
		w.Write([]byte("Hello Teachers Route"))
		fmt.Println("Hello Teachers Route")
	})

	http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Hello Students Route"))
		fmt.Println("Hello Students Route")
	})

	http.HandleFunc("/execs", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Hello Execs Route"))
		fmt.Println("Hello Execs Route")
	})

	fmt.Println("Server is running on port: ", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("Error starting the server", err)
	}
	
}