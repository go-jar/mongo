package mongo

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/go-jar/golog"
	"github.com/goinbox/gomisc"
)

type Client struct {
	config *Config

	conn       *mgo.Session
	db         *mgo.Database
	collection *mgo.Collection

	isConnected  bool
	isConnClosed bool

	logger    golog.ILogger
	traceId   []byte
	logPrefix []byte
}

func NewClient(config *Config, logger golog.ILogger) *Client {
	if logger == nil {
		logger = new(golog.NoopLogger)
	}

	return &Client{
		config: config,

		logger:    logger,
		traceId:   []byte("-"),
		logPrefix: []byte("[mongo " + strings.Join(config.Hosts, ",") + "]\t"),
	}
}

func (c *Client) SetTraceId(traceId []byte) *Client {
	c.traceId = traceId

	return c
}

func (c *Client) SetDebug(debug bool) {
	mgo.SetDebug(debug)
}

func (c *Client) IsConnected() bool {
	return c.isConnected
}

func (c *Client) IsClosed() bool {
	return c.isConnClosed
}

func (c *Client) Free() {
	if c.conn != nil {
		c.conn.Close()
	}

	c.isConnected = false
	c.isConnClosed = true
}

func (c *Client) Connect() error {
	url := "mongodb://"
	if c.config.User == "" && c.config.Passwd == "" {
		url += strings.Join(c.config.Hosts, ",")
	} else {
		url += c.config.User + ":" + c.config.Passwd + "@" + strings.Join(c.config.Hosts, ",")
	}

	session, err := mgo.Dial(url)
	if err != nil {
		return err
	}

	session.SetMode(c.config.Mode, true)
	session.SetSocketTimeout(c.config.SocketTimeout)
	session.SetSyncTimeout(c.config.SyncTimeout)

	c.conn = session
	c.db = session.DB(c.config.DBName)
	c.isConnected = true
	return nil
}

func (c *Client) connectCheck() {
	if !c.isConnected {
		if err := c.Connect(); err != nil {
			panic(err)
		}
	}
}

func (c *Client) DB(name string) *mgo.Database {
	c.db = c.conn.DB(name)
	return c.db
}

func (c *Client) Collection(collectionName string) *mgo.Collection {
	c.connectCheck()
	c.collection = c.db.C(collectionName)
	return c.collection
}

func (c *Client) Count(collectionName string) (int, error) {
	n, err := c.Collection(collectionName).Count()
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (c *Client) BuildQuery(collectionName string, query *Query) *mgo.Query {
	q := c.Collection(collectionName).Find(query.finder).SetMaxTime(c.config.QueryTimeout)

	if query.selector != nil {
		q = q.Select(query.selector)
	}

	if query.sort != nil {
		q = q.Sort(query.sort...)
	}

	if query.limit != 0 {
		q = q.Limit(query.limit)
	}

	if query.skip != 0 {
		q = q.Skip(query.skip)
	}

	if query.setMaxTime != 0 {
		q = q.SetMaxTime(query.setMaxTime)
	}

	return q
}

func (c *Client) Query(collectionName string, query *Query) *mgo.Query {
	c.log("Query", collectionName, query)
	return c.BuildQuery(collectionName, query)
}

func (c *Client) Find(collectionName string, finder interface{}) *mgo.Query {
	c.log("Find", finder)
	return c.Collection(collectionName).Find(finder)
}

func (c *Client) FindId(collectionName string, id interface{}) *mgo.Query {
	c.log("FindId", collectionName, id)
	return c.Collection(collectionName).FindId(id)
}

func (c *Client) FindAndModify(collectionName string, finder interface{}, updater interface{}) (bson.M, error) {
	c.log("FindAndModify", collectionName, finder, updater)
	change := mgo.Change{
		Update:    updater,
		Upsert:    true,
		ReturnNew: true,
	}

	result := bson.M{}
	_, err := c.Collection(collectionName).Find(finder).Apply(change, result)
	return result, err
}

func (c *Client) Indexes(collectionName string) ([]mgo.Index, error) {
	c.log("Indexs", collectionName)
	indexes, err := c.Collection(collectionName).Indexes()
	return indexes, err
}

func (c *Client) Insert(collectionName string, docs ...interface{}) error {
	c.log("Insert", collectionName, docs)
	err := c.Collection(collectionName).Insert(docs...)
	return err
}

func (c *Client) Update(collectionName string, selector, updater interface{}) error {
	c.log("Update", collectionName, selector, updater)
	err := c.Collection(collectionName).Update(selector, updater)
	return err
}

func (c *Client) UpdateAll(collectionName string, selector, updater interface{}) error {
	c.log("UpdateAll", collectionName, selector, updater)
	_, err := c.Collection(collectionName).UpdateAll(selector, updater)
	return err
}

func (c *Client) UpdateId(collectionName string, id, updater interface{}) error {
	c.log("UpdateId", collectionName, id, updater)
	err := c.Collection(collectionName).UpdateId(id, updater)
	return err
}

func (c *Client) Upsert(collectionName string, selector, updater interface{}) error {
	c.log("Upsert", collectionName, selector, updater)
	_, err := c.Collection(collectionName).Upsert(selector, updater)
	return err
}

func (c *Client) Remove(collectionName string, selector interface{}) error {
	c.log("Remove", collectionName, selector)
	err := c.Collection(collectionName).Remove(selector)
	return err
}

func (c *Client) RemoveAll(collectionName string, selector interface{}) error {
	c.log("RemoveAll", collectionName, selector)
	_, err := c.Collection(collectionName).RemoveAll(selector)
	return err
}

func (c *Client) RemoveId(collectionName string, id interface{}) error {
	c.log("RemoveId", collectionName, id)
	err := c.Collection(collectionName).RemoveId(id)
	return err
}

func (c *Client) ConvertBsonToStruct(doc interface{}, entity interface{}) error {
	bsonBytes, _ := bson.Marshal(doc)
	err := bson.Unmarshal(bsonBytes, entity)
	return err
}

func (c *Client) ConvertJsonToStruct(doc interface{}, entity interface{}) error {
	jsonBytes, _ := json.Marshal(doc)
	err := json.Unmarshal(jsonBytes, entity)
	return err
}

func (c *Client) log(query string, args ...interface{}) {
	for _, arg := range args {
		query += " " + fmt.Sprint(arg)
	}

	_ = c.logger.Log(c.config.LogLevel, gomisc.AppendBytes(c.logPrefix, []byte("\t"), c.traceId, []byte("\t"), []byte(query)))
}
