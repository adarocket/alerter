package web

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
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
	//router.GET("/params", paramsHandler)
	//router.GET("/param/:name", changeParamHandlerGet)
	//router.POST("/param/:name", changeParamHandlerPost)

	http.Handle("/", router)

	fmt.Println("Server is listening...  http://127.0.0.1:5400/")
	log.Fatal(http.ListenAndServe(":5400", nil))
}
