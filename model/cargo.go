package model

import (
	"database/sql"

	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

type Cargo struct {
	ID          int    `json:"id"`
	Serial      string `json:"serial"`
	Factory     string `json:"factory"`
	Date        string `json:"date"`
	Batch       string `json:"batch"`
	Carrier     string `json:"carrier"`
	ColdVan     string `json:"cold_van"`
	Distributor string `json:"distributor"`
	Hospital    string `json:"hospital"`
	Patient     string `json:"patient"`
}

func migrate(db *sql.DB) {
	sql := `
		CREATE TABLE cargos (
			id INTEGER PRIMARY KEY,
			serial VARCHAR,
			factory VARCHAR,
			date TIMESTAMP,
			batch VARCHAR,
			carrier VARCHAR,
			cold_van VARCHAR,
			distributor VARCHAR,
			hospital VARCHAR,
			patient VARCHAR
		);
	`

	_, err := db.Exec(sql)
	if err != nil {
		panic(err)
	}
}

func getCargos(token string) (result []Cargo) {
	dbIns, ok := dbMap.Load(token)
	if !ok {
		return result
	}
	db := dbIns.(*sql.DB)
	var sql string
	var user string
	switch token {
	case getUserToken("user1"):
		sql = "SELECT * FROM cargos;"
		user = "user1"
	case getUserToken("user2"):
		sql = "SELECT * FROM cargos;"
		user = "user2"
	case getUserToken("user3"):
		sql = "SELECT * FROM cargos;"
		user = "user3"
	default:
		logrus.Debugf("unknown token: %s", token)
		return
	}
	cargos, err := db.Query(sql)

	if err != nil {
		logrus.Errorf("getCargos: %s", err)
		return
	}

	logrus.Debugf("%s for %s", sql, user)
	defer cargos.Close()

	for cargos.Next() {
		cargo := Cargo{}
		err := cargos.Scan(&cargo.ID, &cargo.Serial, &cargo.Factory, &cargo.Date, &cargo.Batch,
			&cargo.Carrier, &cargo.ColdVan, &cargo.Distributor, &cargo.Hospital, &cargo.Patient)
		if err != nil {
			logrus.Errorf("getCargos: %s", err)
			return
		}
		cargo.Date = cargo.Date[:10]
		result = append(result, cargo)
	}

	return
}

func postCargo(db *sql.DB, cargo Cargo) (int64, error) {
	sql := "INSERT INTO cargos(name, done) VALUES(?, 0)"

	stmt, err := db.Prepare(sql)
	if err != nil {
		panic(err)
	}

	defer stmt.Close()

	result, err := stmt.Exec(name)
	if err != nil {
		panic(err)
	}

	return result.LastInsertId()
}

func putCargo(c echo.Context, token string) (int64, error) {
	var cargo Cargo
	_ = c.Bind(&cargo)

	var sql string

	sql = "UPDATE cargos SET name = ?, done = ? WHERE id = ?"

	stmt, err := db.Prepare(sql)
	if err != nil {
		panic(err)
	}

	defer stmt.Close()

	done := 0
	if cargo.Done {
		done = 1
	}

	result, err := stmt.Exec(cargo.Name, done, cargo.ID)

	if err != nil {
		panic(err)
	}

	return result.LastInsertId()
}

func deleteCargo(db *sql.DB, id int) (int64, error) {
	sql := "DELETE FROM cargos WHERE id = ?"

	stmt, err := db.Prepare(sql)
	if err != nil {
		panic(err)
	}

	result, err := stmt.Exec(id)
	if err != nil {
		panic(err)
	}

	return result.RowsAffected()
}
