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
	alertNodeGroup := router.Group("/alertNode")
	{
		alertNodeGroup.GET("/:id/:uuid/edit", GetAlertNodeByIDAndUuid)
		alertNodeGroup.POST("/edit", updateAlertNode)
		alertNodeGroup.GET("/create", getEmptyAlertNodeTmpl)
		alertNodeGroup.POST("/create", createAlertNode)
		alertNodeGroup.GET("/:id/:uuid/delete", deleteAlertNode)
	}

	router.GET("/alertNodes/:id", getAlertNodesListByID)
	router.GET(homePage, getAlertsList)
	http.Handle("/", router)

	fmt.Println("Server is listening...  http://127.0.0.1:8080/alerts")
	log.Fatal(http.ListenAndServe(webServerAddr, nil))
}
