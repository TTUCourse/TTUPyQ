package main

import "database/sql"
import "github.com/gin-gonic/gin"
import _ "github.com/go-sql-driver/mysql"
import "gopkg.in/gcfg.v1"
import "gopkg.in/gorp.v1"
import "log"
import "net/http"
import "strconv"
import "fmt"

type PYQ struct {
    Id      int64
    Author  string
    Title   string
    Content string
    Time    string
}

type PostInfo struct {
    Id      int64
    Author  string
    Title   string
    Time    string
}

type Config struct {
    MysqlUsername   string
    MysqlPassword   string
    MysqlLocation   string
    MysqlDbname     string
}

type configFile struct {
    Server Config
}

func main() {
    route := gin.Default()
    route.LoadHTMLGlob("tmpl/*")
    route.Static("/assets", "assets")
    route.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.tmpl", gin.H {
            "title": "傳說中的考古題系統",
        })
    })
    route.GET("/posts/:id", getPostsId)
    route.POST("/api/posts/", postApiPostsPage)
    route.Run(":8000")
}

var config = LoadConfiguration()

func LoadConfiguration() Config {
    var err error
    var cfg configFile

    err = gcfg.ReadFileInto(&cfg, "./server.conf")

    if err != nil {
        panic(err.Error())
    }

    return cfg.Server
}

var dbmap = initDb()

func initDb() *gorp.DbMap {
    db, err := sql.Open("mysql", config.MysqlUsername + ":" + config.MysqlPassword +
                        "@" + config.MysqlLocation + "/" + config.MysqlDbname)
    checkErr(err ,"sql.Open failed")
    dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
    dbmap.AddTableWithName(PYQ{}, "PYQ").SetKeys(true, "Id")
    err = dbmap.CreateTablesIfNotExists()
    checkErr(err, "Create table failed")
    return dbmap
}

func checkErr(err error, msg string) {
    if err != nil {
        log.Fatalln(msg, err)
    }
}

func getPostsId(c *gin.Context) {
    id := c.Param("id")
    var content string
    err := dbmap.SelectOne(&content, "SELECT content FROM PYQ WHERE id = ?", id)
    if err == nil {
        c.HTML(http.StatusOK, "posts.tmpl", gin.H {
            "title": "傳說中的考古題系統",
            "content": content,
        })
    }
    else {
        c.String(http.StatusNotFound, "安安沒找到喔")
    }
}

func postApiPostsPage(c *gin.Context) {
    page, _ := strconv.ParseInt(c.PostForm("p"), 0, 64)
    fmt.Printf("%d\n",page)
    var list []PostInfo
    _, err := dbmap.Select(&list, "SELECT id, author, title, time FROM PYQ LIMIT ?,10", (page - 1) * 10)
    if err == nil {
        c.JSON(http.StatusOK, list)
    }
    else {
        checkErr(err, "postApi Error")
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
    }
}
