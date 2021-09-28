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

func getAlertByID(c *gin.Context) {
	tmpl, err := template.ParseFS(WebUI, "data/getAlert.html")
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
	tmpl, err := template.ParseFS(WebUI, "data/getAlerts.html")
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
