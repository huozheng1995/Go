package drivers

import (
	"database/sql"
	"fmt"
	_ "github.com/alexbrainman/odbc"
	"log"
	"testing"
)

func TestDr35970(t *testing.T) {
	db, err := sql.Open("odbc", "DSN=CData DB2 Sys")
	if err != nil {
		log.Fatal(err)
	}

	var (
		DbInsert string
		//DbUpdate string
	)

	DbInsert = `INSERT INTO SCHEMA1."StarChecks" ("first_name", "last_name", "planet", "id") VALUES (?, ?, ?, ?);`
	//DbUpdate = `UPDATE SCHEMA1."StarChecks" SET "first_name" = ?, "last_name" = ?, "planet" = ? WHERE "id" = ?;`

	st, err := db.Prepare(DbInsert)
	if err != nil {
		log.Fatal(err)
	}

	res, err := st.Exec("Ahsoka", "Tano", "Some Colony", "12")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.RowsAffected())

	defer st.Close()
	defer db.Close()
}

func TestDr34449(t *testing.T) {
	db, err := sql.Open("odbc", "DSN=CData DB2 Sys")
	if err != nil {
		log.Fatal(err)
	}

	var (
		DbMergeQuery3 string
		//DbInsert string
		//DbUpdate string
	)

	DbMergeQuery3 = `MERGE INTO SCHEMA1."StarChecks" AS target 
USING ( VALUES (?, ?, ?, ?) ) AS source (first_name, last_name, planet, id) ON target."id" = source.id AND target."first_name" = source.first_name 
WHEN MATCHED THEN UPDATE SET "first_name" = source.first_name, "last_name" = source.last_name, "planet" = source.planet 
WHEN NOT MATCHED THEN INSERT ("id", "first_name", "last_name", "planet") VALUES (source.id, source.first_name, source.last_name, source.planet);`
	//DbInsert = `INSERT INTO SCHEMA1."StarChecks" ("first_name", "last_name", "planet", "id") VALUES (?, ?, ?, ?);`
	//DbUpdate = `UPDATE SCHEMA1."StarChecks" SET "first_name" = ?, "last_name" = ?, "planet" = ? WHERE "id" = ?;`

	st, err := db.Prepare(DbMergeQuery3)
	if err != nil {
		log.Fatal(err)
	}

	res, err := st.Exec("Ahsoka", "Tano", "Some Colony", "11")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.RowsAffected())

	defer st.Close()
	defer db.Close()
}

func TestDr34449_insert(t *testing.T) {
	db, err := sql.Open("odbc", "DSN=CData DB2 Sys")
	if err != nil {
		log.Fatal(err)
	}

	var (
		DbInsert string
		//DbUpdate string
	)

	DbInsert = `INSERT INTO SCHEMA1."StarChecks" ("first_name", "last_name", "planet", "id") VALUES (?, ?, ?, ?);`
	//DbUpdate = `UPDATE SCHEMA1."StarChecks" SET "first_name" = ?, "last_name" = ?, "planet" = ? WHERE "id" = ?;`

	st, err := db.Prepare(DbInsert)
	if err != nil {
		log.Fatal(err)
	}

	res, err := st.Exec("Ahsoka", "Tano", "Some Colony", "12")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.RowsAffected())

	defer st.Close()
	defer db.Close()
}
