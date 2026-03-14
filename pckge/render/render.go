package render

import (
	"fmt"
	"html/template"
	"net/http"
	"web3/pckge/models"
)

var tmplCache = make(map[string]*template.Template)

func RenderTemplate(w http.ResponseWriter, t string, pd *models.PageData) {
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
