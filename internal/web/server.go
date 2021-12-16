package web

import (
	"embed"
	"fmt"
	"github.com/adarocket/alerter/internal/database/controller"
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
			c.Redirect(http.StatusFound, c.Request.URL.Host+"/auth")
			c.Abort()
		}
	}
}

func StartServer(webServerAddr string) {
	alertHandlers := GetAlertHandlersInstance(controller.GetAlertControllerInstance())
	alertNodeController := GetAlertNodeHandlersInstance(controller.GetAlertNodeControllerInstance())

	router := gin.Default()
	router.Use(authMw)
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	alertGroup := router.Group("/alert")
	{
		alertGroup.GET("/:id/edit", alertHandlers.getAlertByID)
		alertGroup.GET("/:id/delete", alertHandlers.deleteAlert)
		alertGroup.GET("/create", alertHandlers.getEmptyAlertTmpl)
		alertGroup.POST("/create", alertHandlers.createAlert)
		alertGroup.POST("/:id/edit", alertHandlers.updateAlert)
	}
	alertNodeGroup := router.Group("/alertNode")
	{
		alertNodeGroup.GET("/:id/:uuid/edit", alertNodeController.GetAlertNodeByIDAndUuid)
		alertNodeGroup.POST("/edit", alertNodeController.updateAlertNode)
		alertNodeGroup.GET("/create", alertNodeController.getEmptyAlertNodeTmpl)
		alertNodeGroup.POST("/create", alertNodeController.createAlertNode)
		alertNodeGroup.GET("/:id/:uuid/delete", alertNodeController.deleteAlertNode)
	}

	router.GET("/alertNodes/:id", alertNodeController.getAlertNodesListByID)
	router.GET("/auth", getAuthHandler)
	router.POST("/auth", postAuthHandler)
	router.GET(homePage, alertHandlers.getAlertsList)
	http.Handle("/", router)

	fmt.Println("Server is listening...  http://127.0.0.1:8080/alerts")
	log.Fatal(http.ListenAndServe(webServerAddr, nil))
}
