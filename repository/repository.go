package repository

import "web3/pckge/models"

type DatabaseRepo interface {
	InsertPost(newPost models.Post) error
	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	AuthenticateUser(email, testPasword string) (int, string, error)
	GetAnArticle() (int, int, string, string, error)
	GetThreeArticles() (models.ArticleList, error)
}
