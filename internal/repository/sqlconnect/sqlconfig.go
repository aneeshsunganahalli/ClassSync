package sqlconnect

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDb() (*sql.DB, error) {
	fmt.Println("Trying to connect to MariaDB")

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("HOST")
	if host == "" {
		host = os.Getenv("DB_HOST")
	}
	dbname := os.Getenv("DB_NAME")
	dbport := os.Getenv("DB_PORT")

	if user == "" || password == "" || host == "" || dbname == "" || dbport == "" {
		return nil, errors.New("missing required database environment variables")
	}

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", user, password, host, dbport, dbname)
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	fmt.Println("Connected to MariaDB")
	return db, nil
}
