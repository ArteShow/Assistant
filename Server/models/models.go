package models

type TaskIdFromRequest struct {
	ID      int64 `json:"task_ID"`
	User_ID int64 `json:"user_ID"`
}
