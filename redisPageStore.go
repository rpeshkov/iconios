package main

import (
	"fmt"
	"strconv"

	"github.com/mediocregopher/radix.v2/pool"
	uuid "github.com/satori/go.uuid"
)

type RedisPagestore struct {
	connPool *pool.Pool
}

func (r *RedisPagestore) FinishPage(p *PageData) error {
	c, err := r.connPool.Get()

	if err != nil {
		return err
	}

	defer r.connPool.Put(c)

	key := "page:" + p.ID.String()
	exists, err := c.Cmd("EXISTS", key).Int()

	if err != nil {
		return err
	}

	if exists == 0 {
		return fmt.Errorf("Page '%s' not exist", key)
	}

	return c.Cmd("HSET", key, "finished", "1").Err
}

func (r *RedisPagestore) NewPage(url string, finished bool) (*PageData, error) {
	p := &PageData{uuid.NewV4(), url, finished}
	if err := r.SavePage(p); err != nil {
		return nil, err
	}

	return p, nil
}

func (r *RedisPagestore) SavePage(p *PageData) error {
	c, err := r.connPool.Get()

	if err != nil {
		return err
	}

	defer r.connPool.Put(c)

	key := "page:" + p.ID.String()
	return c.Cmd("HMSET", key, "url", p.URL, "finished", false).Err
}

func (r *RedisPagestore) GetPage(id string) (*PageData, error) {
	c, err := r.connPool.Get()

	if err != nil {
		return nil, err
	}

	defer r.connPool.Put(c)

	key := "page:" + id
	exists, err := c.Cmd("EXISTS", key).Int()

	if err != nil {
		return nil, err
	}

	if exists == 0 {
		return nil, fmt.Errorf("Page '%s' does not exist", key)
	}

	m, err := c.Cmd("HGETALL", key).Map()
	if err != nil {
		return nil, fmt.Errorf("get page data failed: %v", err)
	}

	f, err := strconv.ParseBool(m["finished"])
	if err != nil {
		f = false
	}

	return &PageData{ID: uuid.FromStringOrNil(id), URL: m["url"], Finished: f}, nil
}

func NewRedisPagestore(network, addr string) (*RedisPagestore, error) {
	c, err := pool.New(network, addr, 10)

	if err != nil {
		return nil, fmt.Errorf("unable to create RedisPagestore: %v", err)
	}

	return &RedisPagestore{c}, nil
}
