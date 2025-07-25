package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var tasks []Task // 仮：タスク一覧をメモリ上に保持（DB未接続のため）
var jwtKey = []byte("your-very-secret-key")

// taken発行関数
func generateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

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
		Username: req.Username, // 追加！
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
	if err := db.Where("username = ? AND password = ?", req.Username, req.Password).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}
	// JWT発行
	token, err := generateJWT(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "JWT generation failed"})
		return
	}
	c.JSON(200, gin.H{"token": token})
}

// タスク追加API
func addTaskHandler(c *gin.Context) {

	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	uid := userID.(int) // ユーザIDを取得

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

	newTask := Task{
		Title:     req.Title,
		About:     req.About,
		Status:    req.Status,
		Limit:     req.Limit,
		CreatedAt: time.Now().Format(time.RFC3339),
		UserID:    uid, // ユーザIDを設定
	}
	db.Create(&newTask)
	c.JSON(http.StatusOK, newTask)
}

// タスクステータス変更API
func updateTaskStatusHandler(c *gin.Context) {
	// 認証ユーザー取得
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	uid := userID.(int)

	idParam := c.Param("id")
	var statusReq struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.BindJSON(&statusReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	if statusReq.Status < 1 || statusReq.Status > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status must be 1, 2, or 3"})
		return
	}

	var task Task
	if err := db.First(&task, idParam).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	// セキュリティ：自分のタスクのみ更新OK
	if task.UserID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your task"})
		return
	}

	task.Status = statusReq.Status
	db.Save(&task)
	c.JSON(http.StatusOK, gin.H{"id": task.ID, "status": task.Status})
}

// タスク一覧取得API
func getTasksHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	uid := userID.(int)

	// DB使用の場合（GORM）
	/*
	   var userTasks []Task
	   if err := db.Where("user_id = ? AND status <> 9", uid).Find(&userTasks).Error; err != nil {
	       c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
	       return
	   }
	   c.JSON(http.StatusOK, userTasks)
	*/

	var userTasks []Task
	for _, t := range tasks {
		if t.UserID == uid && t.Status != 9 {
			userTasks = append(userTasks, t)
		}
	}
	c.JSON(http.StatusOK, userTasks)

}

// タスク論理削除API
func deleteTaskHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	uid := userID.(int)
	/*
	   idParam := c.Param("id")
	   var task Task
	   if err := db.First(&task, idParam).Error; err != nil {
	       c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
	       return
	   }
	   // セキュリティ：自分のタスクだけ削除OK
	   if task.UserID != uid {
	       c.JSON(http.StatusForbidden, gin.H{"error": "not your task"})
	       return
	   }
	   task.Status = 9 // 論理削除
	   db.Save(&task)
	   c.JSON(http.StatusOK, gin.H{"id": task.ID, "status": task.Status})
	*/

	idParam := c.Param("id")
	var targetID int
	_, err := fmt.Sscanf(idParam, "%d", &targetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	for i := range tasks {
		if tasks[i].ID == targetID {
			if tasks[i].UserID != uid {
				c.JSON(http.StatusForbidden, gin.H{"error": "not your task"})
				return
			}
			tasks[i].Status = 9
			c.JSON(http.StatusOK, gin.H{"id": tasks[i].ID, "status": tasks[i].Status})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})

}

// JWT認証ミドルウェア
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header required"})
			return
		}
		tokenString := authHeader[len("Bearer "):]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid claims"})
			return
		}
		c.Set("user_id", int(claims["user_id"].(float64)))
		c.Next()
	}
}
