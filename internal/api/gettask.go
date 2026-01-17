package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/lumenratio/finalTODO/internal/db"
	"github.com/lumenratio/finalTODO/internal/logger"
)

func checkID(id string) (int, error) {
	if len(id) == 0 {
		return 0, fmt.Errorf("id field is empty")
	}
	num, err := strconv.Atoi(id)
	if num <= 0 {
		return 0, fmt.Errorf("given ID is zero or negtive number")
	}
	return num, err
}

func getTaskListHandler(DBConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println("request list of tasks")
		// set content-type
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// count of tasks in list
		num := 10
		data, err := db.GetTasks(DBConn, num)
		if err != nil {
			writeError(w, "can't get tasks from database", http.StatusInternalServerError)
			return
		}
		writeJson(w, TasksResp{Tasks: data})
	}
}

func getTaskHandler(DBConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println("get task by id")
		// set content-type
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// проверяем ID
		id, err := checkID(r.URL.Query().Get("id"))
		if err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Получаем задачу из базы
		task, err := db.GetTaskByID(DBConn, id)
		if err == sql.ErrNoRows {
			writeError(w, err.Error(), http.StatusNotFound)
			return
		}
		// остальные ошибки
		if err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJson(w, task)
	}
}
