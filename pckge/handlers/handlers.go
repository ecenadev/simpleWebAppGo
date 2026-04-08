package handlers

import (
	"fmt"
	"log"
	"net/http"
	"web3/pckge/config"
	"web3/pckge/dbdriver"
	"web3/pckge/form"
	"web3/pckge/models"
	"web3/pckge/render"
	"web3/repository"
	"web3/repository/dbrepo"
)

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

type PageData *models.PageData

var Repo *Repository

func NewRepo(ac *config.AppConfig, db *dbdriver.DB) *Repository {
	return &Repository{
		App: ac,
		DB:  dbrepo.NewPostgresRepo(db.SQL, ac),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) HomeHandler(w http.ResponseWriter, request *http.Request) {

	// id, uid, title, content, err := m.DB.GetAnArticle()

	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// fmt.Println("ID : ", id)
	// fmt.Println("User ID : ", uid)
	// fmt.Println("Title : ", title)
	// fmt.Println("Content : ", content)

	var artList models.ArticleList
	artList, err := m.DB.GetThreeArticles()

	if err != nil {
		log.Println(err)
		return
	}

	for i := range artList.Content {
		fmt.Println(artList.Content[i])
	}

	data := make(map[string]interface{})
	data["articleList"] = artList

	// m.App.Session.Put(request.Context(), "user_id", "01elio") // crates a sesion and add a value

	render.RenderTemplate(w, "home.page.tmpl", &models.PageData{Data: data}, request)
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
	if !m.App.Session.Exists(request.Context(), "user_id") {
		http.Redirect(w, request, "/login", http.StatusTemporaryRedirect)
	}
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

	uID := (m.App.Session.Get(request.Context(), "user_id")).(int)

	// blog_title := reqest.Form.Get("blog_title")
	// blog_article := reqest.Form.Get("blog_article")

	// w.Write([]byte(blog_title))
	// w.Write([]byte(blog_article))

	article := models.Post{
		Title:   request.Form.Get("blog_title"),
		Content: request.Form.Get("blog_article"),
		UserID:  int(uID),
	}

	form := form.New(request.PostForm)
	// form.HasValue("blog_title", request)

	form.HasRequired("blog_title", "blog_aticle")

	form.MinLength("blog_title", 5, request)
	form.MinLength("blog_article", 5, request)

	// form.IsEmail("email")

	// if !form.Valid() {
	// 	data := make(map[string]interface{})
	// 	data["article"] = article

	// 	render.RenderTemplate(w, "article-recieved.page.tmpl", &models.PageData{
	// 		Form: form,
	// 		Data: data,
	// 	}, request)
	// 	return
	// }

	//Write to DB
	err = m.DB.InsertPost(article)
	if err != nil {
		log.Fatal(err)
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

func (m *Repository) PostLoginHandler(w http.ResponseWriter, request *http.Request) {

	_ = m.App.Session.RenewToken(request.Context())

	err := request.ParseForm()
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		fmt.Println(err)
	}

	email := request.Form.Get("email")
	password := request.Form.Get("password")

	form := form.New(request.Form)
	form.HasRequired("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		render.RenderTemplate(w, "login.page.tmpl", &models.PageData{
			Form: form,
		}, request)

		return
	}

	id, _, err := m.DB.AuthenticateUser(email, password)
	if err != nil {
		m.App.Session.Put(request.Context(), "error", "invalid email or password")
		http.Redirect(w, request, "/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(request.Context(), "user_id", id)
	m.App.Session.Put(request.Context(), "flash", "Valid Login")
	http.Redirect(w, request, "/", http.StatusSeeOther)

}

func (m *Repository) LogOutHandler(w http.ResponseWriter, r *http.Request) {
	m.App.Session.Remove(r.Context(), "user_id")
	m.App.Session.Destroy(r.Context())
	m.App.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}
