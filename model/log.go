package model

type Log struct {
	Level  string `bson:"level" json:"level"`
	Time   string `bson:"time" json:"ts"`
	Caller string `bson:"caller" json:"caller"`
	Msg    string `bson:"msg" json:"msg"`
}
