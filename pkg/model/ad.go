package model

const (
	AdCollection = "ads"
)

type Ad struct {
	Id          string `bson:"id" json:"id"`
	Email       string `bson:"email" json:"email"`
	Description string `bson:"description" json:"description"`
	State       string `bson:"state" json:"state"`
	Category    string `bson:"category" json:"category"`
	Image       string `bson:"image" json:"image"`
}
