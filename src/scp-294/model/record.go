package model

import (
	"github.com/edward/scp-294/common"
)

type Record struct {
	Id          int
	Name        string
	ConvertType string
	InputData   string
	OutputData  string
	GroupId     int
}

func ListRecords() (records []Record, err error) {
	sql := "SELECT Id, Name, ConvertType, InputData, OutputData, GroupId FROM record"
	rows, err := common.Db.Query(sql)
	if err != nil {
		return
	}
	for rows.Next() {
		record := Record{}
		err = rows.Scan(&record.Id, &record.Name, &record.ConvertType,
			&record.InputData, &record.OutputData, &record.GroupId)
		if err != nil {
			return
		}

		records = append(records, record)
	}
	return
}

func GetRecord(id string) (record Record, err error) {
	sql := "SELECT Id, Name, ConvertType, InputData, OutputData, GroupId FROM record WHERE Id=$1"
	err = common.Db.QueryRow(sql, id).Scan(&record.Id, &record.Name, &record.ConvertType,
		&record.InputData, &record.OutputData, &record.GroupId)
	return
}

func (record *Record) Insert() (err error) {
	sql := "INSERT INTO record (Name, ConvertType, InputData, OutputData, GroupId) VALUES ($1, $2, $3, $4, $5)"
	stmt, err := common.Db.Prepare(sql)
	if err != nil {
		return
	}
	_, err = stmt.Exec(record.Name, record.ConvertType, record.InputData, record.OutputData, record.GroupId)
	if err != nil {
		return
	}
	return
}

func (record *Record) Update() (err error) {
	sql := "UPDATE record set Name=$1, ConvertType=$2, InputData=$3, OutputData=$4, GroupId=$5 WHERE Id=$6"
	_, err = common.Db.Exec(sql, record.Name, record.ConvertType, record.InputData, record.OutputData,
		record.GroupId, record.Id)
	return
}

func DeleteRecord(id int) (err error) {
	sql := "DELETE FROM record WHERE Id=$1"
	_, err = common.Db.Exec(sql, id)
	return
}
