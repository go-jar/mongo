package mongo

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/globalsign/mgo/bson"
)

const (
	TEST_COLLECTION = "myCollection"
)

func init() {
	client = newTestClient()
}

func TestClient_Remove(t *testing.T) {
	selector := bson.M{"_id": 4}
	err := client.Remove(TEST_COLLECTION, selector)
	if err != nil {
		fmt.Println(err)
	}
}

func TestClient_RemoveId(t *testing.T) {
	id := 3
	err := client.RemoveId(TEST_COLLECTION, id)
	if err != nil {
		fmt.Println(err)
	}
}

func TestClient_RemoveAll(t *testing.T) {
	selector := bson.M{"_id": bson.M{"$gte": 0}}
	err := client.RemoveAll(TEST_COLLECTION, selector)
	if err != nil {
		fmt.Println(err)
	}
}

func TestClient_Insert(t *testing.T) {
	doc := bson.M{"_id": 11, "a": 1, "b": 2}
	err := client.Insert(TEST_COLLECTION, doc)
	if err != nil {
		fmt.Println(err)
	}

	total := 5
	docs := make([]interface{}, total)
	for i := 0; i < total; i++ {
		docs[i] = bson.M{"_id": i, "a": 3, "b": 4}
	}

	err = client.Insert(TEST_COLLECTION, docs...)
	if err != nil {
		fmt.Println(err)
	}
}

func TestClient_Update(t *testing.T) {
	selector := bson.M{"_id": 1}
	updater := bson.M{
		"$inc":         bson.M{"view_count": 1},
		"$currentDate": bson.M{"edit_time": true},
	}
	err := client.Update(TEST_COLLECTION, selector, updater)
	if err != nil {
		fmt.Println(err)
	}
}

func TestClient_UpdateAll(t *testing.T) {
	selector := bson.M{"_id": bson.M{"$gt": 0}}
	updater := bson.M{
		"$inc":         bson.M{"view_count": 1},
		"$currentDate": bson.M{"edit_time": true},
	}
	err := client.UpdateAll(TEST_COLLECTION, selector, updater)
	if err != nil {
		fmt.Println(err)
	}
}

func TestClient_UpdateId(t *testing.T) {
	id := 1
	updator := bson.M{
		"$inc":         bson.M{"view_count": 1},
		"$currentDate": bson.M{"edit_time": true},
	}
	err := client.UpdateId(TEST_COLLECTION, id, updator)
	if err != nil {
		fmt.Println(err)
	}
}

func TestClient_Upsert(t *testing.T) {
	selector := bson.M{"_id": 4}
	updator := bson.M{
		"$inc":         bson.M{"view_count": 1},
		"$currentDate": bson.M{"edit_time": true},
		"$setOnInsert": bson.M{"add_time": "2020-11-28 08:00:99"},
	}
	err := client.Upsert(TEST_COLLECTION, selector, updator)
	if err != nil {
		fmt.Println(err)
	}
}

func TestClient_Query(t *testing.T) {
	result := []bson.M{}
	query := NewQuery().
		Find(bson.M{"_id": bson.M{"$gt": 0}}).
		Sort("-_id").
		Select(bson.M{"edit_time": 0}).
		Skip(0).
		SetMaxTime(1 * time.Second)

	err := client.Query(TEST_COLLECTION, query).One(&result)
	if err != nil {
		fmt.Println(err)
	}

	for _, item := range result {
		jsonData, _ := json.Marshal(item)
		fmt.Printf("%s\n", jsonData)
	}
}

func TestClient_QueryOne(t *testing.T) {
	result := bson.M{}
	query := NewQuery().
		Find(bson.M{"_id": bson.M{"$gt": 0}}).
		Sort("-_id").
		Select(bson.M{"edit_time": 0}).
		Skip(0).
		SetMaxTime(1 * time.Second)

	err := client.Query(TEST_COLLECTION, query).One(&result)
	if err != nil {
		t.Error(err)
	}

	jsonData, _ := json.Marshal(result)
	fmt.Printf("%s\n", jsonData)
}

func TestClient_Count(t *testing.T) {
	query := NewQuery()
	result, err := client.Query(TEST_COLLECTION, query).Count()
	if err != nil {
		fmt.Println(err)
	}
	jsonData, _ := json.Marshal(result)
	fmt.Printf("%s\n", jsonData)
}

func TestClient_Find(t *testing.T) {
	result := []bson.M{}
	err := client.Find(TEST_COLLECTION, bson.M{"_id": bson.M{"$in": []int{4, 11}}}).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	jsonData, _ := json.Marshal(result)
	fmt.Printf("%s\n", jsonData)
}

func TestClient_FindId(t *testing.T) {
	result := bson.M{}
	err := client.FindId(TEST_COLLECTION, 4).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	jsonData, _ := json.Marshal(result)
	fmt.Printf("%s\n", jsonData)
}

func TestClient_FindAndModify(t *testing.T) {
	finder := bson.M{"_id": "app"}
	updater := bson.M{"$inc": bson.M{"max_id": 1}}
	result, err := client.FindAndModify("id_genter", finder, updater)
	if err != nil {
		fmt.Println(err)
	}

	doc := new(TestIdGenerateStruct)
	err = client.ConvertBsonToStruct(result, doc)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(doc)

	err = client.ConvertJsonToStruct(result, doc)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(doc)
}
