package main

import "github.com/gin-gonic/gin"
import "net/http"

func main() {
    r := gin.Default()
    r.LoadHTMLGlob("tmpl/*")
    r.Static("/assets", "assets")
    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.tmpl", gin.H {
            "title": "傳說中的考古題系統",
        })
    })
    r.GET("/posts/:id", func(c *gin.Context) {
        // id := c.Param("id")
        c.String(http.StatusOK, "posts.tmpl", gin.H {
            "title": "傳說中的考古題系統",
            "content": "",
        })
    })
    r.Run(":8000")
}
