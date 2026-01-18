package api

import (
	"encoding/json"
	"net/http"

	"github.com/lumenratio/finalTODO/internal/db"
	"github.com/lumenratio/finalTODO/internal/logger"
)

type ErrorResp struct {
	Error string `json:"error"`
}

type IDResp struct {
	ID string `json:"id"`
}

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

type JWTData struct {
	Token string `json:"token"`
}

func writeError(w http.ResponseWriter, message string, code int) {
	jerr := ErrorResp{Error: message}
	logger.Err.Println(jerr.Error)
	w.WriteHeader(code)
	writeJson(w, jerr)
}

func writeJson(w http.ResponseWriter, data any) {
	if data == nil {
		w.Write([]byte("{}")) // этот хак - самое быстрое, что смог придумать. Что бы прошли тесты и web клиент
		return
	}
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Err.Println("failed to encode JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
