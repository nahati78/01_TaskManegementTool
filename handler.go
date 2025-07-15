package main

import (
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
)

// ユーザ登録APIのハンドラ関数
func signupHandler(c *gin.Context) { //cはginフレームワークがリクエストごとに自動で生成し渡す.
	var user User                             //構造体の変数を作成.
	if err := c.BindJSON(&user); err != nil { //リクエストボディ（JSON）をuser構造体にマッピングする.失敗したらerrにエラーが入る.
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"}) //HTTP 400エラーとして、{"error": "invalid input"}を返す.
		return
	}
	// 仮実装：ID=1を返すだけ
	c.JSON(http.StatusOK, gin.H{ //HTTP 200でJSONレスポンスを返す.
		"id":    1,
		"name":  user.Name,
		"email": user.Email,
	})
}

// ログインAPI
func loginHandler(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&loginData); err != nil { //BindJSON←json形式のデータをGoの構造体に自動で代入してくれる
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	// 仮：fakeUserと一致するなら成功
	if loginData.Email == fakeUser.Email && loginData.Password == fakeUser.Password {
		c.JSON(http.StatusOK, gin.H{
			"token": "xxxxx.yyyyy.zzzzz", // 仮のトークン
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "メールアドレスまたはパスワードが正しくありません",
		})
	}
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
// PATCH /tasks/:id/status
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
