package model

type Api struct {
	Url    string `json:"url" bson:"url"`
	Method string `json:"method" bson:"method"`
}
