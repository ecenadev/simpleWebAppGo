package dbrepo

import (
	"context"
	"errors"
	"fmt"
	"time"
	"web3/pckge/models"

	"golang.org/x/crypto/bcrypt"
)

//functions or accessing database

func (m *postgresDBRepo) InsertPost(newPost models.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	query := `INSERT into posts(title, content, user_id) VALUES($1, $2,$3)`
	fmt.Println("------------- funct")
	_, err := m.DB.ExecContext(ctx, query, newPost.Title,
		newPost.Content, newPost.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `Select name, email, passowrd acct_created, last_login, user_type, id FROM users WHERE id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User

	err := row.Scan(
		&u.Name,
		&u.Email,
		&u.Password,
		&u.AcctCreated,
		&u.LastLogin,
		&u.UserType,
		&u.ID,
	)

	if err != nil {
		return u, err
	}

	return u, nil
}

func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `Update users set name=$1,email=$2, last_login =$3, user_type =$4 `

	_, err := m.DB.ExecContext(ctx, query, u.Name, u.Email, time.Now(), u.UserType)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) AuthenticateUser(email, testPasword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id int
	var hashedPW string

	query := `SELECT id, password from users where email = $1`

	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(&id, &hashedPW)

	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPW), []byte(testPasword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("password is incorrect")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPW, nil
}

func (m *postgresDBRepo) GetAnArticle() (int, int, string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id, uID int
	var aTitle, aContent string

	query := `SELECT id,user_id, title, content from posts Limit 1`

	row := m.DB.QueryRowContext(ctx, query)
	err := row.Scan(&id, &uID, &aTitle, &aContent)
	if err != nil {
		return id, uID, "", "", err
	}

	return id, uID, aTitle, aContent, nil

}

func (m *postgresDBRepo) GetThreeArticles() (models.ArticleList, error) {
	var artList models.ArticleList

	rows, err := m.DB.Query(`Select id, user_id, title, content from posts order by id DESC LIMIT $1`, 3)

	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, uID int
		var title, content string
		err = rows.Scan(&id, &uID, &title, &content)
		if err != nil {
			panic(err)
		}
		artList.ID = append(artList.ID, id)
		artList.UserID = append(artList.UserID, uID)
		artList.Title = append(artList.Title, title)
		artList.Content = append(artList.Content, content)

	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return artList, nil
}
