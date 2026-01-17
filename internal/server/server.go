package server

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"

	"github.com/lumenratio/finalTODO/internal/api"
	"github.com/lumenratio/finalTODO/internal/logger"
)

const defaultPort string = "7540"

//type Server struct {
//	Server http.Server
//}

// check that PORT have only numbers
func validatePort(portEnv string) string {
	// check string
	port, err := strconv.Atoi(portEnv)
	if err == nil && port >= 1024 && port <= 65535 {
		return strconv.Itoa(port)
	}
	if err != nil {
		logger.Err.Println("in TODO_PORT env only numbers is allowed. Use standard port number")
		return defaultPort
	}
	logger.Info.Println("in TODO_PORT env only numbers in 1024-65535 range is allowed. Use standard port number")
	return defaultPort
}

// Here we initialise HTTP server
func InitServer(db *sql.DB) *http.Server {
	// validate port string
	port := defaultPort
	if portStr := os.Getenv("TODO_PORT"); portStr != "" {
		port = validatePort(portStr)
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:     ":" + port,
		Handler:  mux,
		ErrorLog: logger.Err,
		//ReadTimeout:  (time.Second * 5),
		//WriteTimeout: (time.Second * 10),
		//IdleTimeout:  (time.Second * 15),
	}

	//mux.HandleFunc("POST /upload", handlers.Upload)
	// Register hadlers from API package
	api.Init(mux, db)
	return server
}
