package common

import (
	"context"
	"database/sql"
	"github.com/edward/scp-294/logger"
)

var Db *sql.DB

func Connect() {
	var err error
	connStr := "./scp294.db"
	Db, err = sql.Open("sqlite3", connStr)
	if err != nil {
		logger.Log(err.Error())
		return
	}

	ctx := context.Background()
	err = Db.PingContext(ctx)
	if err != nil {
		logger.Log(err.Error())
		return
	}
	logger.Log("Database Connected!")
}

func CreateTableRecord() {
	drop := `DROP TABLE IF EXISTS record`
	_, err := Db.Exec(drop)
	if err != nil {
		logger.Log("Failed to drop table 'Record'")
	}

	create := `CREATE TABLE IF NOT EXISTS record(
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			Name VARCHAR(32) NOT NULL,
			ConvertType VARCHAR(32) NOT NULL,
			InputData TEXT,
			OutputData TEXT, 
			GroupId INTEGER
		)`
	_, err = Db.Exec(create)
	if err != nil {
		logger.Log("Failed to create table 'Record'")
	}

	createIndex := `CREATE UNIQUE INDEX recordIndex ON record (GroupId, Name);`
	_, err = Db.Exec(createIndex)
	if err != nil {
		logger.Log("Failed to index of table 'Record'")
	}
}
