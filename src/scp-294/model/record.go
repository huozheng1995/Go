package model

import (
	"github.com/edward/scp-294/common"
)

type Record struct {
	Id         int
	Name       string
	GroupId    int
	InputType  string
	OutputType string
	InputData  string
	OutputData string
}

func ListRecords() (records []Record, err error) {
	sql := "SELECT Id, Name, GroupId, InputType, OutputType FROM record"
	rows, err := common.Db.Query(sql)
	if err != nil {
		return
	}
	for rows.Next() {
		record := Record{}
		err = rows.Scan(&record.Id, &record.Name, &record.GroupId, &record.InputType, &record.OutputType)
		if err != nil {
			return
		}

		records = append(records, record)
	}
	return
}

func GetRecord(id string) (record Record, err error) {
	sql := "SELECT Id, Name, GroupId, InputType, OutputType, InputData, OutputData FROM record WHERE Id=$1"
	err = common.Db.QueryRow(sql, id).Scan(&record.Id, &record.Name, &record.GroupId,
		&record.InputType, &record.OutputType, &record.InputData, &record.OutputData)
	return
}

func (record *Record) Insert() (err error) {
	sql := "INSERT INTO record (Name, GroupId, InputType, OutputType, InputData, OutputData) VALUES ($1, $2, $3, $4, $5, $6)"
	stmt, err := common.Db.Prepare(sql)
	if err != nil {
		return
	}
	_, err = stmt.Exec(record.Name, record.GroupId, record.InputType, record.OutputType, record.InputData, record.OutputData)
	if err != nil {
		return
	}
	return
}

func (record *Record) Update() (err error) {
	sql := "UPDATE record set Name=$1, GroupId=$2, InputType=$3, OutputType=$4, InputData=$5, OutputData=$6 WHERE Id=$7"
	_, err = common.Db.Exec(sql, record.Name, record.GroupId, record.InputType, record.OutputType,
		record.InputData, record.OutputData, record.Id)
	if err != nil {
		return
	}
	return
}

func DeleteRecord(id string) (err error) {
	sql := "DELETE FROM record WHERE Id=$1"
	_, err = common.Db.Exec(sql, id)
	if err != nil {
		return
	}
	return
}

func DeleteRecordsByGroupId(groupId string) (err error) {
	sql := "DELETE FROM record WHERE GroupId=$1"
	_, err = common.Db.Exec(sql, groupId)
	if err != nil {
		return
	}
	return
}
