package mongo

import "github.com/go-jar/golog"

type TestIdGenerateStruct struct {
	Id    string `bson:"_id"`
	MaxId int32  `bson:"max_id"`
}

var client *Client

func newTestClient() *Client {
	logger, _ := golog.NewSimpleLogger(golog.NewConsoleWriter(), golog.NewConsoleFormat(golog.NewFileInfoFormat(1)))
	config := NewConfig([]string{"127.0.0.1:27017"}, "user", "passwd", "demo")
	return NewClient(config, logger)
}
