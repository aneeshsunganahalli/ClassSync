package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aneeshsunganahalli/ClassSync/internal/api/middlewares"
	mw "github.com/aneeshsunganahalli/ClassSync/internal/api/middlewares"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server Healthy"))
		fmt.Println("Server Healthy")
}

func teachersHandler(w http.ResponseWriter, r *http.Request) {

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

func studentsHandler(w http.ResponseWriter, r *http.Request) {

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

func execsHandler(w http.ResponseWriter, r *http.Request) {

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


func main() {

	port := ":3000"
	cert := "cert.pem"
	key := "key.pem"

	mux := http.NewServeMux()
	rl := middlewares.NewRateLimiter(5, time.Minute)

	mux.HandleFunc("/", rootHandler)

	mux.HandleFunc("/teachers", teachersHandler)
	mux.HandleFunc("/students", studentsHandler)
	mux.HandleFunc("/execs", execsHandler)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	hppOptions := mw.HPPOptions{
		CheckQuery: true,
		CheckParams: true,
		CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
		Whitelist: []string{"sortBy", "name", "max"},
	}
	secureMux := mw.Hpp(hppOptions)(rl.Middleware(mw.CompressionHandler(mw.ResponseTimeMiddleware(mw.SecurityHeaders(mw.Cors(mux))))))

	server := &http.Server{
		Addr: port,
		Handler: secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running on port: ", port)
	err := server.ListenAndServeTLS(cert, key)	
	if err != nil {
		log.Fatalln("Error starting the server", err)
	}
	
}