package main

import (
	"database/sql"
	"log"

	"github.com/TechTrm/Authentication-Api-Services/api"
	db "github.com/TechTrm/Authentication-Api-Services/db/sqlc"
	"github.com/TechTrm/Authentication-Api-Services/util"
	_ "github.com/lib/pq"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot Connect to DB:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
	

}