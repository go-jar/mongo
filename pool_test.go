package mongo

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/globalsign/mgo/bson"
)

func TestPool(t *testing.T) {
	config := &PoolConfig{NewClientFunc: newMongoTestClient}
	config.MaxConns = 100
	config.MaxIdleTime = time.Second * 5

	pool := NewPool(config)

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
