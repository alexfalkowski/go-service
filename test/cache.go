package test

import (
	"time"
)

type Cache struct {
	Value string
}

func (c *Cache) Contains(_ string) bool {
	return true
}

func (c *Cache) Delete(_ string) error {
	return nil
}

func (c *Cache) Fetch(_ string) (string, error) {
	return c.Value, nil
}

func (c *Cache) FetchMulti(_ []string) map[string]string {
	return map[string]string{}
}

func (c *Cache) Flush() error {
	return nil
}

func (c *Cache) Save(_, _ string, _ time.Duration) error {
	return nil
}

type ErrCache struct{}

func (*ErrCache) Contains(_ string) bool {
	return true
}

func (*ErrCache) Delete(_ string) error {
	return ErrFailed
}

func (*ErrCache) Fetch(_ string) (string, error) {
	return "", ErrFailed
}

func (*ErrCache) FetchMulti(_ []string) map[string]string {
	return map[string]string{}
}

func (*ErrCache) Flush() error {
	return nil
}

func (*ErrCache) Save(_, _ string, _ time.Duration) error {
	return ErrFailed
}
