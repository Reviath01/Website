package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	Reset  = "\033[0m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[42m"
	server *gin.Engine
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	server = gin.New()
	server.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf(Cyan+"Web:    "+Reset+"| %v |"+Yellow+" %3d "+Reset+"| %13v | %15s |"+Green+" %-7s  "+Reset+"\"%s\" \n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
		)
	}))
	server.LoadHTMLGlob("public/*.html")
	server.Static("/assets", "./public/assets")
	server.GET("/", func(c *gin.Context) { c.HTML(200, "index.html", gin.H{}) })
	server.POST("/contact", func(c *gin.Context) {
		r := c.Request
		name := r.FormValue("name")
		message := r.FormValue("message")

		type jsonStruct struct {
			Name    string
			Message string
		}
		jsondata, _ := json.Marshal(jsonStruct{
			Name:    name,
			Message: message,
		})

		postBody, _ := json.Marshal(map[string]string{
			"username":   "Reviath",
			"content":    string(jsondata),
			"avatar_url": "https://cdn.discordapp.com/avatars/894273903600484384/a_83f7b181631e78c6f5beaf651f74e65c.gif?size=160",
		})
		resp, err := http.Post("https://discord.com/api/webhooks/948991633519702046/5eVX64T5hJQuU_pM2bZ2oXBR4PrpudFk6tz1m52L_G7WdnshILqsWUXkes0sfiW5KkCq", "application/json", bytes.NewBuffer(postBody))
		if err != nil {
			fmt.Printf("An Error Occured %v", err)
		}
		defer resp.Body.Close()
		c.Redirect(http.StatusMovedPermanently, "/")
	})
	var port string
	if os.Getenv("PORT") == "" {
		port = "8080"
	} else {
		port = os.Getenv("PORT")
	}
	if err := server.Run(":" + port); err != nil {
		print(err.Error() + "\n")
		os.Exit(0)
	}
}
