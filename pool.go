package mongo

import (
	"github.com/go-jar/pool"
)

type PoolConfig struct {
	pool.Config

	NewClientFunc func() (*Client, error)
}

type Pool struct {
	pl *pool.Pool

	config *PoolConfig
}

type NewClientFunc func() (*Client, error)

func NewPool(config *PoolConfig) *Pool {
	p := &Pool{
		config: config,
	}

	if p.config.NewConnFunc == nil {
		p.config.NewConnFunc = p.newConn
	}

	p.pl = pool.NewPool(&p.config.Config)

	return p
}

func (p *Pool) Get() (*Client, error) {
	conn, err := p.pl.Get()
	if err != nil {
		return nil, err
	}
	return conn.(*Client), err
}

func (p *Pool) Put(client *Client) error {
	return p.pl.Put(client)
}

func (p *Pool) newConn() (pool.IConn, error) {
	return p.config.NewClientFunc()
}
