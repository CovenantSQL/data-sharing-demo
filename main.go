package main

import (
	"math/rand"
	"net/url"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/CovenantSQL/data-sharing-demo/model"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	//db := initDb("storage.db")
	//migrate(db)

	// debug proxy
	url1, err := url.Parse("http://localhost:8080")
	if err != nil {
		logrus.Fatal(err)
	}
	targets := []*middleware.ProxyTarget{
		{
			URL: url1,
		},
	}
	e := echo.New()
	g := e.Group("/")
	g.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
	e.Use(middleware.CORS())
	//e.Static("/", "frontend")
	e.POST("/apiv1/login", model.Login())
	e.POST("/apiv1/logout", model.Logout())
	e.GET("/apiv1/cargos", model.GetCargos())
	e.POST("/apiv1/cargo", model.PostCargo())
	e.PUT("/apiv1/cargo", model.PutCargo())
	e.DELETE("/apiv1/cargo/:id", model.DeleteCargo())
	e.Start(":8081")
}
