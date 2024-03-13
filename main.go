package main

import (
	"database/sql"
	"log"

	"github.com/TechTrm/Authentication-Api-Services/api"
	db "github.com/TechTrm/Authentication-Api-Services/db/sqlc"
	_ "github.com/lib/pq"
)


const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:password@localhost:5432/users_db?sslmode=disable"
	serverAddress = "0.0.0.0:8080"

)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot Connect to DB:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
	

}