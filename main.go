package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("./test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	// モデルをマイグレーション
	db.AutoMigrate(&User{}, &Task{})
}

func main() {
	initDB()
	db.Create(&User{Name: "sample", Email: "sample@example.com", Password: "xxxxxx"})
	db.Create(&Task{Title: "AWS環境設定", About: "AWSへの登録/WSLによる操作", Status: 1, Limit: "2025-07-15", CreatedAt: time.Now().Format(time.RFC3339), UserID: 1}) // データベースの初期化
	r := gin.Default()                                                                                                                                   //ginの「デフォルトサーバ」（ミドルウェアとかログの初期化済み）をrという変数で操作する.
	r.POST("/signup", signupHandler)
	r.POST("/login", loginHandler)

	auth := r.Group("/", AuthMiddleware())
	{
		auth.POST("/tasks", addTaskHandler)
		auth.PATCH("/tasks/:id/status", updateTaskStatusHandler)
		auth.GET("/tasks", getTasksHandler)
		auth.PATCH("/tasks/:id/delete", deleteTaskHandler)
	}
}
