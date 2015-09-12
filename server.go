package main

import "github.com/gin-gonic/gin"
import "net/http"

func main() {
    r := gin.Default()
    r.LoadHTMLGlob("tmpl/*")
    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.tmpl", gin.H {
            "title": "Main website",
        })
    })
    r.GET("/user/:name", func(c *gin.Context) {
        name := c.Param("name")
        c.String(http.StatusOK, "Hello %s", name)
    })
    r.GET("/user/:name/*action", func(c *gin.Context) {
        name := c.Param("name")
        action := c.Param("action")
        message := name + " is " + action
        c.String(http.StatusOK, message)
    })
    r.Run(":80")
}