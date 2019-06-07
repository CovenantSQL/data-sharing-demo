package api

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

type SqlRecord struct {
	Id      *int64  `json:"id"`
	SQL     *string `json:"sql"`
	CargoId *int64  `json:"cargo_id"`
	Hash    *string `json:"hash"`
	User    *string `json:"user"`
}

const (
	SqlMigrate = `
		CREATE TABLE sqls (
			id INTEGER PRIMARY KEY,
			sql VARCHAR,
			cargo_id VARCHAR,
			hash VARCHAR,
			user VARCHAR
		);
	`
)

func getCargoSql(c echo.Context, token string) (result []SqlRecord) {

	db := getDB(token)
	if db == nil {
		return
	}

	cId, _ := strconv.Atoi(c.QueryParam("cargo_id"))
	where := ""
	if cId != 0 {
		where = fmt.Sprintf("WHERE cargo_id = %d", cId)
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 20
	}
	sql := fmt.Sprintf("SELECT * FROM sqls %s ORDER BY id DESC LIMIT %d;", where, limit)

	rows, err := db.Query(sql)

	if err != nil {
		logrus.Errorf("getCargoSql: %v", err)
		return
	}

	defer rows.Close()
	logrus.Debugf("getCargoSql: %s", sql)

	for rows.Next() {
		record := SqlRecord{}
		err = rows.Scan(&record.Id, &record.SQL, &record.CargoId, &record.Hash, &record.User)
		if err != nil {
			logrus.Errorf("getCargoSql: %v", err)
			return
		}

		result = append(result, record)
	}

	return
}

func putCargoSql(cId int64, sql string, user string, db *sql.DB) (err error) {
	if db == nil {
		return fmt.Errorf("db is nil")
	}
	var sqlHash string
	err = db.QueryRow("SELECT last_request_hash();").Scan(&sqlHash)
	if err != nil {
		logrus.Errorf("SELECT last_request_hash failed: %v", err)
		return
	}
	insSql := "INSERT INTO sqls(sql, cargo_id, hash, user) VALUES (?, ?, ?, ?);"
	stmt, err := db.Prepare(insSql)
	if err != nil {
		logrus.Errorf("putCargoSql: %v", err)
		return
	}

	defer stmt.Close()
	logrus.Debugf("putCargoSql: %s %d %s\n\t %s %s", user, cId, insSql, sql, sqlHash)
	_, err = stmt.Exec(sql, cId, sqlHash, user)
	if err != nil {
		logrus.Errorf("putCargoSql: %v", err)
		return
	}
	return
}
