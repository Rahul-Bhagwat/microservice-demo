package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/ianschenck/envflag"
	//"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
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

type Message struct {
	Id   int
	User string
	Data string
}

var db *sql.DB

func main() {
	bindAddress := envflag.String(
		"ADDR",
		":8080",
		"Bind address for the api server")

	lolCatAddress = envflag.String(
		"LOLCAT_ADDR",
		"127.0.0.1:8081",
		"Bind address for the api server")

	mysqlConnectionString := envflag.String(
		"MYSQL_DSN",
		"root:@tcp(127.0.0.1:3306)/microservice",
		"Mysql connection information")

	//TODO: add an auth layer
	r := gin.Default()
	r.Use(Logger())
	r.Use(static.Serve("/", static.LocalFile("./static", true)))

	var err error
	db, err = sql.Open("mysql", *mysqlConnectionString)
	if err != nil {
		log.Fatal("can't connect to database")
	}
	defer db.Close()
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}
	migrate(db)

	r.GET("/messages", allMessages)
	r.POST("/message", createMessage)
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error parsing value from http request: %s", err)
		//return "", err
	}

	c.JSON(http.StatusOK, gin.H{
		"lolcat_url": body,
	})
}

func allMessages(c *gin.Context) {
	var (
		message  Message
		messages []Message
	)
	rows, err := db.Query("select id, user, data from messages;")
	if err != nil {
		fmt.Print(err.Error())
	}
	for rows.Next() {
		err = rows.Scan(&message.Id, &message.User, &message.Data)
		messages = append(messages, message)
		if err != nil {
			fmt.Print(err.Error())
		}
	}
	defer rows.Close()
	c.JSON(http.StatusOK, gin.H{
		"result": messages,
		"count":  len(messages),
	})
}

func createMessage(c *gin.Context) {
	var buffer bytes.Buffer
	//	user := c.PostForm("user")
	user := "unknown" //TODO Add auth
	data := c.PostForm("data")
	stmt, err := db.Prepare("insert into messages (user, data) values(?,?);")
	if err != nil {
		fmt.Print(err.Error())
	}
	_, err = stmt.Exec(user, data)

	if err != nil {
		fmt.Print(err.Error())
	}

	// Fastest way to append strings
	buffer.WriteString(user)
	buffer.WriteString(" ")
	buffer.WriteString(data)
	defer stmt.Close()
	name := buffer.String()
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf(" %s successfully created", name),
	})
}

func migrate(db *sql.DB) {
	stmt, err := db.Prepare("CREATE TABLE messages (id int NOT NULL AUTO_INCREMENT, user varchar(40), data varchar(40), PRIMARY KEY (id));")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Message Table successfully migrated....")
	}
}
