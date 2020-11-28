package mongo

import (
	"time"
)

type Query struct {
	finder     interface{}
	selector   interface{}
	sort       []string
	limit      int
	skip       int
	setMaxTime time.Duration
}

func NewQuery() *Query {
	return &Query{}
}

func (q *Query) Find(finder interface{}) *Query {
	q.finder = finder
	return q
}

func (q *Query) Select(selector interface{}) *Query {
	q.selector = selector
	return q
}

func (q *Query) Sort(fields ...string) *Query {
	q.sort = fields
	return q
}

func (q *Query) Limit(n int) *Query {
	q.limit = n
	return q
}

func (q *Query) Skip(n int) *Query {
	q.skip = n
	return q
}

func (q *Query) SetMaxTime(maxTime time.Duration) *Query {
	q.setMaxTime = maxTime
	return q
}
