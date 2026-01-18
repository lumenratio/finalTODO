package internal

import (
	"net/http"

	"github.com/lumenratio/finalTODO/internal/logger"
)

func IndexFS(w http.ResponseWriter, r *http.Request) {
	rootDir := http.Dir("./web")
	logger.Info.Println(rootDir)
	http.FileServer(rootDir)
}
