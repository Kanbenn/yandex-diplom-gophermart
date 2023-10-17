package models

type User struct {
	ID       int    `json:"-"` // hide from json
	Login    string `json:"login"`
	Password string `json:"password"`
}
