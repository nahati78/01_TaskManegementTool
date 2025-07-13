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

func main() {
	r := gin.Default()               //ginの「デフォルトサーバ」（ミドルウェアとかログの初期化済み）をrという変数で操作する.
	r.POST("/signup", signupHandler) //signupというURLにアクセスされたらsingupHandler処理を実行しろ（signupというURLの作成も兼ねている）.
	r.POST("/login", loginHandler)   // ログインAPI追加
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
	if err := c.BindJSON(&loginData); err != nil {
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
