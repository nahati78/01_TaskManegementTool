package main

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title" binding:"required"`
	About     string `json:"about" binding:"required"`
	Status    int    `json:"status" binding:"required"`
	Limit     string `json:"limit" binding:"required"`
	CreatedAt string `json:"created_at"`
	UserID    int    `json:"user_id"`
}
