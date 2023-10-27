package models

type User struct {
	ID       int    `json:"-"` // hide from json.Marshal
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserBalance struct {
	ID        int     `json:"-"`
	Balance   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

type Order struct {
	Number string  `json:"order"`
	User   int     `json:"-"`
	Status string  `json:"-"`
	Sum    float32 `json:"sum" `
	Time   string  `json:"processed_at"`
}

type OrderResponse struct {
	Number string  `json:"number"`
	Status string  `json:"status"`
	Bonus  float32 `json:"accrual,omitempty"`
	Time   string  `json:"uploaded_at"`
}

type AccrualResponse struct {
	User   int     `json:"-" db:"user_id"`
	Number string  `json:"order"`
	Status string  `json:"status"`
	Bonus  float32 `json:"accrual"`
}
