package web

import (
	"github.com/adarocket/alerter/database"
	"github.com/adarocket/alerter/database/structs"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// FIXME: Не делай так. Сделай два разных обработчика для разных методов
func authHandler(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		// FIXME: почему так сложно?
		// FIXME: почему не используешь template.ParseFS() ?

		file, err := WebUI.ReadFile("data/auth.html")
		if err != nil {
			log.Println(err)
			http.Error(c.Writer, "internal server error", 500) // FIXME: не используй хардкод, используй глобальные переменные http.StatusInternalServerError
			return
		}

		// FIXME: что такое example? Не используй такие названия
		tmpl, err := template.New("example").Parse(string(file))
		if err != nil {
			log.Println(err)
			http.Error(c.Writer, "internal server error", 500)
			return
		}

		err = tmpl.Execute(c.Writer, string(file))
		if err != nil {
			log.Println(err)
			http.Error(c.Writer, "internal server error", 500)
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

// FIXME: не информативные имена функций. не понятно что делает то или иной обработчик
// Например getAlertByID
func getAlertByID(c *gin.Context) {
	file, err := WebUI.ReadFile("data/getAlert.html")
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", 500)
		return
	}

	tmpl, err := template.New("example").Parse(string(file))
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", 500)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", 500)
		return
	}

	alerts, err := database.Sqllite.GetDataFromAlert(id)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", 500)
		return
	}

	err = tmpl.Execute(c.Writer, alerts)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", 500)
		return
	}
}

// Например createAlert
func createAlert(c *gin.Context) {
	alertNode := structs.AlertsTable{}
	var err error

	alertNode.Name = c.Request.FormValue("Name")
	alertNode.CheckedField = c.Request.FormValue("CheckedField")
	alertNode.TypeChecker = c.Request.FormValue("TypeChecker")

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", 500)
		return
	}
	alertNode.ID = id

	err = database.Sqllite.UpdateDataInAlertsTable(alertNode)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", 500)
		return
	}

	http.Redirect(c.Writer, c.Request, c.Request.URL.Host+"/alerts", 302)
}

func getAlertsList(c *gin.Context) {
	file, err := WebUI.ReadFile("data/getAlerts.html")
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", 500)
		return
	}

	tmpl, err := template.New("example").Parse(string(file))
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", 500)
		return
	}

	alerts, err := database.Sqllite.GetDataFromAlerts()
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", 500)
		return
	}

	err = tmpl.Execute(c.Writer, alerts)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", 500)
		return
	}
}
