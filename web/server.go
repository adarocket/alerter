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

	alertGroup := router.Group("/alert")
	{
		alertGroup.GET("/:id/edit", getAlertByID)
		alertGroup.GET("/:id/delete", deleteAlert)
		alertGroup.GET("/create", getEmptyAlertTmpl)
		alertGroup.POST("/create", createAlert)
		alertGroup.POST("/:id/edit", updateAlert)
	}
	alertNodeGroup := router.Group("/alertNode/:id")
	{
		alertNodeGroup.GET("/edit", getAlertNodeByID)
		alertNodeGroup.POST("/edit", updateAlertNode)
		alertNodeGroup.GET("/action", actionChose)
		alertNodeGroup.GET("/create", createAlertNode)
		alertNodeGroup.GET("/delete", deleteAlertNode)
	}

	router.GET(homePage, getAlertsList)
	http.Handle("/", router)

	fmt.Println("Server is listening...  http://127.0.0.1:5400/alerts")
	log.Fatal(http.ListenAndServe(webServerAddr, nil))
}
