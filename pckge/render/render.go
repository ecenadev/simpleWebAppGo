package render

import (
	"fmt"
	"html/template"
	"net/http"
	"web3/pckge/config"
	"web3/pckge/models"

	"github.com/justinas/nosurf"
)

var tmplCache = make(map[string]*template.Template)

var app *config.AppConfig

func NewAppConfig(a *config.AppConfig) {
	app = a
}

func AddCSRFData(pd *models.PageData, r *http.Request) *models.PageData {
	pd.CSRFToken = nosurf.Token(r)

	if app.Session.Exists(r.Context(), "user_id") {
		pd.IsAuthenticated = 1
	} else {
		pd.IsAuthenticated = 0
	}

	return pd
}

func RenderTemplate(w http.ResponseWriter, t string, pd *models.PageData, request *http.Request) {
	var tmpl *template.Template
	var err error
	_, inMap := tmplCache[t]
	fmt.Println(tmplCache) //check wether the template is in cache
	if !inMap {
		err = makeTemplateCache(t)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Template in cache")
	}

	//load the template to the template variable
	tmpl = tmplCache[t]
	pd = AddCSRFData(pd, request)
	err = tmpl.Execute(w, pd)
	if err != nil {
		fmt.Println(err)
	}

}

func makeTemplateCache(t string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", t),
		"./templates/base.layout.tmpl",
	}

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}
	tmplCache[t] = tmpl
	return nil
}
