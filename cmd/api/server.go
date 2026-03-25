package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aneeshsunganahalli/ClassSync/internal/api/middlewares"
	mw "github.com/aneeshsunganahalli/ClassSync/internal/api/middlewares"
	"github.com/aneeshsunganahalli/ClassSync/internal/api/router"
	"github.com/aneeshsunganahalli/ClassSync/internal/repository/sqlconnect"
	"github.com/aneeshsunganahalli/ClassSync/pkg/utils"

	"github.com/joho/godotenv"
)


func main() {

	err := godotenv.Load()
	if err != nil {
		return
	}
	
	_, err = sqlconnect.ConnectDb()
	if err != nil {
		fmt.Println("Error")
		return
	}

	port := ":" + os.Getenv("API_PORT")
	cert := "cert.pem"
	key := "key.pem"

	rl := middlewares.NewRateLimiter(10, time.Minute)


	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	hppOptions := mw.HPPOptions{
		CheckQuery: true,
		CheckParams: true,
		CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
		Whitelist: []string{"sortBy", "name", "max", "first_name", "last_name"},
	}

	// Maintain logical order of when to apply which middleware
	// secureMux :=  mw.Cors(rl.Middleware(mw.ResponseTimeMiddleware(mw.SecurityHeaders(mw.CompressionHandler(mw.Hpp(hppOptions)(mux))))))

	router := router.Router()
	secureMux := utils.ApplyMiddlewares(router, mw.Cors, rl.Middleware, mw.ResponseTimeMiddleware, mw.SecurityHeaders, mw.CompressionHandler, mw.Hpp(hppOptions))

	server := &http.Server{
		Addr: port,
		Handler: secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running on port: ", port)
	err = server.ListenAndServeTLS(cert, key)	
	if err != nil {
		log.Fatalln("Error starting the server", err)
	}
	
}



