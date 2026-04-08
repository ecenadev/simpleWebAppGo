package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func updateUserEmail(conn *sql.DB, newEmail string, id int) error {
	query := fmt.Sprintf(`Update users set email='%s' where id = '%d'`, newEmail, id)
	_, err := conn.Exec(query)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func getUserData(conn *sql.DB, id int) error {
	var name, email, pw, uType string
	// var id int

	query := fmt.Sprintf(`Select id, name,email, password, user_type from users where id = %d`, id)

	row := conn.QueryRow(query)
	err := row.Scan(&id, &name, &email, &pw, &uType)

	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("ID : ", id)
	fmt.Println("Name : ", name)
	fmt.Println("Email : ", email)
	return nil

}

func insertNewUser(conn *sql.DB, name string, email string, pw string, uType int) error {
	query := fmt.Sprintf(`INSERT INTO USERS (name, email, password, acct_created, last_login, user_type) VALUES 
	('%s','%s','%s',current_timestamp, current_timestamp,'%v')`, name, email, pw, uType)
	_, err := conn.Exec(query)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func getAllRowData(conn *sql.DB) error {
	rows, err := conn.Query("SELECT id,name,email from users")
	if err != nil {
		log.Println(err)
		return err
	}

	defer rows.Close()
	var name, email string
	var id int

	for rows.Next() {
		err := rows.Scan(&id, &name, &email)
		if err != nil {
			log.Println(err)
		}

		fmt.Println("Data : ", id, name, email)
	}
	if err != nil {
		log.Fatal("ERrro reading data : ", err)
	}
	return nil
}

func deleteUserByID(conn *sql.DB, id int) error {
	query := fmt.Sprintf(`DELETE from users where id='%d'`, id)
	_, err := conn.Exec(query)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

func main() {
	conn, err := sql.Open("pgx", "host=localhost port=5432 dbname=blog_db user=postgres password = Admin1234")
	if err != nil {
		log.Fatalf(fmt.Sprintf("Coulnd COnnect to database : %v\n", err))
	}

	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		log.Fatalf("Couldnt ping ")
	}

	err = getAllRowData(conn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("-------------------------------------")

	// err = insertNewUser(conn, "Elio Trapi", "ss@aol.com",
	// 	"abcdef", 3)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = getAllRowData(conn)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	getUserData(conn, 1)
	fmt.Println("-------------------------------------")

	err = updateUserEmail(conn, "db-2@aol.com", 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("-------------------------------------")

	getUserData(conn, 1)
	err = deleteUserByID(conn, 2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("-------------------------------------")

	err = getAllRowData(conn)
	if err != nil {
		log.Fatal(err)
	}

}
