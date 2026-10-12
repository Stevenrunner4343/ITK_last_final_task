package model

type Data struct {
	Id        int    `bson:"id"`
	Operation string `bson:"operation"`
	Amount    int    `bson:"amount"`
}
