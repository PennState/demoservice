package main

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	echoLog "github.com/neko-neko/echo-logrus/v2/log"

	log "github.com/sirupsen/logrus"
)

var (
	initalized bool
	failures   int
	lastSeen   time.Time
)

func main() {

	initalized = false
	lastSeen = time.Now()
	failures = 0

	e := echo.New()

	// Middleware

	e.Logger = echoLog.Logger()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/foo", foo)
	e.GET("/bar", bar)
	e.GET("/baz", baz)
	e.GET("/delay", delay)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

func foo(c echo.Context) error {
	log.Infof("Handling Request from: %s", c.Path())

	linksEnv := os.Getenv("FOO_LINKS")
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

var mux sync.Mutex

func bar(c echo.Context) error {
	mux.Lock()
	defer mux.Unlock()

	timeout := os.Getenv("BAR_TIMEOUT")
	t, err := strconv.Atoi(timeout)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	time.Sleep(time.Duration(t) * time.Millisecond)

	return c.String(http.StatusOK, "Success")
}

func baz(c echo.Context) error {
	timeout := os.Getenv("BAZ_TIMEOUT")
	t, err := strconv.Atoi(timeout)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	time.Sleep(time.Duration(t) * time.Millisecond)

	return c.String(http.StatusOK, "Success")
}

func delay(c echo.Context) error {
	err := sleep(c)

	if err != nil {
		return err
	}

	return c.String(http.StatusOK, "Success")
}

func sleep(c echo.Context) error {
	mux.Lock()
	defer mux.Unlock()

	now := time.Now()

	if !initalized {
		lastSeen = now
		initalized = true
	}

	log.Printf("%v\n%v\n", lastSeen, now)

	d := now.Sub(lastSeen).Milliseconds()
	log.Printf("d = %d\n", d)

	if d < 1000 {
		r := rand.Intn(3)
		dur := time.Duration(r)
		log.Println("500")
		log.Printf("factor: %d\n", r)
		st := 500 * dur * time.Millisecond
		log.Printf("Sleeping %v\n", st.Seconds())
		time.Sleep(st)
	} else if d < 1500 && d > 1000 {
		r := rand.Intn(2) + 2
		dur := time.Duration(r)
		log.Printf("factor: %d\n", r)
		st := 500 * dur * time.Millisecond
		log.Printf("Sleeping %v\n", st.Seconds())
		time.Sleep(st)
	} else if d < 2000 && d > 1500 {
		r := rand.Intn(3) + 3
		log.Printf("factor: %d\n", r)
		dur := time.Duration(r)
		st := 500 * dur * time.Millisecond
		log.Printf("Sleeping %v\n", st.Seconds())
		time.Sleep(st)
	} else if d < 3000 && d > 2000 {
		r := rand.Intn(3) + 5
		log.Printf("factor: %d\n", r)
		dur := time.Duration(r)
		st := 500 * dur * time.Millisecond
		log.Printf("Sleeping %v\n", st.Seconds())
		time.Sleep(st)
	} else {
		failures++

		log.Printf("Failures: %d\n", failures)
		//Reset if we see 5 failures
		if failures > 5 {
			failures = 0
			lastSeen = now
		}

		return c.String(http.StatusRequestTimeout, "Timout")
	}

	lastSeen = now
	return nil
}
