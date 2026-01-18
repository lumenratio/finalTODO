package db

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TasksInfo struct {
	DB *sql.DB
	Task
}

func AddTask(task *TasksInfo) (int64, error) {
	var id int64
	// сам запрос
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := task.DB.Exec(query, sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat) /*передайте параметры task.Date, task.Title и т.д.*/)
	if err == nil {
		id, err = res.LastInsertId()
	}

	return id, err
}

func GetTasks(db *sql.DB, limit int) ([]*Task, error) {
	taskList := []*Task{}
	// запрос SELECT
	query := `SELECT * FROM scheduler ORDER BY date LIMIT :num`
	res, err := db.Query(query, sql.Named("num", limit))
	defer res.Close()
	if err != nil {
		return nil, err
	}

	// проверяем на пуcтой результат
	if res.Err() == sql.ErrNoRows {
		return taskList, err
	}

	// разбираем ответ
	for res.Next() {
		task := Task{}
		err := res.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		taskList = append(taskList, &task)
	}

	return taskList, err
}

func GetTaskByID(db *sql.DB, id int) (*Task, error) {
	// запрос SELECT
	query := `SELECT * FROM scheduler WHERE ID = :num`
	res := db.QueryRow(query, sql.Named("num", id))
	if res.Err() == sql.ErrNoRows {
		return nil, res.Err()
	}
	// разбираем ответ
	task := Task{}
	err := res.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}

	return &task, err
}

func UpdateTask(task *TasksInfo) error {
	// запрос на одновление
	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment,  repeat = :repeat WHERE id = :id`
	res, err := task.DB.Exec(query, sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat), sql.Named("id", task.ID))
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// был применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return err
}

func UpdateDate(task *TasksInfo) error {
	// запрос на одновление даты
	query := `UPDATE scheduler SET date = :date WHERE id = :id`
	_, err := task.DB.Exec(query, sql.Named("date", task.Date), sql.Named("id", task.ID))
	if err != nil {
		return err
	}
	return err
}

func DeleteTask(db *sql.DB, id int) error {
	//запрос на удаление
	query := `DELETE FROM scheduler WHERE ID = :id`
	res, err := db.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// был применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for delete task`)
	}

	return err
}
