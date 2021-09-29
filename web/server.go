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
const homePage = "/alerts"

var WebUI embed.FS

func authMw(c *gin.Context) {
	cookies := c.Request.Cookies()

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
	router.GET("/alert/:id/delete", deleteAlert)
	router.GET("/alert/create", getEmptyAlertTmpl)
	router.POST("/alert/create", createAlert)
	router.POST("/alert/:id/edit", updateAlert)
	router.GET("/alertNode/:id/edit", getAlertNodeByID)
	router.POST("/alertNode/:id/edit", createAlertNode)
	router.GET(homePage, getAlertsList)

	http.Handle("/", router)

	fmt.Println("Server is listening...  http://127.0.0.1:5400/alerts")
	log.Fatal(http.ListenAndServe(webServerAddr, nil))
}
