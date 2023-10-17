package models

type User struct {
	Id       int    `json:"-"` // hide from json
	Login    string `json:"login"`
	Password string `json:"password"`
}
