package entity

type Car struct {
	Id      string `json:"id" bson:"_id"`
	UserId  string `json:"-" bson:"userId"`
	Name    string `json:"name" bson:"name"`
	Company string `json:"company" bson:"company"`
	Model   int    `json:"model" bson:"model"`
	IsAdmin bool   `json:"isAdmin"`
}

type NewCar struct {
	Name    string `json:"name" bson:"name"`
	Company string `json:"company" bson:"company"`
	Model   int    `json:"model" bson:"model"`
}
