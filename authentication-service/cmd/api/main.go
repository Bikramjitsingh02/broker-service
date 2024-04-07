package main

import (
	data "command-line-arguments/home/bikramjit/micro-go/go-micro/authentication-service/data/models.go"
	"database/sql"
)

const WEB_PORT = "80"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {

}
