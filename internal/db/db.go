package db

import (
	"database/sql"
	"os"

	"github.com/lumenratio/finalTODO/internal/logger"
	_ "modernc.org/sqlite"
)

type DB struct {
	DBConn  *sql.DB
	DBFile  string
	install bool
}

const (
	schema string = `CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(128),
	comment TEXT,
	repeat VARCHAR(128));`
	index string = `CREATE INDEX idx_date ON scheduler (date);`
)

func (db *DB) createDB() error {
	var err error
	db.DBConn, err = sql.Open("sqlite", db.DBFile)

	if err != nil {
		return err
	}
	// Create table
	_, err = db.DBConn.Exec(schema)
	if err != nil {
		logger.Err.Println(err)
		return err
	}
	// Create index
	_, err = db.DBConn.Exec(index)
	if err != nil {
		return err
	}
	return nil
}

func InitDB(dbFile string) (*sql.DB, error) {
	db := DB{DBFile: dbFile}
	_, err := os.Stat(db.DBFile)

	if err != nil {
		logger.Info.Println("DB file not found. Will be create a new one")
		db.install = true
	}
	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX
	if db.install {
		err = db.createDB()
		if err != nil {
			logger.Err.Println("can't crate a new db:", err)
			return nil, err
		}
	}

	// Open already created DB
	db.DBConn, err = sql.Open("sqlite", db.DBFile)
	if err != nil {
		logger.Err.Println("can't open given DB:", err)
		return nil, err
	}

	return db.DBConn, err
}
