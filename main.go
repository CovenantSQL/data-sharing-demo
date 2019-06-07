package main

import (
	"math/rand"
	"net/url"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/CovenantSQL/data-sharing-demo/api"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	rand.Seed(time.Now().UnixNano())

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
	e.POST("/apiv1/login", api.Login())
	e.POST("/apiv1/logout", api.Logout())
	e.GET("/apiv1/cargo", api.GetCargos())
	e.POST("/apiv1/cargo", api.PostCargo())
	e.PUT("/apiv1/cargo", api.PutCargo())
	e.DELETE("/apiv1/cargo/:id", api.DeleteCargo())
	e.POST("/apiv1/attach", api.Upload)
	e.GET("/apiv1/attach/:file", api.Download)
	e.GET("/apiv1/sql", api.GetCargoSql())
	e.Start(":8081")
}
