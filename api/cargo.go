package api

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

type Cargo struct {
	ID          *int64  `json:"id"`
	Serial      *string `json:"serial"`
	Factory     *string `json:"factory"`
	Date        *string `json:"date"`
	Batch       *string `json:"batch"`
	Carrier     *string `json:"carrier"`
	ColdVan     *string `json:"cold_van"`
	Distributor *string `json:"distributor"`
	Hospital    *string `json:"hospital"`
	Patient     *string `json:"patient"`
}

const (
	migrate = `
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
)

func getDB(token string) *sql.DB {
	dbIns, ok := dbMap.Load(token)
	if !ok {
		logrus.Errorf("getDB: failed for %s", token)
		return nil
	}
	return dbIns.(*sql.DB)
}

func getCargos(token string) (result []Cargo) {
	var sql string
	var user string
	switch token {
	case getUserToken("user1"):
		sql = "SELECT * FROM cargos;"
		user = "user1"
	case getUserToken("user2"):
		sql = "SELECT id, serial, factory, date, batch FROM cargos;"
		user = "user2"
	case getUserToken("user3"):
		sql = "SELECT * FROM cargos;"
		user = "user3"
	default:
		logrus.Debugf("unknown token: %s", token)
		return
	}

	db := getDB(token)
	if db == nil {
		return
	}
	cargos, err := db.Query(sql)

	if err != nil {
		logrus.Errorf("getCargos: %v", err)
		return
	}

	logrus.Debugf("%s for %s", sql, user)
	defer cargos.Close()

	for cargos.Next() {
		cargo := Cargo{}
		switch user {
		case "user1":
			err = cargos.Scan(&cargo.ID, &cargo.Serial, &cargo.Factory, &cargo.Date, &cargo.Batch,
				&cargo.Carrier, &cargo.ColdVan, &cargo.Distributor, &cargo.Hospital, &cargo.Patient)
		case "user2":
			err = cargos.Scan(&cargo.ID, &cargo.Serial, &cargo.Factory, &cargo.Date, &cargo.Batch)
		case "user3":
			err = cargos.Scan(&cargo.ID, &cargo.Serial, &cargo.Factory, &cargo.Date, &cargo.Batch,
				&cargo.Carrier, &cargo.ColdVan, &cargo.Distributor, &cargo.Hospital, &cargo.Patient)
		default:
			panic("not possible user")
		}
		if err != nil {
			logrus.Errorf("getCargos: %v", err)
			return
		}
		if cargo.Date != nil {
			date := (*cargo.Date)[:10]
			cargo.Date = &date
		}

		result = append(result, cargo)
	}

	return
}

func postCargo(c echo.Context, token string) (id int64, err error) {
	var (
		cargo Cargo
		sql   string
		sets  string
		args  []interface{}
	)
	id = -1
	err = c.Bind(&cargo)
	if err != nil {
		logrus.Errorf("putCargo: %v", err)
		return
	}
	args = make([]interface{}, 0, 10)
	if cargo.Serial != nil {
		sets += "serial,"
		args = append(args, *cargo.Serial)
	}
	if cargo.Factory != nil {
		sets += "factory,"
		args = append(args, *cargo.Factory)
	}
	if cargo.Date != nil {
		sets += "date,"
		args = append(args, *cargo.Date)
	}
	if cargo.Batch != nil {
		sets += "batch,"
		args = append(args, *cargo.Batch)
	}
	if cargo.Carrier != nil {
		sets += "carrier,"
		args = append(args, *cargo.Carrier)
	}
	if cargo.ColdVan != nil {
		sets += "cold_van,"
		args = append(args, *cargo.ColdVan)
	}
	if cargo.Distributor != nil {
		sets += "distributor,"
		args = append(args, *cargo.Distributor)
	}
	if cargo.Hospital != nil {
		sets += "hospital,"
		args = append(args, *cargo.Hospital)
	}
	if cargo.Patient != nil {
		sets += "patient,"
		args = append(args, *cargo.Patient)
	}
	if len(args) == 0 {
		err = errors.New("nothing updated")
		logrus.Warn("putCargo: nothing updated")
		return
	}
	if len(sets) > 0 && sets[len(sets)-1] == ',' {
		sets = sets[:len(sets)-1]
	}

	placeHolders := strings.Repeat("?,", len(args))
	if len(placeHolders) > 0 && placeHolders[len(placeHolders)-1] == ',' {
		placeHolders = placeHolders[:len(placeHolders)-1]
	}

	sql = fmt.Sprintf("INSERT INTO cargos(%s) VALUES(%s)", sets, placeHolders)

	db := getDB(token)
	if db == nil {
		return
	}

	stmt, err := db.Prepare(sql)
	if err != nil {
		logrus.Errorf("postCargo: %v", err)
		return
	}

	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		logrus.Errorf("postCargo: %v", err)
		return
	}

	return result.LastInsertId()
}

func putCargo(c echo.Context, token string) (id int64, err error) {
	var (
		cargo Cargo
		sql   string
		sets  string
		args  []interface{}
	)
	id = -1
	err = c.Bind(&cargo)
	if err != nil {
		logrus.Errorf("putCargo: %v", err)
		return
	}
	if cargo.ID == nil {
		err = errors.New("no id got in req")
		logrus.Error("putCargo: no id got in req")
		return
	}
	args = make([]interface{}, 0, 10)
	if cargo.Serial != nil {
		sets += "serial = ?,"
		args = append(args, *cargo.Serial)
	}
	if cargo.Factory != nil {
		sets += "factory = ?,"
		args = append(args, *cargo.Factory)
	}
	if cargo.Date != nil {
		sets += "date = ?,"
		args = append(args, *cargo.Date)
	}
	if cargo.Batch != nil {
		sets += "batch = ?,"
		args = append(args, *cargo.Batch)
	}
	if cargo.Carrier != nil {
		sets += "carrier = ?,"
		args = append(args, *cargo.Carrier)
	}
	if cargo.ColdVan != nil {
		sets += "cold_van = ?,"
		args = append(args, *cargo.ColdVan)
	}
	if cargo.Distributor != nil {
		sets += "distributor = ?,"
		args = append(args, *cargo.Distributor)
	}
	if cargo.Hospital != nil {
		sets += "hospital = ?,"
		args = append(args, *cargo.Hospital)
	}
	if cargo.Patient != nil {
		sets += "patient = ?,"
		args = append(args, *cargo.Patient)
	}
	if len(args) == 0 {
		err = errors.New("nothing updated")
		logrus.Warn("putCargo: nothing updated")
		return
	}
	if len(sets) > 0 && sets[len(sets)-1] == ',' {
		sets = sets[:len(sets)-1]
	}

	sql = fmt.Sprintf("UPDATE cargos SET %s WHERE id = %d", sets, *cargo.ID)

	db := getDB(token)
	if db == nil {
		return
	}

	stmt, err := db.Prepare(sql)
	if err != nil {
		logrus.Errorf("putCargo: %v", err)
		return
	}

	defer stmt.Close()
	logrus.Debugf("putCargo: %s %v", sql, args)
	_, err = stmt.Exec(args...)
	if err != nil {
		logrus.Errorf("putCargo: %v", err)
		return
	}

	id = *cargo.ID

	return
}

func deleteCargo(c echo.Context, token string) (id int64, err error) {
	id = -1
	db := getDB(token)
	if db == nil {
		return
	}

	id, _ = strconv.ParseInt(c.Param("id"), 10, 64)

	sql := "DELETE FROM cargos WHERE id = ?"

	stmt, err := db.Prepare(sql)
	if err != nil {
		logrus.Errorf("deleteCargo: %v", err)
		return
	}

	result, err := stmt.Exec(id)
	if err != nil {
		logrus.Errorf("deleteCargo: %v", err)
		return
	}

	return result.RowsAffected()
}
