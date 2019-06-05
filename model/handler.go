package model

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

type H map[string]interface{}

type authReq struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type authResp struct {
	User  string `json:"user"`
	Token string `json:"token"`
}

const (
	salt = "9e1618775ccdeeb19f110e04fbc6c5bc"
	db   = "b77422b30688fdc8facfe84a0c48c1f94aca3444178a9502753b3692a5576f10"
	host = "127.0.0.1:4665"
)

var (
	dbMap sync.Map //map[key]*sql.DB
)

func getUserToken(user string) string {
	buf := make([]byte, 0, len(user)+len(salt))
	buf = append(buf, []byte(user)...)
	buf = append(buf, []byte(salt)...)
	sha := sha256.Sum256(buf)
	return hex.EncodeToString(sha[:16])
}

func Login() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		auth := new(authReq)
		if err = c.Bind(auth); err != nil {
			logrus.Errorf("Login: bad request, %s", err)
			return c.JSON(http.StatusBadRequest, err)
		}
		token := getUserToken(auth.User)
		_, ok := dbMap.Load(token)
		if !ok {
			dsn := fmt.Sprintf(
				"%s:%s@tcp(%s)/%s",
				auth.User, auth.Password, host, db)
			db, err := sql.Open("mysql", dsn)
			if err != nil {
				logrus.Errorf("Login: open db failed, %s", err)
				return c.JSON(http.StatusInternalServerError, err)
			}
			err = db.Ping()
			if err != nil {
				_ = db.Close()
				logrus.Errorf("Login: ping db failed, %s", err)
				return c.JSON(http.StatusUnauthorized, err)
			}
			_, loaded := dbMap.LoadOrStore(token, db)
			if loaded {
				_ = db.Close()
			}
		}

		resp := authResp{
			User:  auth.User,
			Token: token,
		}
		logrus.Debugf("Login: succeed, %s", auth.User)

		return c.JSON(http.StatusOK, resp)
	}
}

func Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			logrus.Errorf("Logout: no token got")
			return c.JSON(http.StatusBadRequest, nil)
		} else {
			dbMap.Delete(token)
			return c.JSON(http.StatusOK, nil)
		}
	}
}

func GetCargos() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			logrus.Errorf("GetCargos: no token got")
			return c.JSON(http.StatusBadRequest, nil)
		} else {
			return c.JSON(http.StatusOK, getCargos(token))
		}
	}
}

func PostCargo(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var cargo Cargo
		_ = c.Bind(&cargo)

		id, err := postCargo(db, cargo)

		if err == nil {
			return c.JSON(http.StatusCreated, H{
				"created": id,
			})
		} else {
			return err
		}

	}
}

func PutCargo() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			logrus.Errorf("PutCargos: no token got")
			return c.JSON(http.StatusBadRequest, nil)
		} else {
			if id, err := putCargo(c, token); err == nil {
				return c.JSON(http.StatusOK, H{
					"updated": id,
				})
			} else {
				return err
			}

		}

	}
}

func DeleteCargo(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		_, err := deleteCargo(db, id)

		if err == nil {
			return c.JSON(http.StatusOK, H{
				"deleted": id,
			})
		} else {
			return err
		}

	}
}
