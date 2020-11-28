package mongo

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/go-jar/pool"
)

const (
	TEST_COLLECTION = "myCollection"
)

func TestPool(t *testing.T) {
	config := &pool.Config{
		MaxConns:    100,
		MaxIdleTime: time.Second * 5,
	}

	pool := NewPool(config, newMongoTestClient)

	testPool(pool, t)
}

func newMongoTestClient() (*Client, error) {
	return newTestClient(), nil
}

func testPool(p *Pool, t *testing.T) {
	client, _ := p.Get()

	query := NewQuery()
	result := []bson.M{}
	err := client.Query(TEST_COLLECTION, query).All(&result)
	if err != nil {
		fmt.Println(err.Error())
	}

	jsonData, _ := json.Marshal(result)
	fmt.Println(jsonData)

	p.Put(client)
}
