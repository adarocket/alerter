package web

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
)

func getAuthHandler(c *gin.Context) {
	if c.Request.Method != "GET" {
		postAuthHandler(c)
		return
	}

	tmpl, err := template.ParseFS(WebUI, "data/auth.html")
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

func postAuthHandler(c *gin.Context) {
	if c.Request.Method != "POST" {
		return
	}

	if c.Request.FormValue("name") != "ada" || c.Request.FormValue("password") != "rocket" {
		http.Error(c.Writer, "wrong username or password", http.StatusUnauthorized)
		return
	}

	tokenStr, err := GenerateToken()
	if err != nil {
		http.Error(c.Writer, "ooops", http.StatusInternalServerError)
		return
	}

	cs := &http.Cookie{Name: tokenName, Value: tokenStr}
	http.SetCookie(c.Writer, cs)

	http.Redirect(c.Writer, c.Request, c.Request.Referer(), http.StatusTemporaryRedirect)
}
