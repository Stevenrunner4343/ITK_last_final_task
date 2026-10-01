package storage

type Data struct {
	Id        int    `json:"id"`
	Operation string `json:"operation"`
	Amount    int    `json:"amount"`
}
