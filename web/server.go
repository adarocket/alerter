package web

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/adarocket/alerter/database"
	"github.com/adarocket/alerter/database/structs"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

const tokenName = "X-sessionToken"

var WebUI embed.FS

func authHandler(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		file, err := WebUI.ReadFile("data/auth.html")
		if err != nil {
			log.Println(err)
			http.Error(c.Writer, err.Error(), 500)
			return
		}

		tmpl, err := template.New("example").Parse(string(file))
		if err != nil {
			log.Println(err)
			http.Error(c.Writer, err.Error(), 500)
			return
		}

		err = tmpl.Execute(c.Writer, string(file))
		if err != nil {
			log.Println(err)
			http.Error(c.Writer, err.Error(), 500)
			return
		}
	case "POST":
		if c.Request.FormValue("name") != "ada" || c.Request.FormValue("password") != "rocket" {
			http.Error(c.Writer, "wrong username or password", 401)
			return
		}

		tokenStr, err := GenerateToken()
		if err != nil {
			http.Error(c.Writer, "ooops", http.StatusInternalServerError)
			return
		}

		cs := &http.Cookie{Name: tokenName, Value: tokenStr}
		http.SetCookie(c.Writer, cs)

		http.Redirect(c.Writer, c.Request, c.Request.Referer(), 302)
	}
}

func alertHandlerGet(c *gin.Context) {
	file, err := WebUI.ReadFile("data/getAlert.html")
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), 500)
		return
	}

	tmpl, err := template.New("example").Parse(string(file))
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), 500)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), 500)
		return
	}

	alerts, err := database.Sqllite.GetDataFromAlert(id)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), 500)
		return
	}

	err = tmpl.Execute(c.Writer, alerts)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), 500)
		return
	}
}

func alertHandlerPost(c *gin.Context) {
	alertNode := structs.AlertsTable{}
	var err error

	alertNode.Name = c.Request.FormValue("Name")
	alertNode.CheckedField = c.Request.FormValue("CheckedField")
	alertNode.TypeChecker = c.Request.FormValue("TypeChecker")

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), 500)
		return
	}
	alertNode.ID = id

	err = database.Sqllite.SetDataInAlertsTable(alertNode)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), 500)
		return
	}

	http.Redirect(c.Writer, c.Request, c.Request.URL.Host+"/alerts", 302)
}

func alertsHandlerGet(c *gin.Context) {
	file, err := WebUI.ReadFile("data/getAlerts.html")
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), 500)
		return
	}

	tmpl, err := template.New("example").Parse(string(file))
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), 500)
		return
	}

	alerts, err := database.Sqllite.GetDataFromAlerts()
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), 500)
		return
	}

	err = tmpl.Execute(c.Writer, alerts)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), 500)
		return
	}
}

func alertNodeHandlerGet(c *gin.Context) {
	file, err := WebUI.ReadFile("data/getAlertNode.html")
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), 500)
		return
	}

	tmpl, err := template.New("example").Parse(string(file))
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), 500)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), 500)
		return
	}

	alertNode, err := database.Sqllite.GetDataFromAlertNode(id)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), 500)
		return
	}

	err = tmpl.Execute(c.Writer, alertNode)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), 500)
		return
	}
}

func alertNodeHandlerPost(c *gin.Context) {
	alertNode := structs.AlertNodeTable{}
	var err error
	if alertNode.NormalFrom, err = strconv.ParseFloat(c.Request.FormValue("NormalFrom"), 64); err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), 500)
		return
	}
	if alertNode.NormalTo, err = strconv.ParseFloat(c.Request.FormValue("NormalTo"), 64); err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), 500)
		return
	}
	if alertNode.CriticalFrom, err = strconv.ParseFloat(c.Request.FormValue("CriticalFrom"), 64); err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), 500)
		return
	}
	if alertNode.CriticalTo, err = strconv.ParseFloat(c.Request.FormValue("CriticalTo"), 64); err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), 500)
		return
	}
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), 500)
		return
	}

	alertNode.Frequency = c.Request.FormValue("Frequency")
	alertNode.AlertID = id

	err = database.Sqllite.SetDataInAlertNode(alertNode)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, err.Error(), 500)
		return
	}

	http.Redirect(c.Writer, c.Request, c.Request.URL.Host+"/alerts", 302)
}

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
	router.GET("/alert/:id", alertHandlerGet)
	router.POST("/alert/:id", alertHandlerPost)
	router.GET("/alertNode/:id", alertNodeHandlerGet)
	router.POST("/alertNode/:id", alertNodeHandlerPost)
	router.GET("/alerts", alertsHandlerGet)

	http.Handle("/", router)

	fmt.Println("Server is listening...  http://127.0.0.1:5400/alerts")
	log.Fatal(http.ListenAndServe(":5400", nil))
}
