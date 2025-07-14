package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct { //ユーザ定義の構造体
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// 仮：登録済みユーザー（ハードコーディング）
var fakeUser = User{
	Name:     "sample",
	Email:    "sample@example.com",
	Password: "xxxxxx",
}

// タスク用構造体
type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title" binding:"required"`
	About     string `json:"about" binding:"required"`
	Status    int    `json:"status" binding:"required"`
	Limit     string `json:"limit" binding:"required"` // 本来はtime.Time型が望ましい
	CreatedAt string `json:"created_at"`
	UserID    int    `json:"user_id"`
}

func main() {
	r := gin.Default()               //ginの「デフォルトサーバ」（ミドルウェアとかログの初期化済み）をrという変数で操作する.
	r.POST("/signup", signupHandler) //signupというURLにアクセスされたらsingupHandler処理を実行しろ（signupというURLの作成も兼ねている）.
	r.POST("/login", loginHandler)   //ログインAPI追加
	r.POST("/tasks", addTaskHandler) //タスク追加API追加
	r.Run(":8080")                   //サーバを8080で開く.
}

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

	// 仮：ID自動発番（実際はDBで発番）
	newTask := Task{
		ID:        15, // ダミー
		Title:     req.Title,
		About:     req.About,
		Status:    req.Status,
		Limit:     req.Limit,
		CreatedAt: "2025-07-08T12:58:00Z", // ダミー
		UserID:    userID,
	}
	c.JSON(http.StatusOK, newTask)
}
