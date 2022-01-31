package main

import (
	"fmt"
	"todoproject/database"
	"todoproject/route"
)

const (
	host     = "localhost"
	dbname   = "todo"
	port     = "5432"
	user     = "postgres"
	password = "1234"
)

func main() {
	err := database.Connect(host, port, dbname, user, password, database.SSLModeDisable)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected....")
	srv := route.Route()
	connErr := srv.Run(":5555")
	if connErr != nil {
		panic(err)
	}

}
