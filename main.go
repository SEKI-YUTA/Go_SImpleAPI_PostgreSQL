package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type People struct {
	id int
	name string
	age int
}

func main() {
	// var db_url string = "postgresql://{ホスト名}:{ポート番号}/{DB名}?user={ユーザ名}&password={パスワード}"
	var db_url string = "postgresql://localhost:5432/go_lang?user=root&password=root"
	// https://vamdemicsystem.black/go/%E3%80%90go%E3%80%91go%E3%81%A7postgresql%E3%81%B8%E6%8E%A5%E7%B6%9A%E3%81%99%E3%82%8B
	fmt.Println("Start Application")
	conn, err := pgx.Connect(context.Background(), db_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to db %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// var first_user string
	// var first_user user

	rows, err := conn.Query(context.Background(),"select * from user_list;")
	// err = conn.QueryRow(context.Background(), "select name from user_list where id = 1;").Scan(&first_user)
	defer rows.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get user name %v\n", err)
		os.Exit(1)
	}

	people := []People{}
	for rows.Next() {
		var p People
		err := rows.Scan(&p.id, &p.name, &p.age)
		if err != nil {
			fmt.Println("failed to scan data")
		}

		people = append(people, p)
	}

	if err := rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr,"failed while iterating rows %v\n", err)
	}

	fmt.Println(len(people))

	for _, pp := range people {
		fmt.Println("name: " + pp.name + " age: ", pp.age)
	}
	// fmt.Println("user name: " + first_user)
}