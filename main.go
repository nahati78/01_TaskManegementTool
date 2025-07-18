package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	r := gin.Default()                                    //ginの「デフォルトサーバ」（ミドルウェアとかログの初期化済み）をrという変数で操作する.
	r.POST("/signup", signupHandler)                      //signupというURLにアクセスされたらsingupHandler処理を実行しろ（signupというURLの作成も兼ねている）.
	r.POST("/login", loginHandler)                        //ログインAPI
	r.POST("/tasks", addTaskHandler)                      //タスク追加API
	r.PATCH("/tasks/:id/status", updateTaskStatusHandler) //タスクステータス変更API
	r.GET("/tasks", getTasksHandler)                      //タスク一覧取得API
	r.PATCH("/tasks/:id/delete", deleteTaskHandler)       //タスク論理削除API
	r.Run(":8080")                                        //サーバを8080で開く.
}
