package main

import (
	"fmt"
	"todoproject/database"
	"todoproject/route"
	"todoproject/utils"
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
	utils.CheckError(err)
	fmt.Println("Connected....")
	srv := route.Route()
	connErr := srv.Run(":5555")
	utils.CheckError(connErr)

}
