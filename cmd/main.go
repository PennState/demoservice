package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	echoLog "github.com/neko-neko/echo-logrus/v2/log"

	log "github.com/sirupsen/logrus"
)

func main() {
	e := echo.New()

	// Middleware

	e.Logger = echoLog.Logger()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/foo", foo)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

func foo(c echo.Context) error {
	log.Infof("Handling Request from: %s", c.Path())

	linksEnv := os.Getenv("DEMO_LINKS")
	links := strings.Split(linksEnv, "|")

	for _, l := range links {
		log.Infof("Requesting: %s", l)
		resp, err := http.Get(l)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		if resp.StatusCode != http.StatusOK {
			return c.String(resp.StatusCode, err.Error())
		}
	}

	return c.String(http.StatusOK, "Success")
}
