package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/punkzberryz/todo/api"
	db "github.com/punkzberryz/todo/db/sqlc"
	"github.com/punkzberryz/todo/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, &store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	addr := fmt.Sprintf(":%s", config.ServerPort)
	http.ListenAndServe(addr, server.Router)
}
