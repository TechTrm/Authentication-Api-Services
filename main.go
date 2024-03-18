package main

import (
	"database/sql"
	"log"

	"github.com/TechTrm/Authentication-Api-Services/api"
	db "github.com/TechTrm/Authentication-Api-Services/db/sqlc"
	"github.com/TechTrm/Authentication-Api-Services/util"
	_ "github.com/lib/pq"
)

// func init(){

// 	// config, err := util.LoadConfig(".")
// 	// if err != nil {
// 	// 	log.Fatal("cannot load config:", err)
// 	// }

// 		// migrationURL  := config.MigrationURL
// 		// dbSource  :=  config.DBSource
// 		migrationURL  := "file://db/migration"
// 		dbSource  :=  "postgresql://root:password@postgres:5432/users_db?sslmode=disable"

// 		migration, err := migrate.New(migrationURL, dbSource)
// 		if err != nil {
// 			log.Fatal("cannot create new migrate instance")
// 		}

// 		if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
// 			log.Fatal("failed to run migrate up")
// 		}

// 		// log.Info("db migrated successfully")

// }

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
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
	

}

