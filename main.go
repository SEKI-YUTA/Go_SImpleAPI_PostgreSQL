package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// とりあえずユーザーを追加する機能と一覧を取得する機能は実装した
// Todo 多分DBへの接続をプールできるはず
// Todo 削除、編集機能を追加
type People struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Age int `json:"age"`
}

func getAllUser(ctx *gin.Context) {
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
	// defer rows.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get user name %v\n", err)
		os.Exit(1)
	}

	people := []People{}
	for rows.Next() {
		var p People
		err := rows.Scan(&p.ID, &p.Name, &p.Age)
		if err != nil {
			fmt.Println("failed to scan data")
		}

		people = append(people, p)
	}

	if err := rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr,"failed while iterating rows %v\n", err)
	}

	// fmt.Println(len(people))

	// for _, pp := range people {
	// 	fmt.Println("name: " + pp.name + " age: ", pp.age)
	// }

	ctx.IndentedJSON(http.StatusOK, people)
}

func addUser(ctx *gin.Context) {
	var db_url string = "postgresql://localhost:5432/go_lang?user=root&password=root"
	// https://vamdemicsystem.black/go/%E3%80%90go%E3%80%91go%E3%81%A7postgresql%E3%81%B8%E6%8E%A5%E7%B6%9A%E3%81%99%E3%82%8B
	fmt.Println("Start Application")
	conn, err := pgx.Connect(context.Background(), db_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to db %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var newUser People

	if err := ctx.BindJSON(&newUser); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to bind to json"})
		return
	}

	_, err = conn.Exec(context.Background(), "insert into user_list (name, age) values ($1, $2);", newUser.Name, newUser.Age)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to add user"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, newUser)
}

func main() {
	// fmt.Println("user name: " + first_user)
	router := gin.Default()
	router.GET("/users", getAllUser)
	router.POST("/users/add", addUser)
	router.Run("localhost:9090")
}