package sqlconnect

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDb() (*sql.DB, error) {
	fmt.Println("Trying to connect to MariaDB")

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("HOST")
	dbname := os.Getenv("DB_NAME")
	dbport := os.Getenv("DB_PORT")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, dbport, dbname)
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		// panic(err)
		return nil, err
	}

	fmt.Println("Connected to MariaDB")
	return db, nil
}