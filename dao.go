package mongo

import (
	"github.com/globalsign/mgo/bson"
)

const (
	COND_EQUAL         = "$eq"
	COND_NOT_EQUAL     = "$ne"
	COND_LESS          = "$lt"
	COND_LESS_EQUAL    = "$lte"
	COND_GREATER       = "$gt"
	COND_GREATER_EQUAL = "$gte"
	COND_IN            = "$in"
	COND_NOT_IN        = "$nin"
	COND_REGEX         = "$regex"
	COND_OPTIONS       = "$options"
)

type Dao struct {
	*Client
}

func NewDao(client *Client) *Dao {
	return &Dao{
		client,
	}
}

func (d *Dao) InsertRow(tableName string, colNames []string, colValues []interface{}) error {
	doc := make(map[string]interface{})
	for i, colName := range colNames {
		doc[colName] = colValues[i]
	}
	return d.Insert(tableName, doc)
}

func (d *Dao) InsertRows(tableName string, colNames []string, colsValues ...[]interface{}) error {
	docs := make([]interface{}, len(colsValues))
	for i, values := range colsValues {
		doc := make(map[string]interface{})
		for j, colName := range colNames {
			doc[colName] = values[j]
		}
		docs[i] = doc
	}
	return d.Insert(tableName, docs...)
}

func (d *Dao) DeleteById(tableName string, id interface{}) error {
	return d.RemoveId(tableName, id)
}

func (d *Dao) UpdateById(tableName string, id interface{}, updater map[string]interface{}) error {
	return d.UpdateId(tableName, id, map[string]interface{}{"$set": updater})
}

func (d *Dao) SelectById(tableName string, id interface{}) (bson.M, error) {
	result := bson.M{}
	err := d.FindId(tableName, id).One(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *Dao) SelectByIds(tableName string, ids ...interface{}) ([]bson.M, error) {
	result := []bson.M{}
	err := d.Find(tableName, bson.M{"_id": bson.M{"$in": ids}}).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *Dao) SelectCount(tableName string, selector interface{}) (int, error) {
	n, err := d.Find(tableName, selector).Count()
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (d *Dao) SelectAll(tableName string, query *Query) ([]bson.M, error) {
	result := []bson.M{}
	err := d.Query(tableName, query).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *Dao) SelectOne(tableName string, query *Query) (bson.M, error) {
	result := bson.M{}
	err := d.Query(tableName, query).One(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
