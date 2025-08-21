package models

type TaskIdFromRequest struct {
	ID      int64  `json:"task_ID"`
	User_ID int64  `json:"user_ID"`
	Status  string `json:"status"`
}

type MoneyDataBase struct {
	Goal   int `json:"goal"`
	Amount int `json:"amount"`
}

type Stats struct {
	MoneyLeft int `json:"money_left"`
	Sum       int `json:"goal"`
	MoneyToGo int `json:"current_money"`
}
