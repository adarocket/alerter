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

func getAlertNodeByID(c *gin.Context) {
	// start
	file, err := WebUI.ReadFile("data/getAlertNode.html")
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
	// end
	// FIXME: этот кусок постоянно повторяется, вынеси в отдельную функцию

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", 500)
		return
	}

	alertNode, err := database.Sqllite.GetDataFromAlertNode(id)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", 500)
		return
	}

	err = tmpl.Execute(c.Writer, alertNode)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", 500)
		return
	}
}

func createAlertNode(c *gin.Context) {
	alertNode := structs.AlertNodeTable{}
	var err error
	if alertNode.NormalFrom, err = strconv.ParseFloat(c.Request.FormValue("NormalFrom"), 64); err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", 500)
		return
	}
	if alertNode.NormalTo, err = strconv.ParseFloat(c.Request.FormValue("NormalTo"), 64); err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", 500)
		return
	}
	if alertNode.CriticalFrom, err = strconv.ParseFloat(c.Request.FormValue("CriticalFrom"), 64); err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", 500)
		return
	}
	if alertNode.CriticalTo, err = strconv.ParseFloat(c.Request.FormValue("CriticalTo"), 64); err != nil {
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

	alertNode.Frequency = c.Request.FormValue("Frequency")
	alertNode.AlertID = id

	err = database.Sqllite.UpdateDataInAlertNode(alertNode)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", 500)
		return
	}

	http.Redirect(c.Writer, c.Request, c.Request.URL.Host+"/alerts", 302)
}
