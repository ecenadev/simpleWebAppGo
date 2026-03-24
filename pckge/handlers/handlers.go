package handlers

import (
	"log"
	"net/http"
	"web3/pckge/config"
	"web3/pckge/form"
	"web3/pckge/models"
	"web3/pckge/render"
)

type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

func NewRepo(ac *config.AppConfig) *Repository {
	return &Repository{
		App: ac,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) HomeHandler(w http.ResponseWriter, request *http.Request) {

	m.App.Session.Put(request.Context(), "userid", "01elio") // crates a sesion and add a value

	render.RenderTemplate(w, "home.page.tmpl", &models.PageData{}, request)
}

func (m *Repository) AboutHandler(w http.ResponseWriter, request *http.Request) {

	strMap := make(map[string]string)
	render.RenderTemplate(w, "about.page.tmpl", &models.PageData{StrMap: strMap}, request)
	// strMap["title"] = "About Us"
	// strMap["intro"] = "This page is where we talk about ourselves"

	// usrid := m.App.Session.GetString(request.Context(), "userid")
	// strMap["userid"] = usrid //puts sessions in the dataMap

}

func (m *Repository) LoginHandler(w http.ResponseWriter, request *http.Request) {
	strMap := make(map[string]string)
	render.RenderTemplate(w, "login.page.tmpl", &models.PageData{StrMap: strMap}, request)

}

func (m *Repository) MakePostHandler(w http.ResponseWriter, request *http.Request) {

	var emptyArticle models.Article

	data := make(map[string]interface{})

	data["article"] = emptyArticle

	render.RenderTemplate(w, "make-post.page.tmpl", &models.PageData{
		Form: form.New(nil),
		Data: data,
	}, request)
}
func (m *Repository) PostMakePostHandler(w http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()

	if err != nil {
		log.Println(err)
		return
	}

	// blog_title := reqest.Form.Get("blog_title")
	// blog_article := reqest.Form.Get("blog_article")

	// w.Write([]byte(blog_title))
	// w.Write([]byte(blog_article))

	article := models.Article{
		BlogTitle:   request.Form.Get("blog_title"),
		BlogArticle: request.Form.Get("blog_article"),
	}

	form := form.New(request.PostForm)
	// form.HasValue("blog_title", request)

	form.HasRequired("blog_title", "blog_aticle")

	form.MinLength("blog_title", 5, request)
	form.MinLength("blog_article", 5, request)

	// form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["article"] = article

		render.RenderTemplate(w, "article-recieved.page.tmpl", &models.PageData{
			Form: form,
			Data: data,
		}, request)
		return
	}
	m.App.Session.Put(request.Context(), "article", article)
	http.Redirect(w, request, "/article-received", http.StatusSeeOther)
}

func (m *Repository) ArticleReceived(w http.ResponseWriter, r *http.Request) {
	article, ok := m.App.Session.Get(r.Context(), "article").(models.Article)
	if !ok {
		log.Println("Cant get data from sessioon")

		m.App.Session.Put(r.Context(), "error", "Cant get data from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

		return
	}
	data := make(map[string]interface{})
	data["article"] = article
	render.RenderTemplate(w, "article-recieved.page.tmpl",
		&models.PageData{
			Data: data,
		}, r)
}

func (m *Repository) PageHandler(w http.ResponseWriter, request *http.Request) {
	strMap := make(map[string]string)
	render.RenderTemplate(w, "page.page.tmpl", &models.PageData{StrMap: strMap}, request)
}
