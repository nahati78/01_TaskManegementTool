package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/signup", func(c *gin.Context) {
		// リクエストボディを受け取る
		type User struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		var user User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		// 仮にID=1を返すだけ（DB未接続の仮実装）
		c.JSON(http.StatusOK, gin.H{
			"id":    1,
			"name":  user.Name,
			"email": user.Email,
		})
	})
	r.Run(":8080")
}
