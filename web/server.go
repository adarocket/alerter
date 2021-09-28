package web

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const tokenName = "X-sessionToken"

var WebUI embed.FS

func authMw(c *gin.Context) {
	cookies := c.Request.Cookies()

	if c.Request.Method == "POST" {
		postAuthHandler(c)
		return
	}

	if len(cookies) < 1 {
		getAuthHandler(c)
		c.Abort()
	} else {
		if isValid := IsValidToken(cookies[0].Value); !isValid {
			getAuthHandler(c)
			c.Abort()
		}
	}
}

func StartServer(webServerAddr string) {
	router := gin.Default()
	router.Use(authMw)

	// FIXME: где полноценный CRUD
	// FIXME добавь группы
	router.GET("/alert/:id/edit", getAlertByID)
	router.POST("/alert/:id/edit", createAlert) // FIXME: это создание алерта?
	router.GET("/alertNode/:id/edit", getAlertNodeByID)
	router.POST("/alertNode/:id/edit", createAlertNode) // FIXME: это создание чего?
	router.GET("/alerts", getAlertsList)

	http.Handle("/", router)

	fmt.Println("Server is listening...  http://127.0.0.1:5400/alerts")
	log.Fatal(http.ListenAndServe(webServerAddr, nil))
}
