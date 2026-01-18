package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/lumenratio/finalTODO/internal/db"
	"github.com/lumenratio/finalTODO/internal/logger"
)

func checkTitle(title string) bool {
	var result bool
	if len(title) == 0 {
		result = true
	}
	return result
}

func returnNowTime() time.Time {
	return time.Now().Truncate(24 * time.Hour)
}

func checkDate(task *db.Task) error {
	var err error
	// проверяем на пустое поле
	now := returnNowTime()
	tmpDate := task.Date
	if len(tmpDate) == 0 {
		task.Date = now.Format(dateFormat)
	}

	// Поле не пустое, проверяем формат
	var tt time.Time
	if len(tmpDate) != 0 {
		tt, err = time.Parse(dateFormat, task.Date)
		if err != nil {
			return err
		}
	}

	// проверяем Repeat, Date и выполняем вычисление даты
	if now.After(tt) {
		if len(task.Repeat) != 0 {
			task.Date, err = NextDate(now.Format(dateFormat), task.Date, task.Repeat)
			if err != nil {
				return err
			}
		}
		task.Date = now.Format(dateFormat)
	}
	return err
}

func manageTaskHandler(DBConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println("new task arrived")
		// set content-type
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Какой метод сейчас
		method := r.Method
		// инициализируем структуру
		task := db.TasksInfo{DB: DBConn}
		//Read body
		var buff bytes.Buffer
		_, err := buff.ReadFrom(r.Body)
		if err != nil {
			writeError(w, error.Error(err), http.StatusInternalServerError)
			return
		}
		// десериализуем JSON
		if err = json.Unmarshal(buff.Bytes(), &task.Task); err != nil {
			writeError(w, "got malformed JSON in request", http.StatusBadRequest)
			return
		}

		// Проверяем поле Title
		if checkTitle(task.Title) {
			writeError(w, "got malformed JSON in request: \"Title\" can't be empty", http.StatusBadRequest)
			return
		}

		//Проверяем дату
		if err = checkDate(&task.Task); err != nil {
			writeError(w, "got malformed JSON in request: \"Date\" format not supported", http.StatusBadRequest)
			return
		}

		if method == "POST" {
			// Добавляем таску в базу
			id := IDResp{}
			taskID, err := db.AddTask(&task)
			if err != nil {
				writeError(w, "can't add task into database", http.StatusInternalServerError)
				return
			}
			id.ID = strconv.FormatInt(taskID, 10)
			// Возвращаем ID
			writeJson(w, id)
		}

		if method == "PUT" {
			// Обновляем таску в базе
			err = db.UpdateTask(&task)
			if err != nil {
				writeError(w, "can't add task into database", http.StatusBadRequest)
				return
			}
			writeJson(w, nil)
		}
	}
}

func doneTaskHandler(DBConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println("do close/promote task by id")
		// set content-type
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// проверяем ID
		id, err := checkID(r.URL.Query().Get("id"))
		if err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Выгружаем задачу из базы
		task, err := db.GetTaskByID(DBConn, id)
		if err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}

		// проверяем на наличие поля repeat
		if len(task.Repeat) == 0 {
			err = db.DeleteTask(DBConn, id)
			if err != nil {
				writeError(w, err.Error(), http.StatusBadRequest)
				return
			}
			writeJson(w, db.Task{})
			return
		}

		// OK, repeat не пустое поле
		task.Date, err = NextDate(returnNowTime().Format(dateFormat), task.Date, task.Repeat)
		if err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}
		tasksInf := db.TasksInfo{DB: DBConn, Task: *task}
		err = db.UpdateDate(&tasksInf)
		if err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}

		writeJson(w, db.Task{})
	}
}

func deleteTaskHandler(DBConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println("got delete task by id request")
		// set content-type
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// проверяем ID и удаляем задачу
		id, err := checkID(r.URL.Query().Get("id"))
		if err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Удаляем задачу из базы
		err = db.DeleteTask(DBConn, id)
		if err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJson(w, nil)
	}
}
