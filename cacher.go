package gox

import (
	`context`
	`encoding/json`
	`reflect`
	`time`

	`github.com/go-redis/redis/v8`
)

type cacher struct {
	prefix string
	bean   interface{}

	client *redis.Client
}

func NewCacher(prefix string, bean interface{}, client *redis.Client) *cacher {
	return &cacher{
		prefix: prefix,
		bean:   bean,
		client: client,
	}
}

func (c *cacher) Put(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(context.Background(), c.prefix+key, data, time.Hour).Err()
}

func (c *cacher) Get(key string) (interface{}, error) {
	data, err := c.client.Get(context.Background(), c.prefix+key).Bytes()
	if err != nil {
		return nil, err
	}

	bean := reflect.New(reflect.TypeOf(c.bean)).Interface()
	if err = json.Unmarshal(data, bean); err != nil {
		return nil, err
	}

	return bean, nil
}

func (c *cacher) Del(key string) error {
	return c.client.Del(context.Background(), c.prefix+key).Err()
}
