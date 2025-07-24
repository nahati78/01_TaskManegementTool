package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var tasks []Task // 仮：タスク一覧をメモリ上に保持（DB未接続のため）

// ユーザ登録APIのハンドラ関数
func signupHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	user := User{
		Name:     req.Username,
		Email:    fmt.Sprintf("%s@example.com", req.Username),
		Password: req.Password,
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Could not create user"})
		return
	}
	c.JSON(200, gin.H{"message": "User created"})
}

// ログインAPI
func loginHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	var user User
	if err := db.Where("email = ? AND password = ?", req.Username, req.Password).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}
	c.JSON(200, gin.H{"message": "Login successful"})
}

// タスク追加API
func addTaskHandler(c *gin.Context) {
	// 仮：認証済みのユーザID=1とする（本来はJWTから取得）
	userID := 1

	var req struct {
		Title  string `json:"title" binding:"required"`
		About  string `json:"about" binding:"required"`
		Status int    `json:"status" binding:"required"`
		Limit  string `json:"limit" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	// ステータス値のバリデーション（1/2/3以外はNG）
	if req.Status < 1 || req.Status > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status must be 1, 2, or 3"})
		return
	}

	// 仮：ID自動発番（実際はDB）
	newTask := Task{
		ID:        len(tasks) + 1, // IDを自動採番
		Title:     req.Title,
		About:     req.About,
		Status:    req.Status,
		Limit:     req.Limit,
		CreatedAt: "2025-07-08T12:58:00Z", // ダミー
		UserID:    userID,
	}
	tasks = append(tasks, newTask) // tasks配列に追加
	c.JSON(http.StatusOK, newTask)
}

// タスクステータス変更API
func updateTaskStatusHandler(c *gin.Context) {
	// URLパラメータからid取得
	idParam := c.Param("id")
	var statusReq struct {
		Status int `json:"status" binding:"required"`
	}
	// JSONバリデーション
	if err := c.BindJSON(&statusReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	// ステータス値のバリデーション
	if statusReq.Status < 1 || statusReq.Status > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status must be 1, 2, or 3"})
		return
	}
	// idParamをintに変換
	var targetID int
	_, err := fmt.Sscanf(idParam, "%d", &targetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	// タスク検索・更新
	for i := range tasks {
		if tasks[i].ID == targetID {
			tasks[i].Status = statusReq.Status
			c.JSON(http.StatusOK, gin.H{"id": tasks[i].ID, "status": tasks[i].Status})
			return
		}
	}
	// タスク見つからなければ404
	c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
}

// タスク一覧取得API
func getTasksHandler(c *gin.Context) {
	c.JSON(http.StatusOK, tasks)
}

// タスク論理削除API
func deleteTaskHandler(c *gin.Context) {
	idParam := c.Param("id")
	var targetID int
	_, err := fmt.Sscanf(idParam, "%d", &targetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	// タスクを検索して論理削除
	for i := range tasks {
		if tasks[i].ID == targetID {
			tasks[i].Status = 9 // 論理削除
			c.JSON(http.StatusOK, gin.H{"id": tasks[i].ID, "status": tasks[i].Status})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
}
