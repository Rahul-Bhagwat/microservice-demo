package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/ianschenck/envflag"
	//"github.com/Sirupsen/logrus"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Set example variable
		c.Set("example", "12345")

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}

var lolCatAddress *string

func main() {
	bindAddress := envflag.String(
		"ADDR",
		":8080",
		"Bind address for the api server")

	lolCatAddress = envflag.String(
		"LOLCAT_ADDR",
		"127.0.0.1:8081",
		"Bind address for the api server")

	//TODO: add an auth layer
	r := gin.Default()
	r.Use(Logger())
	r.Use(static.Serve("/", static.LocalFile("./static", true)))

	//	p := ginprometheus.NewPrometheus("gin")
	//	p.Use(r)

	r.GET("/get_lolcat", getLolCats)

	// Listen and server on 0.0.0.0:8080
	r.Run(*bindAddress)
}

func getLolCats(c *gin.Context) {
	getURL := fmt.Sprintf("http://%s/random_lolcat", *lolCatAddress) //s.dropletID , can we just give a bogus value?

	//log.KV("URL", getURL).Debug("Downloading lolcats")
	req, err := http.NewRequest("GET", getURL, nil)
	if err != nil {
		panic(err)
		return //TODO ERROR
		//		return 0, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
		return //TODO ERROR
		//		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("unexpected status code %d while pushing to %s", resp.StatusCode, getURL)
		panic(err)
		//log.KV("code", resp.StatusCode).KV("url", getURL).Error("failed pushing to wharf")
		return //TODO ERROR
		//		return 0, err
	}
	c.JSON(http.StatusOK, gin.H{
		"lolcat_url": resp.Body,
	})
}
