package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// とりあえずユーザーを追加する機能と一覧を取得する機能は実装した
// Todo 多分DBへの接続をプールできるはず
// Todo 削除、編集機能を追加

// メモ プールを使用せずにDB接続する場合
// var db_url string = "postgresql://{ホスト名}:{ポート番号}/{DB名}?user={ユーザ名}&password={パスワード}"
// var db_url string = "postgresql://localhost:5432/go_lang?user=root&password=root"
// // https://vamdemicsystem.black/go/%E3%80%90go%E3%80%91go%E3%81%A7postgresql%E3%81%B8%E6%8E%A5%E7%B6%9A%E3%81%99%E3%82%8B
// fmt.Println("Start Application")
// conn, err := pgx.Connect(context.Background(), db_url)
// if err != nil {
// 	fmt.Fprintf(os.Stderr, "Unable to connect to db %v\n", err)
// 	os.Exit(1)
// }
// defer conn.Close(context.Background())

type People struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Age int `json:"age"`
}

var pool *pgxpool.Pool

func getAllUser() []People {
	rows, err := pool.Query(context.Background(),"select * from user_list;")
	// rows, err := conn.Query(context.Background(),"select * from user_list;")
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

	return people
}

func responseAllUser(ctx *gin.Context) {
	var peopleList []People
	peopleList = getAllUser()

	ctx.IndentedJSON(http.StatusOK, peopleList)
}

func addUser(ctx *gin.Context) {
	var newUser People

	if err := ctx.BindJSON(&newUser); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to bind to json"})
		return
	}

	_, err := pool.Exec(context.Background(), "insert into user_list (name, age) values ($1, $2);", newUser.Name, newUser.Age)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to add user"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, newUser)
}

func getUserById(id int) (People) {
	var allUser = getAllUser()
	var target People
	for _, u := range allUser {
		if u.ID == id {
			target = u
		}
	}
	return target
}

func editUser(ctx *gin.Context) {
	id := ctx.Param("id")
	idNum, err := strconv.Atoi(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to update user infomation"})
		return
	}
	user := getUserById(idNum)
	
	// 編集対象のユーザーが見つからなかった場合(ID0は存在しないため0が帰ってきたら取得できていないと判断)
	if  user.ID == 0 {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to update user infomation"})
		return
	}
	
	// これが新しい状態のPeopleを入れる変数
	var editedUser People
	if err := ctx.BindJSON(&editedUser); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to update user infomation"})
		return
	}
	
	pool.Exec(context.Background(), "update user_list set name=$1, age=$2 where id=$3",
	editedUser.Name, editedUser.Age, user.ID)

	ctx.IndentedJSON(http.StatusOK, editedUser)

}

func deleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	idNum, err := strconv.Atoi(id)

	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to delete user"})
		return
	}
	
	var user = getUserById(idNum)
	if user.ID == 0 {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to delete user"})
		return
	}
	_, err = pool.Exec(context.Background(), "delete from user_list where id=$1", user.ID)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to delete user"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, user)
}

func responseUserById(ctx *gin.Context) {
	id := ctx.Param("id")
	idNum, err := strconv.Atoi(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to get user"})
		return
	}

	user := getUserById(idNum)

	ctx.IndentedJSON(http.StatusOK, user)
}

func main() {
	// fmt.Println("user name: " + first_user)
	connConfig, err := pgx.ParseConfig("postgresql://root:root@localhost:5432/go_lang")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse db config")
		os.Exit(1)
	}

	poolConfig, err := pgxpool.ParseConfig("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse pool config")
		os.Exit(1)
	}
	poolConfig.ConnConfig = connConfig

	pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to db")
		os.Exit(1)
	}

	defer pool.Close()

	fmt.Println("start app")
	router := gin.Default()
	router.GET("/users", responseAllUser)
	router.GET("/users/:id", responseUserById)
	router.POST("/users/add", addUser)
	router.PATCH("/users/edit/:id", editUser)
	router.PATCH("/users/delete/:id", deleteUser)
	router.Run("localhost:9090")
	fmt.Println("end app")
}