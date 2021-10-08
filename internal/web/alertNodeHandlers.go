package web

import (
	"github.com/adarocket/alerter/internal/controller"
	"github.com/adarocket/alerter/internal/database"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type AlertNodeHandlers struct {
	alertNodeController controller.AlertNode
}

func GetAlertNodeHandlersInstance(cont controller.AlertNode) AlertNodeHandlers {
	return AlertNodeHandlers{alertNodeController: cont}
}

func (a *AlertNodeHandlers) getAlertNodesListByID(c *gin.Context) {
	tmpl, err := template.ParseFS(WebUI, "data/getAlertNodes.html")
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

	alertNodes, err := a.alertNodeController.GetAlertNodesByID(id)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(c.Writer, alertNodes)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (a *AlertNodeHandlers) createAlertNode(c *gin.Context) {
	idStr := c.Request.FormValue("AlertID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	normalFromStr := c.Request.FormValue("NormalFrom")
	normalFrom, err := strconv.ParseFloat(normalFromStr, 10)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}
	normalToStr := c.Request.FormValue("NormalTo")
	normalTo, err := strconv.ParseFloat(normalToStr, 10)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}
	criticalFromStr := c.Request.FormValue("CriticalFrom")
	criticalFrom, err := strconv.ParseFloat(criticalFromStr, 10)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}
	criticalToStr := c.Request.FormValue("CriticalTo")
	criticalTo, err := strconv.ParseFloat(criticalToStr, 10)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}
	frequency := c.Request.FormValue("Frequency")
	nodeUuid := c.Request.FormValue("NodeUuid")

	alertNode := database.AlertNode{
		AlertID:      id,
		NormalFrom:   normalFrom,
		NormalTo:     normalTo,
		CriticalFrom: criticalFrom,
		CriticalTo:   criticalTo,
		Frequency:    frequency,
		NodeUuid:     nodeUuid,
	}

	err = a.alertNodeController.CreateAlertNode(alertNode)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(c.Writer, c.Request, c.Request.URL.Host+homePage, http.StatusFound)
}

func (a *AlertNodeHandlers) GetAlertNodeByIDAndUuid(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	alertNode, err := a.alertNodeController.GetAlertNodeByIdAndNodeUuid(id, c.Param("uuid"))
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFS(WebUI, "data/getAlertNode.html")
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

func (a *AlertNodeHandlers) updateAlertNode(c *gin.Context) {
	alertNode := database.AlertNode{}
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

	idStr := c.Request.FormValue("AlertID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	alertNode.Frequency = c.Request.FormValue("Frequency")
	alertNode.AlertID = id
	alertNode.NodeUuid = c.Request.FormValue("NodeUuid")

	err = a.alertNodeController.UpdateAlertNode(alertNode)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(c.Writer, c.Request, c.Request.URL.Host+homePage, http.StatusFound)
}

func (a *AlertNodeHandlers) getEmptyAlertNodeTmpl(c *gin.Context) {
	tmpl, err := template.ParseFS(WebUI, "data/getEmptyAlertNode.html")
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

func (a *AlertNodeHandlers) deleteAlertNode(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	nodeUuid := c.Param("uuid")

	err = a.alertNodeController.DeleteAlertNode(id, nodeUuid)
	if err != nil {
		log.Println(err)
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(c.Writer, c.Request, c.Request.URL.Host+homePage, http.StatusFound)
}
