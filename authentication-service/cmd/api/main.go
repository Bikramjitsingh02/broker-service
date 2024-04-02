package main

import (
	"authentication/data"
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

const WEB_PORT = "81"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {

	log.Println("starting authentication service on port : 81")

	conn := ConnectToDB()

	if conn == nil {
		log.Panic("cant connect to postGress")
	}

	//set up config

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", WEB_PORT),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Panic("Something gone wrong ", err)
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
			counts++
			log.Println("Postgres not yet ready ......")
		} else {
			log.Println("Connected to postGres")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds........")
		time.Sleep(2 * time.Second)
		continue
	}

}
