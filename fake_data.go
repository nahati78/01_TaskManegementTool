package main

// 仮：登録済みユーザー（ハードコーディング）
var fakeUser = User{
	Name:     "sample",
	Email:    "sample@example.com",
	Password: "xxxxxx",
}

// 仮：タスク一覧をメモリ上に保持（DB未接続のため）
var tasks = []Task{
	{ID: 15, Title: "AWS環境設定", About: "AWSへの登録/WSLによる操作", Status: 1, Limit: "2025-07-15", CreatedAt: "2025-07-08T12:58:00Z", UserID: 1},
	// ...他にもタスクがあれば追加
}
