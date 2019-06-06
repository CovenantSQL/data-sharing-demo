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
	e.GET("/apiv1/cargos", api.GetCargos())
	e.POST("/apiv1/cargo", api.PostCargo())
	e.PUT("/apiv1/cargo", api.PutCargo())
	e.POST("/apiv1/upload", api.Upload)
	e.DELETE("/apiv1/cargo/:id", api.DeleteCargo())
	e.Start(":8081")
}
