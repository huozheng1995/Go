package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
)

var db *sql.DB

const (
	server   = "172.16.85.138"
	port     = 1433
	user     = "sa"
	password = "xA123456"
	database = "test"
)

type app struct {
	id   int
	name string
}

func main() {
	connStr := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)
	var err error
	db, err = sql.Open("sqlserver", connStr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println("Connected!")

	fmt.Println(getMany(0))

	a := app{
		id:   3,
		name: "insert 3",
	}
	//a.Insert()
	a.PrepareInsert()
	fmt.Println(getMany(0))

	a = getOne(3)
	a.name = "update 1"
	a.Update()
	fmt.Println(getOne(3))

	a.Delete()
	fmt.Println(getMany(0))

}

func getOne(id int) (a app) {
	a = app{}
	err := db.QueryRow("select id, name from t_test where id = @id",
		sql.Named("id", id)).Scan(&a.id, &a.name)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return
}

func getMany(id int) (apps []app) {
	rows, err := db.Query("select id, name from t_test where id > @id",
		sql.Named("id", id))
	if err != nil {
		log.Fatalln(err.Error())
	}
	for rows.Next() {
		a := app{}
		err = rows.Scan(&a.id, &a.name)
		if err != nil {
			log.Fatalln(err.Error())
		}
		apps = append(apps, a)
	}
	return
}

func (a *app) Update() {
	_, err := db.Exec("update t_test set name=@name where id=@id",
		sql.Named("name", a.name), sql.Named("id", a.id))
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func (a *app) Delete() {
	_, err := db.Exec("delete from t_test where id=@id",
		sql.Named("id", a.id))
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func (a *app) Insert() {
	_, err := db.Exec("insert into t_test(id, name) values(@id, @name)",
		sql.Named("id", a.id), sql.Named("name", a.name))
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func (a *app) PrepareInsert() {
	stmt, err := db.Prepare(`insert into t_test(id, name) values(@id, @name);
		select isNull(SCOPE_IDENTITY(), -1)`)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer stmt.Close()
	row := stmt.QueryRow(sql.Named("id", a.id), sql.Named("name", a.name))
	err = row.Scan(&a.id)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
