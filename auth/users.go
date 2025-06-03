package auth

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

var Users = []User{
	{Login: "dima",
		Password: "123"},
}
