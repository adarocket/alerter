package web

import (
	"github.com/adarocket/alerter/database"
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
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	alerts, err := database.Db.GetAlertByID(id)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(c.Writer, alerts)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}
}

func updateAlert(c *gin.Context) {
	alertNode := database.Alerts{}
	var err error

	alertNode.Name = c.Request.FormValue("Name")
	alertNode.CheckedField = c.Request.FormValue("CheckedField")
	alertNode.TypeChecker = c.Request.FormValue("TypeChecker")

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}
	alertNode.ID = id

	err = database.Db.UpdateAlert(alertNode)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(c.Writer, c.Request, c.Request.URL.Host+homePage, http.StatusTemporaryRedirect)
}

func getAlertsList(c *gin.Context) {
	tmpl, err := template.ParseFS(WebUI, "data/getAlerts.html")
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	alerts, err := database.Db.GetAlerts()
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(c.Writer, alerts)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}
}

func getEmptyAlertTmpl(c *gin.Context) {
	tmpl, err := template.ParseFS(WebUI, "data/getEmptyAlert.html")
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(c.Writer, nil)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}
}

func createAlert(c *gin.Context) {
	alertNode := database.Alerts{}
	var err error

	alertNode.Name = c.Request.FormValue("Name")
	alertNode.CheckedField = c.Request.FormValue("CheckedField")
	alertNode.TypeChecker = c.Request.FormValue("TypeChecker")

	idStr := c.Request.FormValue("ID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}
	alertNode.ID = id

	err = database.Db.CreateAlert(alertNode)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(c.Writer, c.Request, c.Request.URL.Host+homePage, http.StatusFound)
}

func deleteAlert(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	err = database.Db.DeleteAlert(id)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(c.Writer, c.Request, c.Request.URL.Host+homePage, http.StatusTemporaryRedirect)
}
