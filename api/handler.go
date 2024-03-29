package api

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
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
	db   = "966c4a297ff996772ef1b7f7ba13f1a1fab95bb4589f9e2236d37cbc83ac38bd"
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
			logrus.Errorf("Login: bad request, %v", err)
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
				logrus.Errorf("Login: open db failed, %v", err)
				return c.JSON(http.StatusInternalServerError, err)
			}
			err = db.Ping()
			if err != nil {
				_ = db.Close()
				logrus.Errorf("Login: ping db failed, %v", err)
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
			logrus.Error("Logout: no token got")
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
			logrus.Error("GetCargos: no token got")
			return c.JSON(http.StatusBadRequest, nil)
		} else {
			return c.JSON(http.StatusOK, getCargos(token))
		}
	}
}

func PutCargo() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			logrus.Error("PutCargos: no token got")
			return c.JSON(http.StatusBadRequest, nil)
		} else {
			if sql, id, err := putCargo(c, token); err == nil {
				return c.JSON(http.StatusOK, H{
					"updated": id,
					"cql":     sql,
				})
			} else {
				return c.JSON(http.StatusForbidden, err)
			}
		}
	}
}

func PostCargo() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			logrus.Error("PutCargos: no token got")
			return c.JSON(http.StatusBadRequest, nil)
		} else {
			if sql, id, err := postCargo(c, token); err == nil {
				return c.JSON(http.StatusCreated, H{
					"created": id,
					"cql":     sql,
				})
			} else {
				return c.JSON(http.StatusForbidden, err)
			}
		}
	}
}

func DeleteCargo() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			logrus.Error("deleteCargo: no token got")
			return c.JSON(http.StatusBadRequest, nil)
		} else {
			if sql, id, err := deleteCargo(c, token); err == nil {
				return c.JSON(http.StatusOK, H{
					"deleted": id,
					"cql":     sql,
				})
			} else {
				return c.JSON(http.StatusForbidden, err)
			}
		}
	}
}

func GetCargoSql() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			logrus.Error("GetCargoSql: no token got")
			return c.JSON(http.StatusBadRequest, nil)
		} else {
			return c.JSON(http.StatusOK, getCargoSql(c, token))
		}
	}
}
