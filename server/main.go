package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/punkzberryz/todo/api"
	db "github.com/punkzberryz/todo/db/sqlc"
	"github.com/punkzberryz/todo/session"
	"github.com/punkzberryz/todo/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	dbConn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(dbConn)
	sessionConn, err := session.NewSession(config.RedisAddress)
	if err != nil {
		log.Fatal("cannot connect to redis client:", err)
	}

	server, err := api.NewServer(config, &store, &sessionConn)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	addr := fmt.Sprintf(":%s", config.ServerPort)
	log.Printf("start listening to: http://%s", config.ServerAddress)
	http.ListenAndServe(addr, server.Router)
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance:", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up:", err)
	}
	log.Println("db migrated successfully")
}
