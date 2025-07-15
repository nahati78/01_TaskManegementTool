package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()               //ginの「デフォルトサーバ」（ミドルウェアとかログの初期化済み）をrという変数で操作する.
	r.POST("/signup", signupHandler) //signupというURLにアクセスされたらsingupHandler処理を実行しろ（signupというURLの作成も兼ねている）.
	r.POST("/login", loginHandler)   //ログインAPI追加
	r.POST("/tasks", addTaskHandler) //タスク追加API追加
	r.PATCH("/tasks/:id/status", updateTaskStatusHandler)
	r.Run(":8080") //サーバを8080で開く.
}
