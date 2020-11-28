package mongo

import (
	"fmt"
	"testing"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/goinbox/gomisc"
)

const (
	TEST_COLLECTION_DAO = "myCollectionDao"
)

type mongoTestEntity struct {
	Id       int64  `bson:"_id" json:"_id"`
	Name     string `bson:"name" json:"name"`
	Status   int    `bson:"status" json:"status"`
	AddTime  string `bson:"add_time" json:"add_time"`
	EditTime string `bson:"edit_time" json:"edit_time"`
}

func TestDaoWrite(t *testing.T) {
	dao := &Dao{newTestClient()}

	selector := bson.M{"_id": bson.M{"$gte": 0}}
	err := dao.RemoveAll(TEST_COLLECTION_DAO, selector)
	if err != nil {
		fmt.Println(err)
	}

	colNames := []string{"_id", "add_time", "edit_time", "name", "status"}

	ts := time.Now().Format(gomisc.TimeGeneralLayout())
	colValues := []interface{}{
		11,
		ts,
		ts,
		"a",
		time.Now().Unix() % 10,
	}
	err = dao.InsertRow(TEST_COLLECTION_DAO, colNames, colValues)
	if err != nil {
		fmt.Println(err)
	}

	var colsValues [][]interface{}
	for i, name := range []string{"a", "b", "c"} {
		colValues := []interface{}{
			int64(i + 100),
			ts,
			ts,
			name,
			i % 10,
		}
		colsValues = append(colsValues, colValues)
	}
	err = dao.InsertRows(TEST_COLLECTION_DAO, colNames, colsValues...)
	if err != nil {
		fmt.Println(err)
	}

	err = dao.UpdateById(TEST_COLLECTION_DAO, 101, map[string]interface{}{"name": "hh"})
	if err != nil {
		fmt.Println(err)
	}
}

func TestDaoRead(t *testing.T) {
	dao := &Dao{newTestClient()}

	result, err := dao.SelectById(TEST_COLLECTION_DAO, 11)
	fmt.Println(result, err)
	entity := new(mongoTestEntity)
	err = dao.ConvertBsonToStruct(result, entity)
	fmt.Println(entity, err)

	results, err := dao.SelectByIds(TEST_COLLECTION_DAO, 11, 101)
	fmt.Println(results, err)

	entities := new([]mongoTestEntity)
	err = dao.ConvertJsonToStruct(results, entities)
	fmt.Println(entities, err)

	count, err := dao.SelectCount(TEST_COLLECTION_DAO, bson.M{"_id": bson.M{"$gt": 11}})
	fmt.Println(count, err)

	queryAll := NewQuery().
		Find(bson.M{"_id": bson.M{"$gte": 100}}).
		Sort("_id").
		Select(bson.M{"edit_time": 0}).Skip(0).
		SetMaxTime(1 * time.Second)
	rows, err := dao.SelectAll(TEST_COLLECTION_DAO, queryAll)
	fmt.Println(rows, err)

	row, err := dao.SelectOne(TEST_COLLECTION_DAO, queryAll)
	fmt.Println(row, err)
}
