package model

import "github.com/edward/scp-294/common"

type Group struct {
	Id   int
	Name string
}

func ListGroups() (groups []Group, err error) {
	sql := "SELECT Id, Name FROM t_group"
	rows, err := common.Db.Query(sql)
	if err != nil {
		return
	}
	for rows.Next() {
		group := Group{}
		err = rows.Scan(&group.Id, &group.Name)
		if err != nil {
			return
		}

		groups = append(groups, group)
	}
	return
}

func GetGroup(id string) (group Group, err error) {
	sql := "SELECT Id, Name FROM t_group WHERE Id=$1"
	err = common.Db.QueryRow(sql, id).Scan(&group.Id, &group.Name)
	return
}

func (group *Group) Insert() (err error) {
	sql := "INSERT INTO t_group (Name) VALUES ($1)"
	stmt, err := common.Db.Prepare(sql)
	if err != nil {
		return
	}
	_, err = stmt.Exec(group.Name)
	if err != nil {
		return
	}
	return
}

func (group *Group) Update() (err error) {
	sql := "UPDATE t_group set Name=$1 WHERE Id=$2"
	_, err = common.Db.Exec(sql, group.Name, group.Id)
	if err != nil {
		return
	}
	return
}

func DeleteGroup(id string) (err error) {
	sql := "DELETE FROM t_group WHERE Id=$1"
	_, err = common.Db.Exec(sql, id)
	if err != nil {
		return
	}
	return
}
