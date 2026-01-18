package main

import (
	"os"

	"github.com/lumenratio/finalTODO/internal/db"
	"github.com/lumenratio/finalTODO/internal/logger"
	"github.com/lumenratio/finalTODO/internal/server"
)

func main() {
	//DB code block
	dbFile := "scheduler.db"
	dbEnv := os.Getenv("TODO_DBFILE")
	if dbEnv != "" {
		dbFile = dbEnv
	}

	db, err := db.InitDB(dbFile)
	if err != nil {
		os.Exit(1)
	}
	defer db.Close()

	// Http code block

	// Init HTTP server
	srv := server.InitServer(db)

	// Start a server
	logger.Info.Printf("Starting http server on address %s\n", srv.Addr)
	err = srv.ListenAndServe()
	if err != nil {
		srv.ErrorLog.Fatal("error when try start the server: ", err)
	}
}
