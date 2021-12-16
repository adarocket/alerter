package web

import (
	"github.com/adarocket/alerter/internal/database/controller"
	database2 "github.com/adarocket/alerter/internal/database/model"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type AlertHandlers struct {
	alertController controller.Alert
}

func GetAlertHandlersInstance(cont controller.Alert) AlertHandlers {
	return AlertHandlers{alertController: cont}
}

func (a *AlertHandlers) getAlertByID(c *gin.Context) {
	tmpl, err := template.ParseFS(WebUI, "tmpl/getAlert.html")
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

	alerts, err := a.alertController.GetAlertByID(id)
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

func (a *AlertHandlers) updateAlert(c *gin.Context) {
	alertNode := database2.Alerts{}
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

	err = a.alertController.UpdateAlert(alertNode)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(c.Writer, c.Request, c.Request.URL.Host+homePage, http.StatusTemporaryRedirect)
}

func (a *AlertHandlers) getAlertsList(c *gin.Context) {
	tmpl, err := template.ParseFS(WebUI, "tmpl/getAlerts.html")
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	alerts, err := a.alertController.GetAlerts()
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

func (a *AlertHandlers) getEmptyAlertTmpl(c *gin.Context) {
	tmpl, err := template.ParseFS(WebUI, "tmpl/getEmptyAlert.html")
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

func (a *AlertHandlers) createAlert(c *gin.Context) {
	alertNode := database2.Alerts{}
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

	err = a.alertController.CreateAlert(alertNode)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(c.Writer, c.Request, c.Request.URL.Host+homePage, http.StatusFound)
}

func (a *AlertHandlers) deleteAlert(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	err = a.alertController.DeleteAlert(id)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(c.Writer, c.Request, c.Request.URL.Host+homePage, http.StatusTemporaryRedirect)
}
