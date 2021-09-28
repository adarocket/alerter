package web

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// FIXME: разбей весь файл, вынеси:
// обработчики по созданию алертов отдельно
// обработчики по привязке алерта к ноде отдельно

const tokenName = "X-sessionToken"

// FIXME: почему не в конфиге?
// FIXME: почему глобальная переменная?
var WebServerAddr = ":5400"

var WebUI embed.FS

// FIXME: названия на 10 из 10.
// прослойки можно привязывать к группам
func simpleMw(c *gin.Context) {
	cookies := c.Request.Cookies()
	if len(cookies) < 1 {
		authHandler(c)
		return
	} else {
		if isValid := IsValidToken(cookies[0].Value); !isValid {
			authHandler(c)
			return
		}
	}

	c.Next()
}

func StartServer() {
	router := gin.Default()
	router.Use(simpleMw)

	// FIXME: где полноценный CRUD
	// FIXME добавь группы
	router.GET("/alert/:id", getAlertByID)
	router.POST("/alert/:id", createAlert) // FIXME: это создание алерта?
	router.GET("/alertNode/:id", getAlertNodeByID)
	router.POST("/alertNode/:id", createAlertNode) // FIXME: это создание чего?
	router.GET("/alerts", getAlertsList)

	http.Handle("/", router)

	fmt.Println("Server is listening...  http://127.0.0.1:5400/alerts")
	log.Fatal(http.ListenAndServe(WebServerAddr, nil))
}
