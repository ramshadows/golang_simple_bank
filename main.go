package main

import (
	"database/sql"
	"log"
	"simple_bank/api"
	db "simple_bank/db/sqlc"
	"simple_bank/utils"

	_ "github.com/lib/pq"
)

func main() {
	// load app configs
	config, err := utils.LoadConfigs(".") //"." means current folder

	if err != nil {
		log.Fatal("Cannot load app configs", err)
	}
	// Create a Postgres DB connection
	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}

	// create a new store and pass it the db connection
	store := db.NewStore(conn)

	// create a new server and pass the store
	server, err := api.NewServer(config, store)

	if err != nil {

		//log.Fatal("cannot create server:", err)
		log.Printf("%v cannot create server", err.Error())
	}

	// start the server by calling Start func and passing it the server address
	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Cannot start server. ", err)
	}

}
