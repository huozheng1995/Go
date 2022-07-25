package common

import (
	"context"
	"database/sql"
	"github.com/edward/scp-294/logger"
)

var Db *sql.DB

func Connect() error {
	var err error
	connStr := "./scp294.db"
	Db, err = sql.Open("sqlite3", connStr)
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = Db.PingContext(ctx)
	if err != nil {
		return err
	}
	logger.Log("Database Connected!")
	return nil
}

func InitDatabases(forceClean bool) {
	if forceClean {
		Db.Exec(`DROP TABLE IF EXISTS t_group`)
		Db.Exec(`DROP TABLE IF EXISTS record`)

		Db.Exec(`CREATE TABLE IF NOT EXISTS t_group(
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			Name VARCHAR(32) NOT NULL
		)`)
		Db.Exec(`INSERT INTO t_group (Id, Name) values (0, 'Default');`)
		Db.Exec(`CREATE TABLE IF NOT EXISTS record(
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			Name VARCHAR(32) NOT NULL,
			GroupId INTEGER,
			InputType VARCHAR(32) NOT NULL,
			OutputType VARCHAR(32) NOT NULL,
			InputData TEXT,
			OutputData TEXT
		)`)
		Db.Exec(`CREATE UNIQUE INDEX recordIndex ON record (GroupId, Name);`)
	}
}
