package web

import (
	database2 "github.com/adarocket/alerter/internal/database"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func getAlertNodeByID(c *gin.Context) {
	tmpl, err := template.ParseFS(WebUI, "data/getAlertNode.html")
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

	alertNode, err := database2.Db.GetNodeAlertByID(id)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(c.Writer, alertNode)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}
}

func updateAlertNode(c *gin.Context) {
	alertNode := database2.AlertNode{}
	var err error
	if alertNode.NormalFrom, err = strconv.ParseFloat(c.Request.FormValue("NormalFrom"), 64); err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}
	if alertNode.NormalTo, err = strconv.ParseFloat(c.Request.FormValue("NormalTo"), 64); err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}
	if alertNode.CriticalFrom, err = strconv.ParseFloat(c.Request.FormValue("CriticalFrom"), 64); err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}
	if alertNode.CriticalTo, err = strconv.ParseFloat(c.Request.FormValue("CriticalTo"), 64); err != nil {
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

	alertNode.Frequency = c.Request.FormValue("Frequency")
	alertNode.AlertID = id

	err = database2.Db.UpdateAlertNode(alertNode)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(c.Writer, c.Request, c.Request.URL.Host+homePage, http.StatusFound)
}

func actionChose(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFS(WebUI, "data/actionMenuAlertNode.html")
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(c.Writer, database2.AlertNode{AlertID: id})
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}
}

func createAlertNode(c *gin.Context) {
	tmpl, err := template.ParseFS(WebUI, "data/getAlertNode.html")
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

	err = tmpl.Execute(c.Writer, database2.AlertNode{AlertID: id})
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}
}

func deleteAlertNode(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	err = database2.Db.DeleteAlertNode(id)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(c.Writer, c.Request, c.Request.URL.Host+homePage, http.StatusFound)
}