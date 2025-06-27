package models

type Task struct {
	ID     int64  `json:"id"`
	Desc   string `json:"desc"`
	Status bool   `json:"status"`
	UserID int64  `json:"user_id"`
}
