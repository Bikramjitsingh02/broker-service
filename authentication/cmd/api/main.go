package main

import (
	"authetication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const WEB_PORT = "80"

var count int64

type Config struct {
	Models data.Models
	DB     *sql.DB
}

func main() {
	log.Println("Starting authentication service ")

	conn := ConnectToDB()

	if conn == nil {
		log.Panic("cant connect to Postgres !!")
	}

	//set up config

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", WEB_PORT),
		Handler: app.routes(),
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Panic("Something went wrong : ", err)
	}

}

func OpenDB(dsn string) (*sql.DB, error) {

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ConnectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	for {
		connection, err := OpenDB(dsn)
		if err != nil {
			log.Println("Postgres is not ready yet")
			count++
		} else {
			log.Println("Connected to Postgres")
			return connection
		}

		if count > 10 {
			log.Println(err)
			return nil
		}

		log.Println("backing off for two seconds .......")
		time.Sleep(2 * time.Second)
		continue
	}
}
