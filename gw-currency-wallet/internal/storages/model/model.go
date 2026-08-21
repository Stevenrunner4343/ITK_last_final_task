package model

type Data struct {
	Id        int    `json:"id"`
	Operation string `json:"operation"`
	Amount    int    `json:"amount"`
}

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
