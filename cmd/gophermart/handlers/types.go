package handlers

type User struct {
	Name string `json:"login"`
	Password string `json:"password"`
}

type Auth struct {
	Key string
}

type UserModel struct {
	ID int
	Name string
	Password string
}